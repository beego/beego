// go implementation of upr client.
// See https://github.com/couchbaselabs/cbupr/blob/master/transport-spec.md
// TODO
// 1. Use a pool allocator to avoid garbage
package memcached

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/couchbase/gomemcached"
	"github.com/couchbase/goutils/logging"
	"strconv"
	"sync"
)

const uprMutationExtraLen = 30
const uprDeletetionExtraLen = 18
const uprSnapshotExtraLen = 20
const bufferAckThreshold = 0.2
const opaqueOpen = 0xBEAF0001
const opaqueFailover = 0xDEADBEEF
const uprDefaultNoopInterval = 120

// UprEvent memcached events for UPR streams.
type UprEvent struct {
	Opcode       gomemcached.CommandCode // Type of event
	Status       gomemcached.Status      // Response status
	VBucket      uint16                  // VBucket this event applies to
	DataType     uint8                   // data type
	Opaque       uint16                  // 16 MSB of opaque
	VBuuid       uint64                  // This field is set by downstream
	Flags        uint32                  // Item flags
	Expiry       uint32                  // Item expiration time
	Key, Value   []byte                  // Item key/value
	OldValue     []byte                  // TODO: TBD: old document value
	Cas          uint64                  // CAS value of the item
	Seqno        uint64                  // sequence number of the mutation
	RevSeqno     uint64                  // rev sequence number : deletions
	LockTime     uint32                  // Lock time
	MetadataSize uint16                  // Metadata size
	SnapstartSeq uint64                  // start sequence number of this snapshot
	SnapendSeq   uint64                  // End sequence number of the snapshot
	SnapshotType uint32                  // 0: disk 1: memory
	FailoverLog  *FailoverLog            // Failover log containing vvuid and sequnce number
	Error        error                   // Error value in case of a failure
	ExtMeta      []byte
	AckSize      uint32 // The number of bytes that can be Acked to DCP
}

// UprStream is per stream data structure over an UPR Connection.
type UprStream struct {
	Vbucket   uint16 // Vbucket id
	Vbuuid    uint64 // vbucket uuid
	StartSeq  uint64 // start sequence number
	EndSeq    uint64 // end sequence number
	connected bool
}

// UprFeed represents an UPR feed. A feed contains a connection to a single
// host and multiple vBuckets
type UprFeed struct {
	// lock for feed.vbstreams
	muVbstreams sync.RWMutex
	// lock for feed.closed
	muClosed   sync.RWMutex
	C           <-chan *UprEvent            // Exported channel for receiving UPR events
	vbstreams   map[uint16]*UprStream       // vb->stream mapping
	closer      chan bool                   // closer
	conn        *Client                     // connection to UPR producer
	Error       error                       // error
	bytesRead   uint64                      // total bytes read on this connection
	toAckBytes  uint32                      // bytes client has read
	maxAckBytes uint32                      // Max buffer control ack bytes
	stats       UprStats                    // Stats for upr client
	transmitCh  chan *gomemcached.MCRequest // transmit command channel
	transmitCl  chan bool                   //  closer channel for transmit go-routine
	closed      bool                        // flag indicating whether the feed has been closed
	// flag indicating whether client of upr feed will send ack to upr feed
	// if flag is true, upr feed will use ack from client to determine whether/when to send ack to DCP
	// if flag is false, upr feed will track how many bytes it has sent to client
	// and use that to determine whether/when to send ack to DCP
	ackByClient bool
}

// Exported interface - to allow for mocking
type UprFeedIface interface {
	Close()
	Closed() bool
	CloseStream(vbno, opaqueMSB uint16) error
	GetError() error
	GetUprStats() *UprStats
	IncrementAckBytes(bytes uint32) error
	GetUprEventCh() <-chan *UprEvent
	StartFeed() error
	StartFeedWithConfig(datachan_len int) error
	UprOpen(name string, sequence uint32, bufSize uint32) error
	UprOpenWithXATTR(name string, sequence uint32, bufSize uint32) error
	UprRequestStream(vbno, opaqueMSB uint16, flags uint32, vuuid, startSequence, endSequence, snapStart, snapEnd uint64) error
}

