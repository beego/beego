// Package memcached provides a memcached binary protocol client.
package memcached

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/couchbase/gomemcached"
	"github.com/couchbase/goutils/logging"
	"golang.org/x/crypto/pbkdf2"
	"hash"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ClientIface interface {
	Add(vb uint16, key string, flags int, exp int, body []byte) (*gomemcached.MCResponse, error)
	Append(vb uint16, key string, data []byte) (*gomemcached.MCResponse, error)
	Auth(user, pass string) (*gomemcached.MCResponse, error)
	AuthAdvanced(user, pass string) (*gomemcached.MCResponse, error)
	AuthList() (*gomemcached.MCResponse, error)
	AuthPlain(user, pass string) (*gomemcached.MCResponse, error)
	AuthScramSha(user, pass string) (*gomemcached.MCResponse, error)
	AuthSHA(user, pass string, mech string, hashFunc func() hash.Hash, hashSize int) (*gomemcached.MCResponse, error)
	AuthSHA1(user, pass string) (*gomemcached.MCResponse, error)
	AuthSHA256(user, pass string) (*gomemcached.MCResponse, error)
	AuthSHA512(user, pass string) (*gomemcached.MCResponse, error)
	CASNext(vb uint16, k string, exp int, state *CASState) bool
	CAS(vb uint16, k string, f CasFunc, initexp int) (*gomemcached.MCResponse, error)
	Close() error
	Decr(vb uint16, key string, amt, def uint64, exp int) (uint64, error)
	Del(vb uint16, key string) (*gomemcached.MCResponse, error)
	EnableMutationToken() (*gomemcached.MCResponse, error)
	Get(vb uint16, key string) (*gomemcached.MCResponse, error)
	GetAndTouch(vb uint16, key string, exp int) (*gomemcached.MCResponse, error)
	GetBulk(vb uint16, keys []string, rv map[string]*gomemcached.MCResponse) error
	GetMeta(vb uint16, key string) (*gomemcached.MCResponse, error)
	GetRandomDoc() (*gomemcached.MCResponse, error)
	Hijack() io.ReadWriteCloser
	Incr(vb uint16, key string, amt, def uint64, exp int) (uint64, error)
	Observe(vb uint16, key string) (result ObserveResult, err error)
	ObserveSeq(vb uint16, vbuuid uint64) (result *ObserveSeqResult, err error)
	Receive() (*gomemcached.MCResponse, error)
	Send(req *gomemcached.MCRequest) (rv *gomemcached.MCResponse, err error)
	Set(vb uint16, key string, flags int, exp int, body []byte) (*gomemcached.MCResponse, error)
	SetKeepAliveOptions(interval time.Duration)
	SetReadDeadline(t time.Time)
	SelectBucket(bucket string) (*gomemcached.MCResponse, error)
	SetCas(vb uint16, key string, flags int, exp int, cas uint64, body []byte) (*gomemcached.MCResponse, error)
	Stats(key string) ([]StatValue, error)
	StatsMap(key string) (map[string]string, error)
	StatsMapForSpecifiedStats(key string, statsMap map[string]string) error
	Transmit(req *gomemcached.MCRequest) error
	TransmitResponse(res *gomemcached.MCResponse) error

	// UprFeed Related
	NewUprFeed() (*UprFeed, error)
	NewUprFeedIface() (UprFeedIface, error)
	NewUprFeedWithConfig(ackByClient bool) (*UprFeed, error)
	NewUprFeedWithConfigIface(ackByClient bool) (UprFeedIface, error)
	UprGetFailoverLog(vb []uint16) (map[uint16]*FailoverLog, error)
}

// authentication mechanisms
const (
	AuthMech_SHA512 = "SCRAM-SHA512"
	AuthMech_SHA256 = "SCRAM-SHA256"
	AuthMech_SHA1   = "SCRAM-SHA1"
	AuthMech_PLAIN  = "PLAIN"
)

const bufsize = 1024

var UnHealthy uint32 = 0
var Healthy uint32 = 1

var InvalidResponseError = errors.New("the server response is not according to the spec")
var InvalidStartResponseError = errors.New("the server initial response is not according to the spec")
var InvalidStepResponseError = errors.New("the server final response is not according to the spec")

// The Client itself.
type Client struct {
	conn io.ReadWriteCloser
	// use uint32 type so that it can be accessed through atomic APIs
	healthy uint32
	opaque  uint32

	hdrBuf []byte
}

var (
	DefaultDialTimeout = time.Duration(0) // No timeout

	DefaultWriteTimeout = time.Duration(0) // No timeout

	dialFun = func(prot, dest string) (net.Conn, error) {
		return net.DialTimeout(prot, dest, DefaultDialTimeout)
	}
)

