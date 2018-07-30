package redis

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/internal"
	"github.com/go-redis/redis/internal/pool"
	"github.com/go-redis/redis/internal/proto"
	"github.com/go-redis/redis/internal/util"
)

type Cmder interface {
	Name() string
	Args() []interface{}
	stringArg(int) string

	readReply(*pool.Conn) error
	setErr(error)

	readTimeout() *time.Duration

	Err() error
	fmt.Stringer
}

func setCmdsErr(cmds []Cmder, e error) {
	for _, cmd := range cmds {
		if cmd.Err() == nil {
			cmd.setErr(e)
		}
	}
}

func firstCmdsErr(cmds []Cmder) error {
	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			return err
		}
	}
	return nil
}

func writeCmd(cn *pool.Conn, cmds ...Cmder) error {
	cn.Wb.Reset()
	for _, cmd := range cmds {
		if err := cn.Wb.Append(cmd.Args()); err != nil {
			return err
		}
	}

	_, err := cn.Write(cn.Wb.Bytes())
	return err
}

func cmdString(cmd Cmder, val interface{}) string {
	var ss []string
	for _, arg := range cmd.Args() {
		ss = append(ss, fmt.Sprint(arg))
	}
	s := strings.Join(ss, " ")
	if err := cmd.Err(); err != nil {
		return s + ": " + err.Error()
	}
	if val != nil {
		switch vv := val.(type) {
		case []byte:
			return s + ": " + string(vv)
		default:
			return s + ": " + fmt.Sprint(val)
		}
	}
	return s

}

func cmdFirstKeyPos(cmd Cmder, info *CommandInfo) int {
	switch cmd.Name() {
	case "eval", "evalsha":
		if cmd.stringArg(2) != "0" {
			return 3
		}

		return 0
	case "publish":
		return 1
	}
	if info == nil {
		return 0
	}
	return int(info.FirstKeyPos)
}

//------------------------------------------------------------------------------

type baseCmd struct {
	_args []interface{}
	err   error

	_readTimeout *time.Duration
}

var _ Cmder = (*Cmd)(nil)

func (cmd *baseCmd) Err() error {
	return cmd.err
}

func (cmd *baseCmd) Args() []interface{} {
	return cmd._args
}

func (cmd *baseCmd) stringArg(pos int) string {
	if pos < 0 || pos >= len(cmd._args) {
		return ""
	}
	s, _ := cmd._args[pos].(string)
	return s
}

func (cmd *baseCmd) Name() string {
	if len(cmd._args) > 0 {
		// Cmd name must be lower cased.
		s := internal.ToLower(cmd.stringArg(0))
		cmd._args[0] = s
		return s
	}
	return ""
}

func (cmd *baseCmd) readTimeout() *time.Duration {
	return cmd._readTimeout
}

func (cmd *baseCmd) setReadTimeout(d time.Duration) {
	cmd._readTimeout = &d
}

func (cmd *baseCmd) setErr(e error) {
	cmd.err = e
}

//------------------------------------------------------------------------------

type Cmd struct {
	baseCmd

	val interface{}
}

func NewCmd(args ...interface{}) *Cmd {
	return &Cmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *Cmd) Val() interface{} {
	return cmd.val
}

func (cmd *Cmd) Result() (interface{}, error) {
	return cmd.val, cmd.err
}

func (cmd *Cmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *Cmd) readReply(cn *pool.Conn) error {
	cmd.val, cmd.err = cn.Rd.ReadReply(sliceParser)
	if cmd.err != nil {
		return cmd.err
	}
	if b, ok := cmd.val.([]byte); ok {
		// Bytes must be copied, because underlying memory is reused.
		cmd.val = string(b)
	}
	return nil
}

//------------------------------------------------------------------------------

type SliceCmd struct {
	baseCmd

	val []interface{}
}

var _ Cmder = (*SliceCmd)(nil)

func NewSliceCmd(args ...interface{}) *SliceCmd {
	return &SliceCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *SliceCmd) Val() []interface{} {
	return cmd.val
}

