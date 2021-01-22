package error

import (
	"github.com/beego/beego/v2/core/codes"
	"reflect"
	"testing"
)

func TestErr(t *testing.T) {
	type args struct {
		c   codes.Code
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "1", args: args{codes.SessionSessionStartError, codes.CodeToStr[codes.SessionSessionStartError]}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Err(tt.args.c, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Err() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		Code codes.Code
		Msg  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{name: "1", fields: fields{codes.SessionSessionStartError, codes.CodeToStr[codes.SessionSessionStartError]}, want: "beego error: code = 5001001 desc = \"SESSION_MODULE_SESSION_START_ERROR\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Code: tt.fields.Code,
				Msg:  tt.fields.Msg,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_GetCode(t *testing.T) {
	type fields struct {
		Code codes.Code
		Msg  string
	}
	tests := []struct {
		name   string
		fields fields
		want   codes.Code
	}{
		// TODO: Add test cases.
		{name: "1", fields: fields{codes.SessionSessionStartError, codes.CodeToStr[codes.SessionSessionStartError]}, want: codes.SessionSessionStartError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Code: tt.fields.Code,
				Msg:  tt.fields.Msg,
			}
			if got := e.GetCode(); got != tt.want {
				t.Errorf("GetCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_GetMessage(t *testing.T) {
	type fields struct {
		Code codes.Code
		Msg  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{name: "1", fields: fields{codes.SessionSessionStartError, codes.CodeToStr[codes.SessionSessionStartError]}, want: codes.CodeToStr[codes.SessionSessionStartError]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Code: tt.fields.Code,
				Msg:  tt.fields.Msg,
			}
			if got := e.GetMessage(); got != tt.want {
				t.Errorf("GetMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		c      codes.Code
		format string
		a      []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "1", args: args{codes.SessionSessionStartError, "%s", []interface{}{codes.CodeToStr[codes.SessionSessionStartError]}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Errorf(tt.args.c, tt.args.format, tt.args.a...); (err != nil) != tt.wantErr {
				t.Errorf("Errorf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		c   codes.Code
		msg string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
		{name: "1", args: args{codes.SessionSessionStartError, codes.CodeToStr[codes.SessionSessionStartError]}, want: &Error{Code:codes.SessionSessionStartError,  Msg:codes.CodeToStr[codes.SessionSessionStartError]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.c, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