// Connect to a memcached server.
func Connect(prot, dest string) (rv *Client, err error) {
	conn, err := dialFun(prot, dest)
	if err != nil {
		return nil, err
	}
	return Wrap(conn)
}

func (c *Client) SetKeepAliveOptions(interval time.Duration) {
	c.conn.(*net.TCPConn).SetKeepAlive(true)
	c.conn.(*net.TCPConn).SetKeepAlivePeriod(interval)
}

func (c *Client) SetReadDeadline(t time.Time) {
	c.conn.(*net.TCPConn).SetReadDeadline(t)
}

// Wrap an existing transport.
func Wrap(rwc io.ReadWriteCloser) (rv *Client, err error) {
	client := &Client{
		conn:   rwc,
		hdrBuf: make([]byte, gomemcached.HDR_LEN),
		opaque: uint32(1),
	}
	client.setHealthy(true)
	return client, nil
}

// Close the connection when you're done.
func (c *Client) Close() error {
	return c.conn.Close()
}

// IsHealthy returns true unless the client is belived to have
// difficulty communicating to its server.
//
// This is useful for connection pools where we want to
// non-destructively determine that a connection may be reused.
func (c Client) IsHealthy() bool {
	healthyState := atomic.LoadUint32(&c.healthy)
	return healthyState == Healthy
}

// Send a custom request and get the response.
func (c *Client) Send(req *gomemcached.MCRequest) (rv *gomemcached.MCResponse, err error) {
	err = c.Transmit(req)
	if err != nil {
		return
	}
	resp, _, err := getResponse(c.conn, c.hdrBuf)
	c.setHealthy(!gomemcached.IsFatal(err))
	return resp, err
}

// Transmit send a request, but does not wait for a response.
func (c *Client) Transmit(req *gomemcached.MCRequest) error {
	if DefaultWriteTimeout > 0 {
		c.conn.(net.Conn).SetWriteDeadline(time.Now().Add(DefaultWriteTimeout))
	}
	_, err := transmitRequest(c.conn, req)
	// clear write deadline to avoid interference with future write operations
	if DefaultWriteTimeout > 0 {
		c.conn.(net.Conn).SetWriteDeadline(time.Time{})
	}
	if err != nil {
		c.setHealthy(false)
	}
	return err
}

// TransmitResponse send a response, does not wait.
func (c *Client) TransmitResponse(res *gomemcached.MCResponse) error {
	if DefaultWriteTimeout > 0 {
		c.conn.(net.Conn).SetWriteDeadline(time.Now().Add(DefaultWriteTimeout))
	}
	_, err := transmitResponse(c.conn, res)
	// clear write deadline to avoid interference with future write operations
	if DefaultWriteTimeout > 0 {
		c.conn.(net.Conn).SetWriteDeadline(time.Time{})
	}
	if err != nil {
		c.setHealthy(false)
	}
	return err
}

// Receive a response
func (c *Client) Receive() (*gomemcached.MCResponse, error) {
	resp, _, err := getResponse(c.conn, c.hdrBuf)
	if err != nil && resp.Status != gomemcached.KEY_ENOENT {
		c.setHealthy(false)
	}
	return resp, err
}

func appendFeatureCode(bytes []byte) []byte {
	bytes = append(bytes, 0, 0)
	binary.BigEndian.PutUint16(bytes[len(bytes)-2:], uint16(0x04))
	return bytes
}

//Send a hello command to enable MutationTokens
func (c *Client) EnableMutationToken() (*gomemcached.MCResponse, error) {
	var payload []byte
	payload = appendFeatureCode(payload)

	return c.Send(&gomemcached.MCRequest{
		Opcode: gomemcached.HELLO,
		Key:    []byte("GoMemcached"),
		Body:   payload,
	})

}

// Get the value for a key.
func (c *Client) Get(vb uint16, key string) (*gomemcached.MCResponse, error) {
	return c.Send(&gomemcached.MCRequest{
		Opcode:  gomemcached.GET,
		VBucket: vb,
		Key:     []byte(key),
	})
}

// Get the value for a key, and update expiry
func (c *Client) GetAndTouch(vb uint16, key string, exp int) (*gomemcached.MCResponse, error) {
	extraBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(extraBuf[0:], uint32(exp))
	return c.Send(&gomemcached.MCRequest{
		Opcode:  gomemcached.GAT,
		VBucket: vb,
		Key:     []byte(key),
		Extras:  extraBuf,
	})
}

// Get metadata for a key
func (c *Client) GetMeta(vb uint16, key string) (*gomemcached.MCResponse, error) {
	return c.Send(&gomemcached.MCRequest{
		Opcode:  gomemcached.GET_META,
		VBucket: vb,
		Key:     []byte(key),
	})
}

