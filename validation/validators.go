package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"time"
)

type Validator interface {
	IsSatisfied(interface{}) bool
	DefaultMessage() string
	GetKey() string
}

type Required struct {
	Key string
}

func (r Required) IsSatisfied(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if str, ok := obj.(string); ok {
		return len(str) > 0
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	if i, ok := obj.(int); ok {
		return i != 0
	}
	if t, ok := obj.(time.Time); ok {
		return !t.IsZero()
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() > 0
	}
	return true
}

func (r Required) DefaultMessage() string {
	return "Required"
}

func (r Required) GetKey() string {
	return r.Key
}

type Min struct {
	Min int
	Key string
}

func (m Min) IsSatisfied(obj interface{}) bool {
	num, ok := obj.(int)
	if ok {
		return num >= m.Min
	}
	return false
}

func (m Min) DefaultMessage() string {
	return fmt.Sprint("Minimum is ", m.Min)
}

func (m Min) GetKey() string {
	return m.Key
}

type Max struct {
	Max int
	Key string
}

func (m Max) IsSatisfied(obj interface{}) bool {
	num, ok := obj.(int)
	if ok {
		return num <= m.Max
	}
	return false
}

func (m Max) DefaultMessage() string {
	return fmt.Sprint("Maximum is ", m.Max)
}

func (m Max) GetKey() string {
	return m.Key
}

// Requires an integer to be within Min, Max inclusive.
type Range struct {
	Min
	Max
	Key string
}

func (r Range) IsSatisfied(obj interface{}) bool {
	return r.Min.IsSatisfied(obj) && r.Max.IsSatisfied(obj)
}

func (r Range) DefaultMessage() string {
	return fmt.Sprint("Range is ", r.Min.Min, " to ", r.Max.Max)
}

func (r Range) GetKey() string {
	return r.Key
}

// Requires an array or string to be at least a given length.
type MinSize struct {
	Min int
	Key string
}

func (m MinSize) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return len(str) >= m.Min
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() >= m.Min
	}
	return false
}

func (m MinSize) DefaultMessage() string {
	return fmt.Sprint("Minimum size is ", m.Min)
}

func (m MinSize) GetKey() string {
	return m.Key
}

// Requires an array or string to be at most a given length.
type MaxSize struct {
	Max int
	Key string
}

func (m MaxSize) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return len(str) <= m.Max
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() <= m.Max
	}
	return false
}

func (m MaxSize) DefaultMessage() string {
	return fmt.Sprint("Maximum size is ", m.Max)
}

func (m MaxSize) GetKey() string {
	return m.Key
}

// Requires an array or string to be exactly a given length.
type Length struct {
	N   int
	Key string
}

func (l Length) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return len(str) == l.N
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() == l.N
	}
	return false
}

func (l Length) DefaultMessage() string {
	return fmt.Sprint("Required length is ", l.N)
}

func (l Length) GetKey() string {
	return l.Key
}

type Alpha struct {
	Key string
}

func (a Alpha) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') {
				return false
			}
		}
		return true
	}
	return false
}

func (a Alpha) DefaultMessage() string {
	return fmt.Sprint("Must be valid alpha characters")
}

func (a Alpha) GetKey() string {
	return a.Key
}

type Numeric struct {
	Key string
}

func (n Numeric) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if '9' < v || v < '0' {
				return false
			}
		}
		return true
	}
	return false
}

func (n Numeric) DefaultMessage() string {
	return fmt.Sprint("Must be valid numeric characters")
}

func (n Numeric) GetKey() string {
	return n.Key
}

type AlphaNumeric struct {
	Key string
}

func (a AlphaNumeric) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
				return false
			}
		}
		return true
	}
	return false
}

func (a AlphaNumeric) DefaultMessage() string {
	return fmt.Sprint("Must be valid alpha or numeric characters")
}

func (a AlphaNumeric) GetKey() string {
	return a.Key
}

// Requires a string to match a given regex.
type Match struct {
	Regexp *regexp.Regexp
	Key    string
}

func (m Match) IsSatisfied(obj interface{}) bool {
	str := obj.(string)
	return m.Regexp.MatchString(str)
}

func (m Match) DefaultMessage() string {
	return fmt.Sprint("Must match ", m.Regexp)
}

func (m Match) GetKey() string {
	return m.Key
}

// Requires a string to not match a given regex.
type NoMatch struct {
	Match
	Key string
}

func (n NoMatch) IsSatisfied(obj interface{}) bool {
	return !n.Match.IsSatisfied(obj)
}

func (n NoMatch) DefaultMessage() string {
	return fmt.Sprint("Must not match ", n.Regexp)
}

func (n NoMatch) GetKey() string {
	return n.Key
}

var alphaDashPattern = regexp.MustCompile("[^\\d\\w-_]")

type AlphaDash struct {
	NoMatch
	Key string
}

func (a AlphaDash) DefaultMessage() string {
	return fmt.Sprint("Must be valid alpha or numeric or dash(-_) characters")
}

func (a AlphaDash) GetKey() string {
	return a.Key
}

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

type Email struct {
	Match
	Key string
}

func (e Email) DefaultMessage() string {
	return fmt.Sprint("Must be a valid email address")
}

func (e Email) GetKey() string {
	return e.Key
}

var ipPattern = regexp.MustCompile("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$")

type IP struct {
	Match
	Key string
}

func (i IP) DefaultMessage() string {
	return fmt.Sprint("Must be a valid ip address")
}

func (i IP) GetKey() string {
	return i.Key
}

var base64Pattern = regexp.MustCompile("^(?:[A-Za-z0-99+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$")

type Base64 struct {
	Match
	Key string
}

func (b Base64) DefaultMessage() string {
	return fmt.Sprint("Must be valid base64 characters")
}

func (b Base64) GetKey() string {
	return b.Key
}
