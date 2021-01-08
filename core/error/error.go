package error

import (
	"fmt"
	"github.com/beego/beego/v2/core/codes"
)

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

func (e *Error) Error() string {
	return fmt.Sprintf("beego error: code = %s desc = %s", e.GetCode(), e.GetMessage())
}

func (x *Error) GetCode() codes.Code {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Msg
	}
	return ""
}