func (cmd *SliceCmd) Result() ([]interface{}, error) {
	return cmd.val, cmd.err
}

func (cmd *SliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *SliceCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(sliceParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.([]interface{})
	return nil
}

//------------------------------------------------------------------------------

type StatusCmd struct {
	baseCmd

	val string
}

var _ Cmder = (*StatusCmd)(nil)

func NewStatusCmd(args ...interface{}) *StatusCmd {
	return &StatusCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *StatusCmd) Val() string {
	return cmd.val
}

func (cmd *StatusCmd) Result() (string, error) {
	return cmd.val, cmd.err
}

func (cmd *StatusCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *StatusCmd) readReply(cn *pool.Conn) error {
	cmd.val, cmd.err = cn.Rd.ReadStringReply()
	return cmd.err
}

//------------------------------------------------------------------------------

type IntCmd struct {
	baseCmd

	val int64
}

var _ Cmder = (*IntCmd)(nil)

func NewIntCmd(args ...interface{}) *IntCmd {
	return &IntCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *IntCmd) Val() int64 {
	return cmd.val
}

func (cmd *IntCmd) Result() (int64, error) {
	return cmd.val, cmd.err
}

func (cmd *IntCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *IntCmd) readReply(cn *pool.Conn) error {
	cmd.val, cmd.err = cn.Rd.ReadIntReply()
	return cmd.err
}

//------------------------------------------------------------------------------

type DurationCmd struct {
	baseCmd

	val       time.Duration
	precision time.Duration
}

var _ Cmder = (*DurationCmd)(nil)

func NewDurationCmd(precision time.Duration, args ...interface{}) *DurationCmd {
	return &DurationCmd{
		baseCmd:   baseCmd{_args: args},
		precision: precision,
	}
}

func (cmd *DurationCmd) Val() time.Duration {
	return cmd.val
}

func (cmd *DurationCmd) Result() (time.Duration, error) {
	return cmd.val, cmd.err
}

func (cmd *DurationCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *DurationCmd) readReply(cn *pool.Conn) error {
	var n int64
	n, cmd.err = cn.Rd.ReadIntReply()
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = time.Duration(n) * cmd.precision
	return nil
}

//------------------------------------------------------------------------------

type TimeCmd struct {
	baseCmd

	val time.Time
}

var _ Cmder = (*TimeCmd)(nil)

func NewTimeCmd(args ...interface{}) *TimeCmd {
	return &TimeCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *TimeCmd) Val() time.Time {
	return cmd.val
}

func (cmd *TimeCmd) Result() (time.Time, error) {
	return cmd.val, cmd.err
}

func (cmd *TimeCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *TimeCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(timeParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.(time.Time)
	return nil
}

//------------------------------------------------------------------------------

type BoolCmd struct {
	baseCmd

	val bool
}

var _ Cmder = (*BoolCmd)(nil)

func NewBoolCmd(args ...interface{}) *BoolCmd {
	return &BoolCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *BoolCmd) Val() bool {
	return cmd.val
}

func (cmd *BoolCmd) Result() (bool, error) {
	return cmd.val, cmd.err
}

func (cmd *BoolCmd) String() string {
	return cmdString(cmd, cmd.val)
}

var ok = []byte("OK")

func (cmd *BoolCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadReply(nil)
	// `SET key value NX` returns nil when key already exists. But
	// `SETNX key value` returns bool (0/1). So convert nil to bool.
	// TODO: is this okay?
	if cmd.err == Nil {
		cmd.val = false
		cmd.err = nil
		return nil
	}
	if cmd.err != nil {
		return cmd.err
	}
	switch v := v.(type) {
	case int64:
		cmd.val = v == 1
		return nil
	case []byte:
		cmd.val = bytes.Equal(v, ok)
		return nil
	default:
		cmd.err = fmt.Errorf("got %T, wanted int64 or string", v)
		return cmd.err
	}
}

//------------------------------------------------------------------------------

type StringCmd struct {
	baseCmd

	val []byte
}

var _ Cmder = (*StringCmd)(nil)

func NewStringCmd(args ...interface{}) *StringCmd {
	return &StringCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *StringCmd) Val() string {
	return util.BytesToString(cmd.val)
}

func (cmd *StringCmd) Result() (string, error) {
	return cmd.Val(), cmd.err
}

func (cmd *StringCmd) Bytes() ([]byte, error) {
	return cmd.val, cmd.err
}

func (cmd *StringCmd) Int64() (int64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseInt(cmd.Val(), 10, 64)
}

func (cmd *StringCmd) Uint64() (uint64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseUint(cmd.Val(), 10, 64)
}

func (cmd *StringCmd) Float64() (float64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	return strconv.ParseFloat(cmd.Val(), 64)
}

func (cmd *StringCmd) Scan(val interface{}) error {
	if cmd.err != nil {
		return cmd.err
	}
	return proto.Scan(cmd.val, val)
}

func (cmd *StringCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *StringCmd) readReply(cn *pool.Conn) error {
	cmd.val, cmd.err = cn.Rd.ReadBytesReply()
	return cmd.err
}

//------------------------------------------------------------------------------

type FloatCmd struct {
	baseCmd

	val float64
}

var _ Cmder = (*FloatCmd)(nil)

func NewFloatCmd(args ...interface{}) *FloatCmd {
	return &FloatCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *FloatCmd) Val() float64 {
	return cmd.val
}

func (cmd *FloatCmd) Result() (float64, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *FloatCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *FloatCmd) readReply(cn *pool.Conn) error {
	cmd.val, cmd.err = cn.Rd.ReadFloatReply()
	return cmd.err
}

//------------------------------------------------------------------------------

type StringSliceCmd struct {
	baseCmd

	val []string
}

var _ Cmder = (*StringSliceCmd)(nil)

func NewStringSliceCmd(args ...interface{}) *StringSliceCmd {
	return &StringSliceCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *StringSliceCmd) Val() []string {
	return cmd.val
}

func (cmd *StringSliceCmd) Result() ([]string, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *StringSliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *StringSliceCmd) ScanSlice(container interface{}) error {
	return proto.ScanSlice(cmd.Val(), container)
}

func (cmd *StringSliceCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(stringSliceParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.([]string)
	return nil
}

//------------------------------------------------------------------------------

type BoolSliceCmd struct {
	baseCmd

	val []bool
}

var _ Cmder = (*BoolSliceCmd)(nil)

func NewBoolSliceCmd(args ...interface{}) *BoolSliceCmd {
	return &BoolSliceCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *BoolSliceCmd) Val() []bool {
	return cmd.val
}

func (cmd *BoolSliceCmd) Result() ([]bool, error) {
	return cmd.val, cmd.err
}

func (cmd *BoolSliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *BoolSliceCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(boolSliceParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.([]bool)
	return nil
}

//------------------------------------------------------------------------------

type StringStringMapCmd struct {
	baseCmd

	val map[string]string
}

var _ Cmder = (*StringStringMapCmd)(nil)

func NewStringStringMapCmd(args ...interface{}) *StringStringMapCmd {
	return &StringStringMapCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *StringStringMapCmd) Val() map[string]string {
	return cmd.val
}

func (cmd *StringStringMapCmd) Result() (map[string]string, error) {
	return cmd.val, cmd.err
}

func (cmd *StringStringMapCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *StringStringMapCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(stringStringMapParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.(map[string]string)
	return nil
}

//------------------------------------------------------------------------------

type StringIntMapCmd struct {
	baseCmd

	val map[string]int64
}

var _ Cmder = (*StringIntMapCmd)(nil)

func NewStringIntMapCmd(args ...interface{}) *StringIntMapCmd {
	return &StringIntMapCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *StringIntMapCmd) Val() map[string]int64 {
	return cmd.val
}

func (cmd *StringIntMapCmd) Result() (map[string]int64, error) {
	return cmd.val, cmd.err
}

func (cmd *StringIntMapCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *StringIntMapCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(stringIntMapParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.(map[string]int64)
	return nil
}

//------------------------------------------------------------------------------

type StringStructMapCmd struct {
	baseCmd

	val map[string]struct{}
}

var _ Cmder = (*StringStructMapCmd)(nil)

func NewStringStructMapCmd(args ...interface{}) *StringStructMapCmd {
	return &StringStructMapCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *StringStructMapCmd) Val() map[string]struct{} {
	return cmd.val
}

func (cmd *StringStructMapCmd) Result() (map[string]struct{}, error) {
	return cmd.val, cmd.err
}

func (cmd *StringStructMapCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *StringStructMapCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(stringStructMapParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.(map[string]struct{})
	return nil
}

//------------------------------------------------------------------------------

type XStream struct {
	Stream   string
	Messages []*XMessage
}

type XMessage struct {
	ID     string
	Values map[string]interface{}
}

//------------------------------------------------------------------------------

type XStreamSliceCmd struct {
	baseCmd

	val []*XStream
}

var _ Cmder = (*XStreamSliceCmd)(nil)

func NewXStreamSliceCmd(args ...interface{}) *XStreamSliceCmd {
	return &XStreamSliceCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *XStreamSliceCmd) Val() []*XStream {
	return cmd.val
}

func (cmd *XStreamSliceCmd) Result() ([]*XStream, error) {
	return cmd.val, cmd.err
}

func (cmd *XStreamSliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *XStreamSliceCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(xStreamSliceParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.([]*XStream)
	return nil
}

// Implements proto.MultiBulkParse
func xStreamSliceParser(rd *proto.Reader, n int64) (interface{}, error) {
	xx := make([]*XStream, n)
	for i := int64(0); i < n; i++ {
		v, err := rd.ReadArrayReply(xStreamParser)
		if err != nil {
			return nil, err
		}
		xx[i] = v.(*XStream)
	}
	return xx, nil
}

// Implements proto.MultiBulkParse
func xStreamParser(rd *proto.Reader, n int64) (interface{}, error) {
	if n != 2 {
		return nil, fmt.Errorf("got %d, wanted 2", n)
	}

	stream, err := rd.ReadStringReply()
	if err != nil {
		return nil, err
	}

	v, err := rd.ReadArrayReply(xMessageSliceParser)
	if err != nil {
		return nil, err
	}

	return &XStream{
		Stream:   stream,
		Messages: v.([]*XMessage),
	}, nil
}

//------------------------------------------------------------------------------

type XMessageSliceCmd struct {
	baseCmd

	val []*XMessage
}

var _ Cmder = (*XMessageSliceCmd)(nil)

func NewXMessageSliceCmd(args ...interface{}) *XMessageSliceCmd {
	return &XMessageSliceCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *XMessageSliceCmd) Val() []*XMessage {
	return cmd.val
}

func (cmd *XMessageSliceCmd) Result() ([]*XMessage, error) {
	return cmd.val, cmd.err
}

func (cmd *XMessageSliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *XMessageSliceCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(xMessageSliceParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.([]*XMessage)
	return nil
}

// Implements proto.MultiBulkParse
func xMessageSliceParser(rd *proto.Reader, n int64) (interface{}, error) {
	msgs := make([]*XMessage, n)
	for i := int64(0); i < n; i++ {
		v, err := rd.ReadArrayReply(xMessageParser)
		if err != nil {
			return nil, err
		}
		msgs[i] = v.(*XMessage)
	}
	return msgs, nil
}

// Implements proto.MultiBulkParse
func xMessageParser(rd *proto.Reader, n int64) (interface{}, error) {
	id, err := rd.ReadStringReply()
	if err != nil {
		return nil, err
	}

	v, err := rd.ReadArrayReply(xKeyValueParser)
	if err != nil {
		return nil, err
	}

	return &XMessage{
		ID:     id,
		Values: v.(map[string]interface{}),
	}, nil
}

// Implements proto.MultiBulkParse
func xKeyValueParser(rd *proto.Reader, n int64) (interface{}, error) {
	values := make(map[string]interface{}, n)
	for i := int64(0); i < n; i += 2 {
		key, err := rd.ReadStringReply()
		if err != nil {
			return nil, err
		}

		value, err := rd.ReadStringReply()
		if err != nil {
			return nil, err
		}

		values[key] = value
	}
	return values, nil
}

//------------------------------------------------------------------------------

type ZSliceCmd struct {
	baseCmd

	val []Z
}

var _ Cmder = (*ZSliceCmd)(nil)

func NewZSliceCmd(args ...interface{}) *ZSliceCmd {
	return &ZSliceCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *ZSliceCmd) Val() []Z {
	return cmd.val
}

func (cmd *ZSliceCmd) Result() ([]Z, error) {
	return cmd.val, cmd.err
}

func (cmd *ZSliceCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *ZSliceCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(zSliceParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.([]Z)
	return nil
}

//------------------------------------------------------------------------------

type ScanCmd struct {
	baseCmd

	page   []string
	cursor uint64

	process func(cmd Cmder) error
}

var _ Cmder = (*ScanCmd)(nil)

func NewScanCmd(process func(cmd Cmder) error, args ...interface{}) *ScanCmd {
	return &ScanCmd{
		baseCmd: baseCmd{_args: args},
		process: process,
	}
}

func (cmd *ScanCmd) Val() (keys []string, cursor uint64) {
	return cmd.page, cmd.cursor
}

func (cmd *ScanCmd) Result() (keys []string, cursor uint64, err error) {
	return cmd.page, cmd.cursor, cmd.err
}

func (cmd *ScanCmd) String() string {
	return cmdString(cmd, cmd.page)
}

func (cmd *ScanCmd) readReply(cn *pool.Conn) error {
	cmd.page, cmd.cursor, cmd.err = cn.Rd.ReadScanReply()
	return cmd.err
}

// Iterator creates a new ScanIterator.
func (cmd *ScanCmd) Iterator() *ScanIterator {
	return &ScanIterator{
		cmd: cmd,
	}
}

//------------------------------------------------------------------------------

type ClusterNode struct {
	Id   string
	Addr string
}

type ClusterSlot struct {
	Start int
	End   int
	Nodes []ClusterNode
}

type ClusterSlotsCmd struct {
	baseCmd

	val []ClusterSlot
}

var _ Cmder = (*ClusterSlotsCmd)(nil)

func NewClusterSlotsCmd(args ...interface{}) *ClusterSlotsCmd {
	return &ClusterSlotsCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *ClusterSlotsCmd) Val() []ClusterSlot {
	return cmd.val
}

func (cmd *ClusterSlotsCmd) Result() ([]ClusterSlot, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *ClusterSlotsCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *ClusterSlotsCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(clusterSlotsParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.([]ClusterSlot)
	return nil
}

//------------------------------------------------------------------------------

// GeoLocation is used with GeoAdd to add geospatial location.
type GeoLocation struct {
	Name                      string
	Longitude, Latitude, Dist float64
	GeoHash                   int64
}

// GeoRadiusQuery is used with GeoRadius to query geospatial index.
type GeoRadiusQuery struct {
	Radius float64
	// Can be m, km, ft, or mi. Default is km.
	Unit        string
	WithCoord   bool
	WithDist    bool
	WithGeoHash bool
	Count       int
	// Can be ASC or DESC. Default is no sort order.
	Sort      string
	Store     string
	StoreDist string
}

type GeoLocationCmd struct {
	baseCmd

	q         *GeoRadiusQuery
	locations []GeoLocation
}

var _ Cmder = (*GeoLocationCmd)(nil)

func NewGeoLocationCmd(q *GeoRadiusQuery, args ...interface{}) *GeoLocationCmd {
	args = append(args, q.Radius)
	if q.Unit != "" {
		args = append(args, q.Unit)
	} else {
		args = append(args, "km")
	}
	if q.WithCoord {
		args = append(args, "withcoord")
	}
	if q.WithDist {
		args = append(args, "withdist")
	}
	if q.WithGeoHash {
		args = append(args, "withhash")
	}
	if q.Count > 0 {
		args = append(args, "count", q.Count)
	}
	if q.Sort != "" {
		args = append(args, q.Sort)
	}
	if q.Store != "" {
		args = append(args, "store")
		args = append(args, q.Store)
	}
	if q.StoreDist != "" {
		args = append(args, "storedist")
		args = append(args, q.StoreDist)
	}
	return &GeoLocationCmd{
		baseCmd: baseCmd{_args: args},
		q:       q,
	}
}

func (cmd *GeoLocationCmd) Val() []GeoLocation {
	return cmd.locations
}

func (cmd *GeoLocationCmd) Result() ([]GeoLocation, error) {
	return cmd.locations, cmd.err
}

func (cmd *GeoLocationCmd) String() string {
	return cmdString(cmd, cmd.locations)
}

func (cmd *GeoLocationCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(newGeoLocationSliceParser(cmd.q))
	if cmd.err != nil {
		return cmd.err
	}
	cmd.locations = v.([]GeoLocation)
	return nil
}

//------------------------------------------------------------------------------

type GeoPos struct {
	Longitude, Latitude float64
}

type GeoPosCmd struct {
	baseCmd

	positions []*GeoPos
}

var _ Cmder = (*GeoPosCmd)(nil)

func NewGeoPosCmd(args ...interface{}) *GeoPosCmd {
	return &GeoPosCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *GeoPosCmd) Val() []*GeoPos {
	return cmd.positions
}

func (cmd *GeoPosCmd) Result() ([]*GeoPos, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *GeoPosCmd) String() string {
	return cmdString(cmd, cmd.positions)
}

func (cmd *GeoPosCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(geoPosSliceParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.positions = v.([]*GeoPos)
	return nil
}

//------------------------------------------------------------------------------

type CommandInfo struct {
	Name        string
	Arity       int8
	Flags       []string
	FirstKeyPos int8
	LastKeyPos  int8
	StepCount   int8
	ReadOnly    bool
}

type CommandsInfoCmd struct {
	baseCmd

	val map[string]*CommandInfo
}

var _ Cmder = (*CommandsInfoCmd)(nil)

func NewCommandsInfoCmd(args ...interface{}) *CommandsInfoCmd {
	return &CommandsInfoCmd{
		baseCmd: baseCmd{_args: args},
	}
}

func (cmd *CommandsInfoCmd) Val() map[string]*CommandInfo {
	return cmd.val
}

func (cmd *CommandsInfoCmd) Result() (map[string]*CommandInfo, error) {
	return cmd.Val(), cmd.Err()
}

func (cmd *CommandsInfoCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *CommandsInfoCmd) readReply(cn *pool.Conn) error {
	var v interface{}
	v, cmd.err = cn.Rd.ReadArrayReply(commandInfoSliceParser)
	if cmd.err != nil {
		return cmd.err
	}
	cmd.val = v.(map[string]*CommandInfo)
	return nil
}

//------------------------------------------------------------------------------

type cmdsInfoCache struct {
	fn func() (map[string]*CommandInfo, error)

	once internal.Once
	cmds map[string]*CommandInfo
}

func newCmdsInfoCache(fn func() (map[string]*CommandInfo, error)) *cmdsInfoCache {
	return &cmdsInfoCache{
		fn: fn,
	}
}

func (c *cmdsInfoCache) Get() (map[string]*CommandInfo, error) {
	err := c.once.Do(func() error {
		cmds, err := c.fn()
		if err != nil {
			return err
		}
		c.cmds = cmds
		return nil
	})
	return c.cmds, err
}
