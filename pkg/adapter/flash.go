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

package adapter

import (
	"github.com/astaxie/beego/pkg/server/web"
)

// FlashData is a tools to maintain data when using across request.
type FlashData web.FlashData

// NewFlash return a new empty FlashData struct.
func NewFlash() *FlashData {
	return (*FlashData)(web.NewFlash())
}

// Set message to flash
func (fd *FlashData) Set(key string, msg string, args ...interface{}) {
	(*web.FlashData)(fd).Set(key, msg, args...)
}

// Success writes success message to flash.
func (fd *FlashData) Success(msg string, args ...interface{}) {
	(*web.FlashData)(fd).Success(msg, args...)
}

// Notice writes notice message to flash.
func (fd *FlashData) Notice(msg string, args ...interface{}) {
	(*web.FlashData)(fd).Notice(msg, args...)
}

// Warning writes warning message to flash.
func (fd *FlashData) Warning(msg string, args ...interface{}) {
	(*web.FlashData)(fd).Warning(msg, args...)
}

// Error writes error message to flash.
func (fd *FlashData) Error(msg string, args ...interface{}) {
	(*web.FlashData)(fd).Error(msg, args...)
}

// Store does the saving operation of flash data.
// the data are encoded and saved in cookie.
func (fd *FlashData) Store(c *Controller) {
	(*web.FlashData)(fd).Store((*web.Controller)(c))
}

// ReadFromRequest parsed flash data from encoded values in cookie.
func ReadFromRequest(c *Controller) *FlashData {
	return (*FlashData)(web.ReadFromRequest((*web.Controller)(c)))
}