type UprStats struct {
	TotalBytes         uint64
	TotalMutation      uint64
	TotalBufferAckSent uint64
	TotalSnapShot      uint64
}

// FailoverLog containing vvuid and sequnce number
type FailoverLog [][2]uint64

// error codes
var ErrorInvalidLog = errors.New("couchbase.errorInvalidLog")

func (flogp *FailoverLog) Latest() (vbuuid, seqno uint64, err error) {
	if flogp != nil {
		flog := *flogp
		latest := flog[len(flog)-1]
		return latest[0], latest[1], nil
	}
	return vbuuid, seqno, ErrorInvalidLog
}

func makeUprEvent(rq gomemcached.MCRequest, stream *UprStream) *UprEvent {
	event := &UprEvent{
		Opcode:   rq.Opcode,
		VBucket:  stream.Vbucket,
		VBuuid:   stream.Vbuuid,
		Key:      rq.Key,
		Value:    rq.Body,
		Cas:      rq.Cas,
		ExtMeta:  rq.ExtMeta,
		DataType: rq.DataType,
		AckSize:  uint32(rq.Size()),
	}
	// 16 LSBits are used by client library to encode vbucket number.
	// 16 MSBits are left for application to multiplex on opaque value.
	event.Opaque = appOpaque(rq.Opaque)

	if len(rq.Extras) >= uprMutationExtraLen &&
		event.Opcode == gomemcached.UPR_MUTATION {

		event.Seqno = binary.BigEndian.Uint64(rq.Extras[:8])
		event.RevSeqno = binary.BigEndian.Uint64(rq.Extras[8:16])
		event.Flags = binary.BigEndian.Uint32(rq.Extras[16:20])
		event.Expiry = binary.BigEndian.Uint32(rq.Extras[20:24])
		event.LockTime = binary.BigEndian.Uint32(rq.Extras[24:28])
		event.MetadataSize = binary.BigEndian.Uint16(rq.Extras[28:30])

	} else if len(rq.Extras) >= uprDeletetionExtraLen &&
		event.Opcode == gomemcached.UPR_DELETION ||
		event.Opcode == gomemcached.UPR_EXPIRATION {

		event.Seqno = binary.BigEndian.Uint64(rq.Extras[:8])
		event.RevSeqno = binary.BigEndian.Uint64(rq.Extras[8:16])
		event.MetadataSize = binary.BigEndian.Uint16(rq.Extras[16:18])

	} else if len(rq.Extras) >= uprSnapshotExtraLen &&
		event.Opcode == gomemcached.UPR_SNAPSHOT {

		event.SnapstartSeq = binary.BigEndian.Uint64(rq.Extras[:8])
		event.SnapendSeq = binary.BigEndian.Uint64(rq.Extras[8:16])
		event.SnapshotType = binary.BigEndian.Uint32(rq.Extras[16:20])
	}

	return event
}

func (event *UprEvent) String() string {
	name := gomemcached.CommandNames[event.Opcode]
	if name == "" {
		name = fmt.Sprintf("#%d", event.Opcode)
	}
	return name
}

func (feed *UprFeed) sendCommands(mc *Client) {
	transmitCh := feed.transmitCh
	transmitCl := feed.transmitCl
loop:
	for {
		select {
		case command := <-transmitCh:
			if err := mc.Transmit(command); err != nil {
				logging.Errorf("Failed to transmit command %s. Error %s", command.Opcode.String(), err.Error())
				// get feed to close and runFeed routine to exit
				feed.Close()
				break loop
			}

		case <-transmitCl:
			break loop
		}
	}

	// After sendCommands exits, write to transmitCh will block forever
	// when we write to transmitCh, e.g., at CloseStream(), we need to check feed closure to have an exit route

	logging.Infof("sendCommands exiting")
}

// NewUprFeed creates a new UPR Feed.
// TODO: Describe side-effects on bucket instance and its connection pool.
func (mc *Client) NewUprFeed() (*UprFeed, error) {
	return mc.NewUprFeedWithConfig(false /*ackByClient*/)
}