// Del deletes a key.
func (c *Client) Del(vb uint16, key string) (*gomemcached.MCResponse, error) {
	return c.Send(&gomemcached.MCRequest{
		Opcode:  gomemcached.DELETE,
		VBucket: vb,
		Key:     []byte(key)})
}

// Get a random document
func (c *Client) GetRandomDoc() (*gomemcached.MCResponse, error) {
	return c.Send(&gomemcached.MCRequest{
		Opcode: 0xB6,
	})
}

// AuthList lists SASL auth mechanisms.
func (c *Client) AuthList() (*gomemcached.MCResponse, error) {
	return c.Send(&gomemcached.MCRequest{
		Opcode: gomemcached.SASL_LIST_MECHS})
}

// Auth performs SASL PLAIN authentication against the server.
func (c *Client) Auth(user, pass string) (*gomemcached.MCResponse, error) {
	res, err := c.AuthList()

	if err != nil {
		return res, err
	}

	authMech := string(res.Body)
	if strings.Index(authMech, AuthMech_PLAIN) != -1 {
		return c.AuthPlain(user, pass)
	}
	return nil, fmt.Errorf("auth mechanism PLAIN not supported")
}

// uses SCRAM-SHA authentication. returns error if SCRAM-SHA is not supported
func (c *Client) AuthScramSha(user, pass string) (*gomemcached.MCResponse, error) {
	res, err := c.AuthList()

	if err != nil {
		return res, err
	}

	authMech := string(res.Body)
	if strings.Index(authMech, AuthMech_SHA512) != -1 {
		return c.AuthSHA512(user, pass)
	} else if strings.Index(authMech, AuthMech_SHA256) != -1 {
		return c.AuthSHA256(user, pass)
	} else if strings.Index(authMech, AuthMech_SHA1) != -1 {
		return c.AuthSHA1(user, pass)
	} else {
		return nil, fmt.Errorf("auth mechanism SCRAM-SHA not supported")
	}
}

// uses SCRAM-SHA authentication if supported. uses plain authentication otherwise
func (c *Client) AuthAdvanced(user, pass string) (*gomemcached.MCResponse, error) {
	res, err := c.AuthList()

	if err != nil {
		return res, err
	}

	authMech := string(res.Body)
	if strings.Index(authMech, AuthMech_SHA512) != -1 {
		return c.AuthSHA512(user, pass)
	} else if strings.Index(authMech, AuthMech_SHA256) != -1 {
		return c.AuthSHA256(user, pass)
	} else if strings.Index(authMech, AuthMech_SHA1) != -1 {
		return c.AuthSHA1(user, pass)
	} else if strings.Index(authMech, AuthMech_PLAIN) != -1 {
		return c.AuthPlain(user, pass)
	}

	return nil, fmt.Errorf("auth mechanism not found")
}

func (c *Client) AuthPlain(user, pass string) (*gomemcached.MCResponse, error) {
	logging.Infof("Using plain authentication for user %v", user)
	return c.Send(&gomemcached.MCRequest{
		Opcode: gomemcached.SASL_AUTH,
		Key:    []byte("PLAIN"),
		Body:   []byte(fmt.Sprintf("\x00%s\x00%s", user, pass))})
}

func (c *Client) AuthSHA512(user, pass string) (*gomemcached.MCResponse, error) {
	logging.Infof("Using scram-sha512 authentication for user %v", user)
	return c.AuthSHA(user, pass, AuthMech_SHA512, sha512.New, 64)
}

func (c *Client) AuthSHA256(user, pass string) (*gomemcached.MCResponse, error) {
	logging.Infof("Using scram-sha256 authentication for user %v", user)
	return c.AuthSHA(user, pass, AuthMech_SHA256, sha256.New, 32)
}

func (c *Client) AuthSHA1(user, pass string) (*gomemcached.MCResponse, error) {
	logging.Infof("Using scram-sha1 authentication for user %v", user)
	return c.AuthSHA(user, pass, AuthMech_SHA1, sha1.New, 20)
}

func (c *Client) AuthSHA(user, pass string, mech string, hashFunc func() hash.Hash, hashSize int) (*gomemcached.MCResponse, error) {
	serverNonce, authMessage, salt, i, err := c.sendStartRequest(user, pass, mech, hashFunc, hashSize)
	if err != nil {
		return nil, err
	}
	return c.sendStepRequest(pass, mech, serverNonce, authMessage, salt, i, hashFunc, hashSize)
}

