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

package bean

import (
	"context"
	"reflect"
	"strings"

	"github.com/astaxie/beego/core/logs"

	"github.com/astaxie/beego/client/orm"
	"github.com/astaxie/beego/core/bean"
)

// DefaultValueFilterChainBuilder only works for InsertXXX method,
// But InsertOrUpdate and InsertOrUpdateWithCtx is more dangerous than other methods.
// so we won't handle those two methods unless you set includeInsertOrUpdate to true
// And if the element is not pointer, this filter doesn't work
type DefaultValueFilterChainBuilder struct {
	factory                bean.AutoWireBeanFactory
	compatibleWithOldStyle bool

	// only the includeInsertOrUpdate is true, this filter will handle those two methods
	includeInsertOrUpdate bool
}

// NewDefaultValueFilterChainBuilder will create an instance of DefaultValueFilterChainBuilder
// In beego v1.x, the default value config looks like orm:default(xxxx)
// But the default value in 2.x is default:xxx
// so if you want to be compatible with v1.x, please pass true as compatibleWithOldStyle
func NewDefaultValueFilterChainBuilder(typeAdapters map[string]bean.TypeAdapter,
	includeInsertOrUpdate bool,
	compatibleWithOldStyle bool) *DefaultValueFilterChainBuilder {
	factory := bean.NewTagAutoWireBeanFactory()

	if compatibleWithOldStyle {
		newParser := factory.FieldTagParser
		factory.FieldTagParser = func(field reflect.StructField) *bean.FieldMetadata {
			if newParser != nil && field.Tag.Get(bean.DefaultValueTagKey) != "" {
				return newParser(field)
			} else {
				res := &bean.FieldMetadata{}
				ormMeta := field.Tag.Get("orm")
				ormMetaParts := strings.Split(ormMeta, ";")
				for _, p := range ormMetaParts {
					if strings.HasPrefix(p, "default(") && strings.HasSuffix(p, ")") {
						res.DftValue = p[8 : len(p)-1]
					}
				}
				return res
			}
		}
	}

	for k, v := range typeAdapters {
		factory.Adapters[k] = v
	}

	return &DefaultValueFilterChainBuilder{
		factory:                factory,
		compatibleWithOldStyle: compatibleWithOldStyle,
		includeInsertOrUpdate:  includeInsertOrUpdate,
	}
}

func (d *DefaultValueFilterChainBuilder) FilterChain(next orm.Filter) orm.Filter {
	return func(ctx context.Context, inv *orm.Invocation) []interface{} {
		switch inv.Method {
		case "Insert", "InsertWithCtx":
			d.handleInsert(ctx, inv)
			break
		case "InsertOrUpdate", "InsertOrUpdateWithCtx":
			d.handleInsertOrUpdate(ctx, inv)
			break
		case "InsertMulti", "InsertMultiWithCtx":
			d.handleInsertMulti(ctx, inv)
			break
		}
		return next(ctx, inv)
	}
}

func (d *DefaultValueFilterChainBuilder) handleInsert(ctx context.Context, inv *orm.Invocation) {
	d.setDefaultValue(ctx, inv.Args[0])
}

func (d *DefaultValueFilterChainBuilder) handleInsertOrUpdate(ctx context.Context, inv *orm.Invocation) {
	if d.includeInsertOrUpdate {
		ins := inv.Args[0]
		if ins == nil {
			return
		}

		pkName := inv.GetPkFieldName()
		pkField := reflect.Indirect(reflect.ValueOf(ins)).FieldByName(pkName)

		if pkField.IsZero() {
			d.setDefaultValue(ctx, ins)
		}
	}
}

func (d *DefaultValueFilterChainBuilder) handleInsertMulti(ctx context.Context, inv *orm.Invocation) {
	mds := inv.Args[1]

	if t := reflect.TypeOf(mds).Kind(); t != reflect.Array && t != reflect.Slice {
		// do nothing
		return
	}

	mdsArr := reflect.Indirect(reflect.ValueOf(mds))
	for i := 0; i < mdsArr.Len(); i++ {
		d.setDefaultValue(ctx, mdsArr.Index(i).Interface())
	}
	logs.Warn("%v", mdsArr.Index(0).Interface())
}

func (d *DefaultValueFilterChainBuilder) setDefaultValue(ctx context.Context, ins interface{}) {
	err := d.factory.AutoWire(ctx, nil, ins)
	if err != nil {
		logs.Error("try to wire the bean for orm.Insert failed. "+
			"the default value is not set: %v, ", err)
	}
}