func (mc *Client) NewUprFeedWithConfig(ackByClient bool) (*UprFeed, error) {

	feed := &UprFeed{
		conn:        mc,
		closer:      make(chan bool, 1),
		vbstreams:   make(map[uint16]*UprStream),
		transmitCh:  make(chan *gomemcached.MCRequest),
		transmitCl:  make(chan bool),
		ackByClient: ackByClient,
	}

	go feed.sendCommands(mc)
	return feed, nil
}

func (mc *Client) NewUprFeedIface() (UprFeedIface, error) {
	return mc.NewUprFeed()
}

func (mc *Client) NewUprFeedWithConfigIface(ackByClient bool) (UprFeedIface, error) {
	return mc.NewUprFeedWithConfig(ackByClient)
}

func doUprOpen(mc *Client, name string, sequence uint32, enableXATTR bool) error {

	rq := &gomemcached.MCRequest{
		Opcode: gomemcached.UPR_OPEN,
		Key:    []byte(name),
		Opaque: opaqueOpen,
	}

	rq.Extras = make([]byte, 8)
	binary.BigEndian.PutUint32(rq.Extras[:4], sequence)

	// opens a producer type connection
	flags := gomemcached.DCP_PRODUCER
	if enableXATTR {
		// set DCP_OPEN_INCLUDE_XATTRS bit in flags
		flags = flags | gomemcached.DCP_OPEN_INCLUDE_XATTRS
	}
	binary.BigEndian.PutUint32(rq.Extras[4:], flags)

	if err := mc.Transmit(rq); err != nil {
		return err
	}

	if res, err := mc.Receive(); err != nil {
		return err
	} else if res.Opcode != gomemcached.UPR_OPEN {
		return fmt.Errorf("unexpected #opcode %v", res.Opcode)
	} else if rq.Opaque != res.Opaque {
		return fmt.Errorf("opaque mismatch, %v over %v", res.Opaque, res.Opaque)
	} else if res.Status != gomemcached.SUCCESS {
		return fmt.Errorf("error %v", res.Status)
	}

	return nil
}

// UprOpen to connect with a UPR producer.
// Name: name of te UPR connection
// sequence: sequence number for the connection
// bufsize: max size of the application
func (feed *UprFeed) UprOpen(name string, sequence uint32, bufSize uint32) error {
	return feed.uprOpen(name, sequence, bufSize, false /*enableXATTR*/)
}

// UprOpen with XATTR enabled.
func (feed *UprFeed) UprOpenWithXATTR(name string, sequence uint32, bufSize uint32) error {
	return feed.uprOpen(name, sequence, bufSize, true /*enableXATTR*/)
}

func (feed *UprFeed) uprOpen(name string, sequence uint32, bufSize uint32, enableXATTR bool) error {
	mc := feed.conn

	var err error
	if err = doUprOpen(mc, name, sequence, enableXATTR); err != nil {
		return err
	}

	// send a UPR control message to set the window size for the this connection
	if bufSize > 0 {
		rq := &gomemcached.MCRequest{
			Opcode: gomemcached.UPR_CONTROL,
			Key:    []byte("connection_buffer_size"),
			Body:   []byte(strconv.Itoa(int(bufSize))),
		}
		err = feed.writeToTransmitCh(rq)
		if err != nil {
			return err
		}
		feed.maxAckBytes = uint32(bufferAckThreshold * float32(bufSize))
	}

	// enable noop and set noop interval
	rq := &gomemcached.MCRequest{
		Opcode: gomemcached.UPR_CONTROL,
		Key:    []byte("enable_noop"),
		Body:   []byte("true"),
	}
	err = feed.writeToTransmitCh(rq)
	if err != nil {
		return err
	}

	rq = &gomemcached.MCRequest{
		Opcode: gomemcached.UPR_CONTROL,
		Key:    []byte("set_noop_interval"),
		Body:   []byte(strconv.Itoa(int(uprDefaultNoopInterval))),
	}
	err = feed.writeToTransmitCh(rq)
	if err != nil {
		return err
	}

	return nil
}