func (c *Client) sendStartRequest(user, pass string, mech string, hashFunc func() hash.Hash, hashSize int) (serverNonce, authMessage string, salt []byte, i int, err error) {
	var clientNonce string
	clientNonce, err = generateClientNonce(24)
	if err != nil {
		return
	}

	initialMessage := fmt.Sprintf("n,,n=%s,r=%s", user, clientNonce)
	startRequest := &gomemcached.MCRequest{
		Opcode: 0x21,
		Key:    []byte(mech),
		Body:   []byte(initialMessage)}

	// send start request
	var startResponse *gomemcached.MCResponse
	startResponse, err = c.Send(startRequest)
	if err != nil {
		return
	}

	// validate start response
	var encodedSalt string
	serverNonce, encodedSalt, i, err = parseStartResponse(startResponse.Body)
	if err != nil {
		return
	}
	if !strings.HasPrefix(serverNonce, clientNonce) {
		err = errors.New("server nonce is not correct")
		return
	}

	salt, err = base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return
	}

	authMessage = initialMessage[3:] + "," + string(startResponse.Body) + ",c=biws,r=" + serverNonce
	return
}

func (c *Client) sendStepRequest(pass, mech, serverNonce, authMessage string, salt []byte, i int, hashFunc func() hash.Hash, hashSize int) (*gomemcached.MCResponse, error) {
	saltedPassword, clientFinalMessage := c.prepareStepRequest(pass, serverNonce, authMessage, salt, i, hashFunc, hashSize)

	// send step request
	stepRequest := &gomemcached.MCRequest{
		Opcode: 0x22,
		Key:    []byte(mech),
		Body:   []byte(clientFinalMessage)}
	stepResponse, err := c.Send(stepRequest)
	if err != nil {
		return nil, err
	}

	// validate step response
	serverFinalMessage, err := parseStepResponse(stepResponse.Body)
	if err != nil {
		return nil, err
	}
	err = validateServerFinalMessage(serverFinalMessage, saltedPassword, authMessage, hashFunc)
	if err != nil {
		return nil, err
	}

	return stepResponse, nil
}

func (c *Client) prepareStepRequest(pass, serverNonce, authMessage string, salt []byte, i int, hashFunc func() hash.Hash, hashSize int) (saltedPassword []byte, clientFinalMessage string) {
	// prepare step request
	clientFinalMessageBare := "c=biws,r=" + serverNonce
	saltedPassword = pbkdf2.Key([]byte(pass), salt, i, hashSize, hashFunc)
	clientKey := HMACHash([]byte("Client Key"), saltedPassword, hashFunc)
	storedKey := SHAHash(clientKey, hashFunc)
	clientSignature := HMACHash([]byte(authMessage), storedKey, hashFunc)
	clientProof := make([]byte, len(clientSignature))
	for i := 0; i < len(clientSignature); i++ {
		clientProof[i] = clientKey[i] ^ clientSignature[i]
	}
	clientFinalMessage = clientFinalMessageBare + ",p=" + base64.StdEncoding.EncodeToString(clientProof)
	return
}

