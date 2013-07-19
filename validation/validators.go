package revel

import (
	"fmt"
	"reflect"
	"regexp"
	"time"
)

type Validator interface {
	IsSatisfied(interface{}) bool
	DefaultMessage() string
}

type Required struct{}

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

// TODO
// type Unique struct {}

type Min struct {
	Min int
}

func (m Min) IsSatisfied(obj interface{}) bool {
	num, ok := obj.(int)
	if ok {
		return num >= m.Min
	}
	return false
}

func (m Min) DefaultMessage() string {
	return fmt.Sprintln("Minimum is", m.Min)
}

type Max struct {
	Max int
}

func (m Max) IsSatisfied(obj interface{}) bool {
	num, ok := obj.(int)
	if ok {
		return num <= m.Max
	}
	return false
}

func (m Max) DefaultMessage() string {
	return fmt.Sprintln("Maximum is", m.Max)
}

// Requires an integer to be within Min, Max inclusive.
type Range struct {
	Min
	Max
}

func (r Range) IsSatisfied(obj interface{}) bool {
	return r.Min.IsSatisfied(obj) && r.Max.IsSatisfied(obj)
}

func (r Range) DefaultMessage() string {
	return fmt.Sprintln("Range is", r.Min.Min, "to", r.Max.Max)
}

// Requires an array or string to be at least a given length.
type MinSize struct {
	Min int
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
	return fmt.Sprintln("Minimum size is", m.Min)
}

// Requires an array or string to be at most a given length.
type MaxSize struct {
	Max int
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
	return fmt.Sprintln("Maximum size is", m.Max)
}

// Requires an array or string to be exactly a given length.
type Length struct {
	N int
}

func (s Length) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return len(str) == s.N
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() == s.N
	}
	return false
}

func (s Length) DefaultMessage() string {
	return fmt.Sprintln("Required length is", s.N)
}

type Alpha struct{}

func (a Alpha) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') {
				return false
			}
		}
	}
	return false
}

func (a Alpha) DefaultMessage() string {
	return fmt.Sprintln("Must be valid alpha characters")
}

type Numeric struct{}

func (n Numeric) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if '9' < v || v < '0' {
				return false
			}
		}
	}
	return false
}

func (n Numeric) DefaultMessage() string {
	return fmt.Sprintln("Must be valid numeric characters")
}

type AlphaNumeric struct{}

func (a AlphaNumeric) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
				return false
			}
		}
	}
	return false
}

func (a AlphaNumeric) DefaultMessage() string {
	return fmt.Sprintln("Must be valid alpha or numeric characters")
}

// Requires a string to match a given regex.
type Match struct {
	Regexp *regexp.Regexp
}

func (m Match) IsSatisfied(obj interface{}) bool {
	str := obj.(string)
	return m.Regexp.MatchString(str)
}

func (m Match) DefaultMessage() string {
	return fmt.Sprintln("Must match", m.Regexp)
}

// Requires a string to not match a given regex.
type NoMatch struct {
	Match
}

func (m NoMatch) IsSatisfied(obj interface{}) bool {
	return !m.Match.IsSatisfied(obj)
}

func (m NoMatch) DefaultMessage() string {
	return fmt.Sprintln("Must no match", m.Regexp)
}

var alphaDashPattern = regexp.MustCompile("[^\\d\\w-_]")

type AlphaDash struct {
	NoMatch
}

func (a AlphaDash) DefaultMessage() string {
	return fmt.Sprintln("Must be valid characters")
}

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

type Email struct {
	Match
}

func (e Email) DefaultMessage() string {
	return fmt.Sprintln("Must be a valid email address")
}

var ipPattern = regexp.MustCompile("((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)")

type IP struct {
	Match
}

func (i IP) DefaultMessage() string {
	return fmt.Sprintln("Must be a valid ip address")
}

var base64Pattern = regexp.MustCompile("^(?:[A-Za-z0-99+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$")

type Base64 struct {
	Match
}

func (b Base64) DefaultMessage() string {
	return fmt.Sprintln("Must be valid base64 characters")
}
