// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// And Modify By Wendal (wendal1985@gmail.com)

// Package errors implements functions to manipulate errors.
package errors

import (
	"runtime/debug"
)

var AddStack = true

// New returns an error that formats as the given text.
func New(text string) error {
	if AddStack {
		text += "\n" + string(debug.Stack())
	}
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