func generateClientNonce(size int) (string, error) {
	randomBytes := make([]byte, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(randomBytes), nil
}

func parseAuthResponse(body []byte) (map[string]string, error) {
	attributes := make(map[string]string)
	parts := strings.Split(string(body), ",")
	for _, part := range parts {
		// find the index of the FIRST "="
		index := strings.Index(part, "=")
		if index <= 0 {
			return nil, InvalidResponseError
		}
		attributes[part[:index]] = part[index+1:]
	}
	return attributes, nil
}

func parseStartResponse(body []byte) (serverNonce string, encodedSalt string, i int, err error) {
	var attributes map[string]string
	attributes, err = parseAuthResponse(body)
	if err != nil {
		err = InvalidStartResponseError
		return
	}

	i = -1

	for key, value := range attributes {
		switch key {
		case "r":
			serverNonce = value
		case "s":
			encodedSalt = value
		case "i":
			i, err = strconv.Atoi(value)
			if err != nil {
				err = InvalidStartResponseError
				return
			}
		default:
			err = InvalidStartResponseError
			return
		}
	}

	if len(serverNonce) == 0 || len(encodedSalt) == 0 || i < 0 {
		err = InvalidStartResponseError
		return
	}

	return

}

func parseStepResponse(body []byte) (serverFinalMessage string, err error) {
	var attributes map[string]string
	attributes, err = parseAuthResponse(body)
	if err != nil {
		err = InvalidStepResponseError
		return
	}

	for key, value := range attributes {
		switch key {
		case "v":
			serverFinalMessage = value
		default:
			err = InvalidStepResponseError
			return
		}
	}

	if len(serverFinalMessage) == 0 {
		err = InvalidStepResponseError
		return
	}

	return

}

func HMACHash(message []byte, secret []byte, hashFunc func() hash.Hash) []byte {
	h := hmac.New(hashFunc, secret)
	h.Write(message)
	return h.Sum(nil)
}

func SHAHash(message []byte, hashFunc func() hash.Hash) []byte {
	h := hashFunc()
	h.Write(message)
	return h.Sum(nil)
}

func validateServerFinalMessage(serverFinalMessage string, saltedPassword []byte, authMessage string, hashFunc func() hash.Hash) error {
	decodedMessage, err := base64.StdEncoding.DecodeString(serverFinalMessage)
	if err != nil {
		return err
	}
	serverKey := HMACHash([]byte("Server Key"), saltedPassword, hashFunc)
	serverSignature := HMACHash([]byte(authMessage), serverKey, hashFunc)
	if string(decodedMessage) != string(serverSignature) {
		return errors.New("server signature is not correct")
	}
	return nil
}

// select bucket
func (c *Client) SelectBucket(bucket string) (*gomemcached.MCResponse, error) {

	return c.Send(&gomemcached.MCRequest{
		Opcode: gomemcached.SELECT_BUCKET,
		Key:    []byte(fmt.Sprintf("%s", bucket))})
}

func (c *Client) store(opcode gomemcached.CommandCode, vb uint16,
	key string, flags int, exp int, body []byte) (*gomemcached.MCResponse, error) {

	req := &gomemcached.MCRequest{
		Opcode:  opcode,
		VBucket: vb,
		Key:     []byte(key),
		Cas:     0,
		Opaque:  0,
		Extras:  []byte{0, 0, 0, 0, 0, 0, 0, 0},
		Body:    body}

	binary.BigEndian.PutUint64(req.Extras, uint64(flags)<<32|uint64(exp))
	return c.Send(req)
}

func (c *Client) storeCas(opcode gomemcached.CommandCode, vb uint16,
	key string, flags int, exp int, cas uint64, body []byte) (*gomemcached.MCResponse, error) {

	req := &gomemcached.MCRequest{
		Opcode:  opcode,
		VBucket: vb,
		Key:     []byte(key),
		Cas:     cas,
		Opaque:  0,
		Extras:  []byte{0, 0, 0, 0, 0, 0, 0, 0},
		Body:    body}

	binary.BigEndian.PutUint64(req.Extras, uint64(flags)<<32|uint64(exp))
	return c.Send(req)
}

// Incr increments the value at the given key.
func (c *Client) Incr(vb uint16, key string,
	amt, def uint64, exp int) (uint64, error) {

	req := &gomemcached.MCRequest{
		Opcode:  gomemcached.INCREMENT,
		VBucket: vb,
		Key:     []byte(key),
		Extras:  make([]byte, 8+8+4),
	}
	binary.BigEndian.PutUint64(req.Extras[:8], amt)
	binary.BigEndian.PutUint64(req.Extras[8:16], def)
	binary.BigEndian.PutUint32(req.Extras[16:20], uint32(exp))

	resp, err := c.Send(req)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(resp.Body), nil
}

// Decr decrements the value at the given key.
func (c *Client) Decr(vb uint16, key string,
	amt, def uint64, exp int) (uint64, error) {

	req := &gomemcached.MCRequest{
		Opcode:  gomemcached.DECREMENT,
		VBucket: vb,
		Key:     []byte(key),
		Extras:  make([]byte, 8+8+4),
	}
	binary.BigEndian.PutUint64(req.Extras[:8], amt)
	binary.BigEndian.PutUint64(req.Extras[8:16], def)
	binary.BigEndian.PutUint32(req.Extras[16:20], uint32(exp))

	resp, err := c.Send(req)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(resp.Body), nil
}

// Add a value for a key (store if not exists).
func (c *Client) Add(vb uint16, key string, flags int, exp int,
	body []byte) (*gomemcached.MCResponse, error) {
	return c.store(gomemcached.ADD, vb, key, flags, exp, body)
}

// Set the value for a key.
func (c *Client) Set(vb uint16, key string, flags int, exp int,
	body []byte) (*gomemcached.MCResponse, error) {
	return c.store(gomemcached.SET, vb, key, flags, exp, body)
}

// SetCas set the value for a key with cas
func (c *Client) SetCas(vb uint16, key string, flags int, exp int, cas uint64,
	body []byte) (*gomemcached.MCResponse, error) {
	return c.storeCas(gomemcached.SET, vb, key, flags, exp, cas, body)
}

// Append data to the value of a key.
func (c *Client) Append(vb uint16, key string, data []byte) (*gomemcached.MCResponse, error) {
	req := &gomemcached.MCRequest{
		Opcode:  gomemcached.APPEND,
		VBucket: vb,
		Key:     []byte(key),
		Cas:     0,
		Opaque:  0,
		Body:    data}

	return c.Send(req)
}

// GetBulk gets keys in bulk
func (c *Client) GetBulk(vb uint16, keys []string, rv map[string]*gomemcached.MCResponse) error {
	stopch := make(chan bool)
	var wg sync.WaitGroup

	defer func() {
		close(stopch)
		wg.Wait()
	}()

	if (math.MaxInt32 - c.opaque) < (uint32(len(keys)) + 1) {
		c.opaque = uint32(1)
	}

	opStart := c.opaque

	errch := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logging.Infof("Recovered in f %v", r)
			}
			errch <- nil
			wg.Done()
		}()

		ok := true
		for ok {

			select {
			case <-stopch:
				return
			default:
				res, err := c.Receive()
				if err != nil {
					if res.Status != gomemcached.KEY_ENOENT {
						errch <- err
						return
					}
					// continue receiving in case of KEY_ENOENT
				} else if res.Opcode == gomemcached.GET {
					opaque := res.Opaque - opStart
					if opaque < 0 || opaque >= uint32(len(keys)) {
						// Every now and then we seem to be seeing an invalid opaque
						// value returned from the server. When this happens log the error
						// and the calling function will retry the bulkGet. MB-15140
						logging.Errorf(" Invalid opaque Value. Debug info : Res.opaque : %v(%v), Keys %v, Response received %v \n key list %v this key %v", res.Opaque, opaque, len(keys), res, keys, string(res.Body))
						errch <- fmt.Errorf("Out of Bounds error")
						return
					}

					rv[keys[opaque]] = res
				}

				if res.Opcode == gomemcached.NOOP {
					ok = false
				}

			}
		}
	}()

	for _, k := range keys {
		err := c.Transmit(&gomemcached.MCRequest{
			Opcode:  gomemcached.GET,
			VBucket: vb,
			Key:     []byte(k),
			Opaque:  c.opaque,
		})
		if err != nil {
			logging.Errorf(" Transmit failed in GetBulkAll %v", err)
			return err
		}
		c.opaque++
	}

	// finally transmit a NOOP
	err := c.Transmit(&gomemcached.MCRequest{
		Opcode:  gomemcached.NOOP,
		VBucket: vb,
		Opaque:  c.opaque,
	})

	if err != nil {
		logging.Errorf(" Transmit of NOOP failed  %v", err)
		return err
	}
	c.opaque++

	return <-errch
}

