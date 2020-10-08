// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"io"

	"github.com/astaxie/beego/core/utils"
)

// Email is the type used for email messages
type Email utils.Email

// Attachment is a struct representing an email attachment.
// Based on the mime/multipart.FileHeader struct, Attachment contains the name, MIMEHeader, and content of the attachment in question
type Attachment utils.Attachment

// NewEMail create new Email struct with config json.
// config json is followed from Email struct fields.
func NewEMail(config string) *Email {
	return (*Email)(utils.NewEMail(config))
}

// Bytes Make all send information to byte
func (e *Email) Bytes() ([]byte, error) {
	return (*utils.Email)(e).Bytes()
}

// AttachFile Add attach file to the send mail
func (e *Email) AttachFile(args ...string) (*Attachment, error) {
	a, err := (*utils.Email)(e).AttachFile(args...)
	if err != nil {
		return nil, err
	}
	return (*Attachment)(a), err
}

// Attach is used to attach content from an io.Reader to the email.
// Parameters include an io.Reader, the desired filename for the attachment, and the Content-Type.
func (e *Email) Attach(r io.Reader, filename string, args ...string) (*Attachment, error) {
	a, err := (*utils.Email)(e).Attach(r, filename, args...)
	if err != nil {
		return nil, err
	}
	return (*Attachment)(a), err
}

// Send will send out the mail
func (e *Email) Send() error {
	return (*utils.Email)(e).Send()
}
