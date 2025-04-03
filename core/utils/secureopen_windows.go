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

//go:build windows
// +build windows

package utils

import (
	"os"
	"syscall"
)

// OpenFileSecure Secure file opening function
// Check whether it is a symbolic connection in Windows
func OpenFileSecure(name string, flag int, perm os.FileMode) (*os.File, error) {
	// Check if it's a symbolic link
	if fi, err := os.Lstat(name); err == nil {
		if isSymlink(fi) {
			return nil, &os.PathError{
				Op:   "open",
				Path: name,
				Err:  syscall.Errno(0x5B4), // ERROR_CANT_ACCESS_FILE
			}
		}
	}
	// Open the file in the normal way
	return os.OpenFile(name, flag, perm)
}

// isSymlink Check if the file is a symbolic link
func isSymlink(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeSymlink != 0 ||
		fi.Mode()&os.ModeDevice != 0 &&
			fi.Mode()&os.ModeCharDevice == 0
}