// ObservedStatus is the type reported by the Observe method
type ObservedStatus uint8

// Observation status values.
const (
	ObservedNotPersisted     = ObservedStatus(0x00) // found, not persisted
	ObservedPersisted        = ObservedStatus(0x01) // found, persisted
	ObservedNotFound         = ObservedStatus(0x80) // not found (or a persisted delete)
	ObservedLogicallyDeleted = ObservedStatus(0x81) // pending deletion (not persisted yet)
)

// ObserveResult represents the data obtained by an Observe call
type ObserveResult struct {
	Status          ObservedStatus // Whether the value has been persisted/deleted
	Cas             uint64         // Current value's CAS
	PersistenceTime time.Duration  // Node's average time to persist a value
	ReplicationTime time.Duration  // Node's average time to replicate a value
}

// Observe gets the persistence/replication/CAS state of a key
func (c *Client) Observe(vb uint16, key string) (result ObserveResult, err error) {
	// http://www.couchbase.com/wiki/display/couchbase/Observe
	body := make([]byte, 4+len(key))
	binary.BigEndian.PutUint16(body[0:2], vb)
	binary.BigEndian.PutUint16(body[2:4], uint16(len(key)))
	copy(body[4:4+len(key)], key)

	res, err := c.Send(&gomemcached.MCRequest{
		Opcode:  gomemcached.OBSERVE,
		VBucket: vb,
		Body:    body,
	})
	if err != nil {
		return
	}

	// Parse the response data from the body:
	if len(res.Body) < 2+2+1 {
		err = io.ErrUnexpectedEOF
		return
	}
	outVb := binary.BigEndian.Uint16(res.Body[0:2])
	keyLen := binary.BigEndian.Uint16(res.Body[2:4])
	if len(res.Body) < 2+2+int(keyLen)+1+8 {
		err = io.ErrUnexpectedEOF
		return
	}
	outKey := string(res.Body[4 : 4+keyLen])
	if outVb != vb || outKey != key {
		err = fmt.Errorf("observe returned wrong vbucket/key: %d/%q", outVb, outKey)
		return
	}
	result.Status = ObservedStatus(res.Body[4+keyLen])
	result.Cas = binary.BigEndian.Uint64(res.Body[5+keyLen:])
	// The response reuses the Cas field to store time statistics:
	result.PersistenceTime = time.Duration(res.Cas>>32) * time.Millisecond
	result.ReplicationTime = time.Duration(res.Cas&math.MaxUint32) * time.Millisecond
	return
}

