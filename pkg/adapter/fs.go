// Copyright 2020
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adapter

import (
	"net/http"
	"path/filepath"

	"github.com/astaxie/beego/pkg/server/web"
)

type FileSystem web.FileSystem

func (d FileSystem) Open(name string) (http.File, error) {
	return (web.FileSystem)(d).Open(name)
}

// Walk walks the file tree rooted at root in filesystem, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn.
func Walk(fs http.FileSystem, root string, walkFn filepath.WalkFunc) error {
	return web.Walk(fs, root, walkFn)
}