// UprGetFailoverLog for given list of vbuckets.
func (mc *Client) UprGetFailoverLog(
	vb []uint16) (map[uint16]*FailoverLog, error) {

	rq := &gomemcached.MCRequest{
		Opcode: gomemcached.UPR_FAILOVERLOG,
		Opaque: opaqueFailover,
	}

	if err := doUprOpen(mc, "FailoverLog", 0, false); err != nil {
		return nil, fmt.Errorf("UPR_OPEN Failed %s", err.Error())
	}

	failoverLogs := make(map[uint16]*FailoverLog)
	for _, vBucket := range vb {
		rq.VBucket = vBucket
		if err := mc.Transmit(rq); err != nil {
			return nil, err
		}
		res, err := mc.Receive()

		if err != nil {
			return nil, fmt.Errorf("failed to receive %s", err.Error())
		} else if res.Opcode != gomemcached.UPR_FAILOVERLOG || res.Status != gomemcached.SUCCESS {
			return nil, fmt.Errorf("unexpected #opcode %v", res.Opcode)
		}

		flog, err := parseFailoverLog(res.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to parse failover logs for vb %d", vb)
		}
		failoverLogs[vBucket] = flog
	}

	return failoverLogs, nil
}

// UprRequestStream for a single vbucket.
func (feed *UprFeed) UprRequestStream(vbno, opaqueMSB uint16, flags uint32,
	vuuid, startSequence, endSequence, snapStart, snapEnd uint64) error {

	rq := &gomemcached.MCRequest{
		Opcode:  gomemcached.UPR_STREAMREQ,
		VBucket: vbno,
		Opaque:  composeOpaque(vbno, opaqueMSB),
	}

	rq.Extras = make([]byte, 48) // #Extras
	binary.BigEndian.PutUint32(rq.Extras[:4], flags)
	binary.BigEndian.PutUint32(rq.Extras[4:8], uint32(0))
	binary.BigEndian.PutUint64(rq.Extras[8:16], startSequence)
	binary.BigEndian.PutUint64(rq.Extras[16:24], endSequence)
	binary.BigEndian.PutUint64(rq.Extras[24:32], vuuid)
	binary.BigEndian.PutUint64(rq.Extras[32:40], snapStart)
	binary.BigEndian.PutUint64(rq.Extras[40:48], snapEnd)

	stream := &UprStream{
		Vbucket:  vbno,
		Vbuuid:   vuuid,
		StartSeq: startSequence,
		EndSeq:   endSequence,
	}

	feed.muVbstreams.Lock()
	// Any client that has ever called this method, regardless of return code,
	// should expect a potential UPR_CLOSESTREAM message due to this new map entry prior to Transmit.
	feed.vbstreams[vbno] = stream
	feed.muVbstreams.Unlock()

	if err := feed.conn.Transmit(rq); err != nil {
		logging.Errorf("Error in StreamRequest %s", err.Error())
		// If an error occurs during transmit, then the UPRFeed will keep the stream
		// in the vbstreams map. This is to prevent nil lookup from any previously
		// sent stream requests.
		return err
	}

	return nil
}

// CloseStream for specified vbucket.
func (feed *UprFeed) CloseStream(vbno, opaqueMSB uint16) error {

	err := feed.validateCloseStream(vbno)
	if err != nil {
		logging.Infof("CloseStream for %v has been skipped because of error %v", vbno, err)
		return err
	}

	closeStream := &gomemcached.MCRequest{
		Opcode:  gomemcached.UPR_CLOSESTREAM,
		VBucket: vbno,
		Opaque:  composeOpaque(vbno, opaqueMSB),
	}

	feed.writeToTransmitCh(closeStream)

	return nil
}

func (feed *UprFeed) GetUprEventCh() <-chan *UprEvent {
	return feed.C
}

func (feed *UprFeed) GetError() error {
	return feed.Error
}