// CheckPersistence checks whether a stored value has been persisted to disk yet.
func (result ObserveResult) CheckPersistence(cas uint64, deletion bool) (persisted bool, overwritten bool) {
	switch {
	case result.Status == ObservedNotFound && deletion:
		persisted = true
	case result.Cas != cas:
		overwritten = true
	case result.Status == ObservedPersisted:
		persisted = true
	}
	return
}

// Sequence number based Observe Implementation
type ObserveSeqResult struct {
	Failover           uint8  // Set to 1 if a failover took place
	VbId               uint16 // vbucket id
	Vbuuid             uint64 // vucket uuid
	LastPersistedSeqNo uint64 // last persisted sequence number
	CurrentSeqNo       uint64 // current sequence number
	OldVbuuid          uint64 // Old bucket vbuuid
	LastSeqNo          uint64 // last sequence number received before failover
}

func (c *Client) ObserveSeq(vb uint16, vbuuid uint64) (result *ObserveSeqResult, err error) {
	// http://www.couchbase.com/wiki/display/couchbase/Observe
	body := make([]byte, 8)
	binary.BigEndian.PutUint64(body[0:8], vbuuid)

	res, err := c.Send(&gomemcached.MCRequest{
		Opcode:  gomemcached.OBSERVE_SEQNO,
		VBucket: vb,
		Body:    body,
		Opaque:  0x01,
	})
	if err != nil {
		return
	}

	if res.Status != gomemcached.SUCCESS {
		return nil, fmt.Errorf(" Observe returned error %v", res.Status)
	}

	// Parse the response data from the body:
	if len(res.Body) < (1 + 2 + 8 + 8 + 8) {
		err = io.ErrUnexpectedEOF
		return
	}

	result = &ObserveSeqResult{}
	result.Failover = res.Body[0]
	result.VbId = binary.BigEndian.Uint16(res.Body[1:3])
	result.Vbuuid = binary.BigEndian.Uint64(res.Body[3:11])
	result.LastPersistedSeqNo = binary.BigEndian.Uint64(res.Body[11:19])
	result.CurrentSeqNo = binary.BigEndian.Uint64(res.Body[19:27])

	// in case of failover processing we can have old vbuuid and the last persisted seq number
	if result.Failover == 1 && len(res.Body) >= (1+2+8+8+8+8+8) {
		result.OldVbuuid = binary.BigEndian.Uint64(res.Body[27:35])
		result.LastSeqNo = binary.BigEndian.Uint64(res.Body[35:43])
	}

	return
}

// CasOp is the type of operation to perform on this CAS loop.
type CasOp uint8

const (
	// CASStore instructs the server to store the new value normally
	CASStore = CasOp(iota)
	// CASQuit instructs the client to stop attempting to CAS, leaving value untouched
	CASQuit
	// CASDelete instructs the server to delete the current value
	CASDelete
)

// User specified termination is returned as an error.
func (c CasOp) Error() string {
	switch c {
	case CASStore:
		return "CAS store"
	case CASQuit:
		return "CAS quit"
	case CASDelete:
		return "CAS delete"
	}
	panic("Unhandled value")
}

//////// CAS TRANSFORM

// CASState tracks the state of CAS over several operations.
//
// This is used directly by CASNext and indirectly by CAS
type CASState struct {
	initialized bool   // false on the first call to CASNext, then true
	Value       []byte // Current value of key; update in place to new value
	Cas         uint64 // Current CAS value of key
	Exists      bool   // Does a value exist for the key? (If not, Value will be nil)
	Err         error  // Error, if any, after CASNext returns false
	resp        *gomemcached.MCResponse
}

// CASNext is a non-callback, loop-based version of CAS method.
//
//  Usage is like this:
//
// var state memcached.CASState
// for client.CASNext(vb, key, exp, &state) {
//     state.Value = some_mutation(state.Value)
// }
// if state.Err != nil { ... }
func (c *Client) CASNext(vb uint16, k string, exp int, state *CASState) bool {
	if state.initialized {
		if !state.Exists {
			// Adding a new key:
			if state.Value == nil {
				state.Cas = 0
				return false // no-op (delete of non-existent value)
			}
			state.resp, state.Err = c.Add(vb, k, 0, exp, state.Value)
		} else {
			// Updating / deleting a key:
			req := &gomemcached.MCRequest{
				Opcode:  gomemcached.DELETE,
				VBucket: vb,
				Key:     []byte(k),
				Cas:     state.Cas}
			if state.Value != nil {
				req.Opcode = gomemcached.SET
				req.Opaque = 0
				req.Extras = []byte{0, 0, 0, 0, 0, 0, 0, 0}
				req.Body = state.Value

				flags := 0
				binary.BigEndian.PutUint64(req.Extras, uint64(flags)<<32|uint64(exp))
			}
			state.resp, state.Err = c.Send(req)
		}

		// If the response status is KEY_EEXISTS or NOT_STORED there's a conflict and we'll need to
		// get the new value (below). Otherwise, we're done (either success or failure) so return:
		if !(state.resp != nil && (state.resp.Status == gomemcached.KEY_EEXISTS ||
			state.resp.Status == gomemcached.NOT_STORED)) {
			state.Cas = state.resp.Cas
			return false // either success or fatal error
		}
	}

	// Initial call, or after a conflict: GET the current value and CAS and return them:
	state.initialized = true
	if state.resp, state.Err = c.Get(vb, k); state.Err == nil {
		state.Exists = true
		state.Value = state.resp.Body
		state.Cas = state.resp.Cas
	} else if state.resp != nil && state.resp.Status == gomemcached.KEY_ENOENT {
		state.Err = nil
		state.Exists = false
		state.Value = nil
		state.Cas = 0
	} else {
		return false // fatal error
	}
	return true // keep going...
}

