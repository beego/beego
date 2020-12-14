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

package orm

import (
	"github.com/beego/beego/v2/client/orm"
)

type baseQuerySetter struct {
}

func (b *baseQuerySetter) ForceIndex(indexes ...string) orm.QuerySeter {
	panic("you should not invoke this method.")
}

func (b *baseQuerySetter) UseIndex(indexes ...string) orm.QuerySeter {
	panic("you should not invoke this method.")
}

func (b *baseQuerySetter) IgnoreIndex(indexes ...string) orm.QuerySeter {
	panic("you should not invoke this method.")
}