func (feed *UprFeed) validateCloseStream(vbno uint16) error {
	feed.muVbstreams.RLock()
	defer feed.muVbstreams.RUnlock()

	if feed.vbstreams[vbno] == nil {
		return fmt.Errorf("Stream for vb %d has not been requested", vbno)
	}

	return nil
}

func (feed *UprFeed) writeToTransmitCh(rq *gomemcached.MCRequest) error {
	// write to transmitCh may block forever if sendCommands has exited
	// check for feed closure to have an exit route in this case
	select {
	case <-feed.closer:
		errMsg := fmt.Sprintf("Abort sending request to transmitCh because feed has been closed. request=%v", rq)
		logging.Infof(errMsg)
		return errors.New(errMsg)
	case feed.transmitCh <- rq:
	}
	return nil
}

// StartFeed to start the upper feed.
func (feed *UprFeed) StartFeed() error {
	return feed.StartFeedWithConfig(10)
}

func (feed *UprFeed) StartFeedWithConfig(datachan_len int) error {
	ch := make(chan *UprEvent, datachan_len)
	feed.C = ch
	go feed.runFeed(ch)
	return nil
}

func parseFailoverLog(body []byte) (*FailoverLog, error) {

	if len(body)%16 != 0 {
		err := fmt.Errorf("invalid body length %v, in failover-log", len(body))
		return nil, err
	}
	log := make(FailoverLog, len(body)/16)
	for i, j := 0, 0; i < len(body); i += 16 {
		vuuid := binary.BigEndian.Uint64(body[i : i+8])
		seqno := binary.BigEndian.Uint64(body[i+8 : i+16])
		log[j] = [2]uint64{vuuid, seqno}
		j++
	}
	return &log, nil
}

func handleStreamRequest(
	res *gomemcached.MCResponse,
	headerBuf []byte,
) (gomemcached.Status, uint64, *FailoverLog, error) {

	var rollback uint64
	var err error

	switch {
	case res.Status == gomemcached.ROLLBACK:
		logging.Infof("Rollback response. body=%v, headerBuf=%v\n", res.Body, headerBuf)
		rollback = binary.BigEndian.Uint64(res.Body)
		logging.Infof("Rollback %v for vb %v\n", rollback, res.Opaque)
		return res.Status, rollback, nil, nil

	case res.Status != gomemcached.SUCCESS:
		err = fmt.Errorf("unexpected status %v, for %v", res.Status, res.Opaque)
		return res.Status, 0, nil, err
	}

	flog, err := parseFailoverLog(res.Body[:])
	return res.Status, rollback, flog, err
}

// generate stream end responses for all active vb streams
func (feed *UprFeed) doStreamClose(ch chan *UprEvent) {
	feed.muVbstreams.RLock()

	uprEvents := make([]*UprEvent, len(feed.vbstreams))
	index := 0
	for vbno, stream := range feed.vbstreams {
		uprEvent := &UprEvent{
			VBucket: vbno,
			VBuuid:  stream.Vbuuid,
			Opcode:  gomemcached.UPR_STREAMEND,
		}
		uprEvents[index] = uprEvent
		index++
	}

	// release the lock before sending uprEvents to ch, which may block
	feed.muVbstreams.RUnlock()

loop:
	for _, uprEvent := range uprEvents {
		select {
		case ch <- uprEvent:
		case <-feed.closer:
			logging.Infof("Feed has been closed. Aborting doStreamClose.")
			break loop
		}
	}
}