// CasFunc is type type of function to perform a CAS transform.
//
// Input is the current value, or nil if no value exists.
// The function should return the new value (if any) to set, and the store/quit/delete operation.
type CasFunc func(current []byte) ([]byte, CasOp)

// CAS performs a CAS transform with the given function.
//
// If the value does not exist, a nil current value will be sent to f.
func (c *Client) CAS(vb uint16, k string, f CasFunc,
	initexp int) (*gomemcached.MCResponse, error) {
	var state CASState
	for c.CASNext(vb, k, initexp, &state) {
		newValue, operation := f(state.Value)
		if operation == CASQuit || (operation == CASDelete && state.Value == nil) {
			return nil, operation
		}
		state.Value = newValue
	}
	return state.resp, state.Err
}

// StatValue is one of the stats returned from the Stats method.
type StatValue struct {
	// The stat key
	Key string
	// The stat value
	Val string
}

// Stats requests server-side stats.
//
// Use "" as the stat key for toplevel stats.
func (c *Client) Stats(key string) ([]StatValue, error) {
	rv := make([]StatValue, 0, 128)

	req := &gomemcached.MCRequest{
		Opcode: gomemcached.STAT,
		Key:    []byte(key),
		Opaque: 918494,
	}

	err := c.Transmit(req)
	if err != nil {
		return rv, err
	}

	for {
		res, _, err := getResponse(c.conn, c.hdrBuf)
		if err != nil {
			return rv, err
		}
		k := string(res.Key)
		if k == "" {
			break
		}
		rv = append(rv, StatValue{
			Key: k,
			Val: string(res.Body),
		})
	}
	return rv, nil
}

// StatsMap requests server-side stats similarly to Stats, but returns
// them as a map.
//
// Use "" as the stat key for toplevel stats.
func (c *Client) StatsMap(key string) (map[string]string, error) {
	rv := make(map[string]string)

	req := &gomemcached.MCRequest{
		Opcode: gomemcached.STAT,
		Key:    []byte(key),
		Opaque: 918494,
	}

	err := c.Transmit(req)
	if err != nil {
		return rv, err
	}

	for {
		res, _, err := getResponse(c.conn, c.hdrBuf)
		if err != nil {
			return rv, err
		}
		k := string(res.Key)
		if k == "" {
			break
		}
		rv[k] = string(res.Body)
	}

	return rv, nil
}

// instead of returning a new statsMap, simply populate passed in statsMap, which contains all the keys
// for which stats needs to be retrieved
func (c *Client) StatsMapForSpecifiedStats(key string, statsMap map[string]string) error {

	// clear statsMap
	for key, _ := range statsMap {
		statsMap[key] = ""
	}

	req := &gomemcached.MCRequest{
		Opcode: gomemcached.STAT,
		Key:    []byte(key),
		Opaque: 918494,
	}

	err := c.Transmit(req)
	if err != nil {
		return err
	}

	for {
		res, _, err := getResponse(c.conn, c.hdrBuf)
		if err != nil {
			return err
		}
		k := string(res.Key)
		if k == "" {
			break
		}
		if _, ok := statsMap[k]; ok {
			statsMap[k] = string(res.Body)
		}
	}

	return nil
}

// Hijack exposes the underlying connection from this client.
//
// It also marks the connection as unhealthy since the client will
// have lost control over the connection and can't otherwise verify
// things are in good shape for connection pools.
func (c *Client) Hijack() io.ReadWriteCloser {
	c.setHealthy(false)
	return c.conn
}

func (c *Client) setHealthy(healthy bool) {
	healthyState := UnHealthy
	if healthy {
		healthyState = Healthy
	}
	atomic.StoreUint32(&c.healthy, healthyState)
}
