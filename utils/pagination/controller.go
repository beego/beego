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

package pagination

import (
	"github.com/astaxie/beego/context"
)

// Instantiates a Paginator and assigns it to context.Input.Data["paginator"].
func SetPaginator(context *context.Context, per int, nums int64) (paginator *Paginator) {
	paginator = NewPaginator(context.Request, per, nums)
	context.Input.Data["paginator"] = &paginator
	return
}