func (feed *UprFeed) runFeed(ch chan *UprEvent) {
	defer close(ch)
	var headerBuf [gomemcached.HDR_LEN]byte
	var pkt gomemcached.MCRequest
	var event *UprEvent

	mc := feed.conn.Hijack()
	uprStats := &feed.stats

loop:
	for {
		select {
		case <-feed.closer:
			logging.Infof("Feed has been closed. Exiting.")
			break loop
		default:
			sendAck := false
			bytes, err := pkt.Receive(mc, headerBuf[:])
			if err != nil {
				logging.Errorf("Error in receive %s", err.Error())
				feed.Error = err
				// send all the stream close messages to the client
				feed.doStreamClose(ch)
				break loop
			} else {
				event = nil
				res := &gomemcached.MCResponse{
					Opcode: pkt.Opcode,
					Cas:    pkt.Cas,
					Opaque: pkt.Opaque,
					Status: gomemcached.Status(pkt.VBucket),
					Extras: pkt.Extras,
					Key:    pkt.Key,
					Body:   pkt.Body,
				}

				vb := vbOpaque(pkt.Opaque)
				uprStats.TotalBytes = uint64(bytes)

				feed.muVbstreams.RLock()
				stream := feed.vbstreams[vb]
				feed.muVbstreams.RUnlock()

				switch pkt.Opcode {
				case gomemcached.UPR_STREAMREQ:
					if stream == nil {
						logging.Infof("Stream not found for vb %d: %#v", vb, pkt)
						break loop
					}
					status, rb, flog, err := handleStreamRequest(res, headerBuf[:])
					if status == gomemcached.ROLLBACK {
						event = makeUprEvent(pkt, stream)
						event.Status = status
						// rollback stream
						logging.Infof("UPR_STREAMREQ with rollback %d for vb %d Failed: %v", rb, vb, err)
						// delete the stream from the vbmap for the feed
						feed.muVbstreams.Lock()
						delete(feed.vbstreams, vb)
						feed.muVbstreams.Unlock()

					} else if status == gomemcached.SUCCESS {
						event = makeUprEvent(pkt, stream)
						event.Seqno = stream.StartSeq
						event.FailoverLog = flog
						event.Status = status
						stream.connected = true
						logging.Infof("UPR_STREAMREQ for vb %d successful", vb)

					} else if err != nil {
						logging.Errorf("UPR_STREAMREQ for vbucket %d erro %s", vb, err.Error())
						event = &UprEvent{
							Opcode:  gomemcached.UPR_STREAMREQ,
							Status:  status,
							VBucket: vb,
							Error:   err,
						}
						// delete the stream
						feed.muVbstreams.Lock()
						delete(feed.vbstreams, vb)
						feed.muVbstreams.Unlock()
					}

				case gomemcached.UPR_MUTATION,
					gomemcached.UPR_DELETION,
					gomemcached.UPR_EXPIRATION:
					if stream == nil {
						logging.Infof("Stream not found for vb %d: %#v", vb, pkt)
						break loop
					}
					event = makeUprEvent(pkt, stream)
					uprStats.TotalMutation++
					sendAck = true

				case gomemcached.UPR_STREAMEND:
					if stream == nil {
						logging.Infof("Stream not found for vb %d: %#v", vb, pkt)
						break loop
					}
					//stream has ended
					event = makeUprEvent(pkt, stream)
					logging.Infof("Stream Ended for vb %d", vb)
					sendAck = true

					feed.muVbstreams.Lock()
					delete(feed.vbstreams, vb)
					feed.muVbstreams.Unlock()

				case gomemcached.UPR_SNAPSHOT:
					if stream == nil {
						logging.Infof("Stream not found for vb %d: %#v", vb, pkt)
						break loop
					}
					// snapshot marker
					event = makeUprEvent(pkt, stream)
					uprStats.TotalSnapShot++
					sendAck = true

				case gomemcached.UPR_FLUSH:
					if stream == nil {
						logging.Infof("Stream not found for vb %d: %#v", vb, pkt)
						break loop
					}
					// special processing for flush ?
					event = makeUprEvent(pkt, stream)

				case gomemcached.UPR_CLOSESTREAM:
					if stream == nil {
						logging.Infof("Stream not found for vb %d: %#v", vb, pkt)
						break loop
					}
					event = makeUprEvent(pkt, stream)
					event.Opcode = gomemcached.UPR_STREAMEND // opcode re-write !!
					logging.Infof("Stream Closed for vb %d StreamEnd simulated", vb)
					sendAck = true

					feed.muVbstreams.Lock()
					delete(feed.vbstreams, vb)
					feed.muVbstreams.Unlock()

				case gomemcached.UPR_ADDSTREAM:
					logging.Infof("Opcode %v not implemented", pkt.Opcode)

				case gomemcached.UPR_CONTROL, gomemcached.UPR_BUFFERACK:
					if res.Status != gomemcached.SUCCESS {
						logging.Infof("Opcode %v received status %d", pkt.Opcode.String(), res.Status)
					}

				case gomemcached.UPR_NOOP:
					// send a NOOP back
					noop := &gomemcached.MCResponse{
						Opcode: gomemcached.UPR_NOOP,
						Opaque: pkt.Opaque,
					}

					if err := feed.conn.TransmitResponse(noop); err != nil {
						logging.Warnf("failed to transmit command %s. Error %s", noop.Opcode.String(), err.Error())
					}
				default:
					logging.Infof("Recived an unknown response for vbucket %d", vb)
				}
			}

			if event != nil {
				select {
				case ch <- event:
				case <-feed.closer:
					logging.Infof("Feed has been closed. Skip sending events. Exiting.")
					break loop
				}

				feed.muVbstreams.RLock()
				l := len(feed.vbstreams)
				feed.muVbstreams.RUnlock()

				if event.Opcode == gomemcached.UPR_CLOSESTREAM && l == 0 {
					logging.Infof("No more streams")
				}
			}

			if !feed.ackByClient {
				// if client does not ack, use the size of data sent to client to determine if ack to dcp is needed
				feed.sendBufferAckIfNeeded(sendAck, uint32(bytes))
			}
		}
	}

	// make sure that feed is closed before we signal transmitCl and exit runFeed
	feed.Close()

	close(feed.transmitCl)
	logging.Infof("runFeed exiting")
}

