package error

import (
	"fmt"
	"github.com/beego/beego/v2/core/codes"
	"strconv"
)

// Error type defines custom error for Beego. It is used by every module
// in Beego. Each `Error` message contains three pieces of data: error code,
// error message.
// More docs http://beego.me/docs/module/error.md.
type Error struct {
    Code codes.Code
    Msg string
}

// New returns a Error representing c and msg.
func New(c codes.Code, msg string) *Error {
    return &Error{Code: c, Msg: msg}
}

// Err returns an error representing c and msg.  If c is OK, returns nil.
func Err(c codes.Code, msg string) error {
	return New(c, msg)
}

// Errorf returns Error(c, fmt.Sprintf(format, a...)).
func Errorf(c codes.Code, format string, a ...interface{}) error {
	return Err(c, fmt.Sprintf(format, a...))
}

// Error returns formatted message for user.
func (e *Error) Error() string {
	codeSrt := strconv.FormatUint(uint64(e.GetCode()), 10)
	return fmt.Sprintf("beego error: code = %s desc = %s", codeSrt, e.GetMessage())
}

// GetCode returns Error's Code.
func (e *Error) GetCode() codes.Code {
	if e != nil {
		return e.Code
	}
	return 0
}

// GetMessage returns Error's Msg.
func (e *Error) GetMessage() string {
	if e != nil {
		return e.Msg
	}
	return ""
}