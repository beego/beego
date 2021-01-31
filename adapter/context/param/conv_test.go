// Copyright 2020 beego
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

package param

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/adapter/context"
)


// Demo is used to test, it's empty
func Demo(i int)  {

}

func TestConvertParams(t *testing.T) {
	res := ConvertParams(nil, reflect.TypeOf(Demo), context.NewContext())
	assert.Equal(t, 0, len(res))
	ctx := context.NewContext()
	ctx.Input.RequestBody = []byte("11")
	res = ConvertParams([]*MethodParam{
		New("A", InBody),
	}, reflect.TypeOf(Demo), ctx)
	assert.Equal(t, int64(11), res[0].Int())
}