// Client, after setting ackByClient flag to true in NewUprFeedWithConfig() call,
// can call this API to notify gomemcached that the client has completed processing
// of a number of bytes
// This API is not thread safe. Caller should NOT have more than one go rountine calling this API
func (feed *UprFeed) IncrementAckBytes(bytes uint32) error {
	if !feed.ackByClient {
		return errors.New("Upr feed does not have ackByclient flag set")
	}
	feed.sendBufferAckIfNeeded(true, bytes)
	return nil
}

// send buffer ack if enough ack bytes have been accumulated
func (feed *UprFeed) sendBufferAckIfNeeded(sendAck bool, bytes uint32) {
	if sendAck {
		totalBytes := feed.toAckBytes + bytes
		if totalBytes > feed.maxAckBytes {
			feed.toAckBytes = 0
			feed.sendBufferAck(totalBytes)
		} else {
			feed.toAckBytes = totalBytes
		}
	}
}

// send buffer ack to dcp
func (feed *UprFeed) sendBufferAck(sendSize uint32) {
	bufferAck := &gomemcached.MCRequest{
		Opcode: gomemcached.UPR_BUFFERACK,
	}
	bufferAck.Extras = make([]byte, 4)
	binary.BigEndian.PutUint32(bufferAck.Extras[:4], uint32(sendSize))
	feed.writeToTransmitCh(bufferAck)
	feed.stats.TotalBufferAckSent++
}

func (feed *UprFeed) GetUprStats() *UprStats {
	return &feed.stats
}

func composeOpaque(vbno, opaqueMSB uint16) uint32 {
	return (uint32(opaqueMSB) << 16) | uint32(vbno)
}

func appOpaque(opq32 uint32) uint16 {
	return uint16((opq32 & 0xFFFF0000) >> 16)
}

func vbOpaque(opq32 uint32) uint16 {
	return uint16(opq32 & 0xFFFF)
}

// Close this UprFeed.
func (feed *UprFeed) Close() {
	feed.muClosed.Lock()
	defer feed.muClosed.Unlock()
	if !feed.closed {
		close(feed.closer)
		feed.closed = true
	}
}

// check if the UprFeed has been closed
func (feed *UprFeed) Closed() bool {
	feed.muClosed.RLock()
	defer feed.muClosed.RUnlock()
	return feed.closed
}
