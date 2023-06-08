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

package orm

import (
	"context"
	"reflect"

	"github.com/beego/beego/v2/client/orm/internal/models"
)

// model to model struct
type queryM2M struct {
	md  interface{}
	mi  *models.ModelInfo
	fi  *models.FieldInfo
	qs  *querySet
	ind reflect.Value
}

// add models to origin models when creating queryM2M.
// example:
//
//		m2m := orm.QueryM2M(post,"Tag")
//		m2m.Add(&Tag1{},&Tag2{})
//	 for _,tag := range post.Tags{}
//
// make sure the relation is defined in post model struct tag.
func (o *queryM2M) Add(mds ...interface{}) (int64, error) {
	return o.AddWithCtx(context.Background(), mds...)
}

func (o *queryM2M) AddWithCtx(ctx context.Context, mds ...interface{}) (int64, error) {
	fi := o.fi
	mi := fi.RelThroughModelInfo
	mfi := fi.ReverseFieldInfo
	rfi := fi.ReverseFieldInfoTwo

	orm := o.qs.orm
	dbase := orm.alias.DbBaser

	var models []interface{}
	var otherValues []interface{}
	var otherNames []string

	for _, colname := range mi.Fields.DBcols {
		if colname != mfi.Column && colname != rfi.Column && colname != fi.Mi.Fields.Pk.Column &&
			mi.Fields.Columns[colname] != mi.Fields.Pk {
			otherNames = append(otherNames, colname)
		}
	}
	for i, md := range mds {
		if reflect.Indirect(reflect.ValueOf(md)).Kind() != reflect.Struct && i > 0 {
			otherValues = append(otherValues, md)
			mds = append(mds[:i], mds[i+1:]...)
		}
	}
	for _, md := range mds {
		val := reflect.ValueOf(md)
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			for i := 0; i < val.Len(); i++ {
				v := val.Index(i)
				if v.CanInterface() {
					models = append(models, v.Interface())
				}
			}
		} else {
			models = append(models, md)
		}
	}

	_, v1, exist := getExistPk(o.mi, o.ind)
	if !exist {
		panic(ErrMissPK)
	}

	names := []string{mfi.Column, rfi.Column}

	values := make([]interface{}, 0, len(models)*2)
	for _, md := range models {

		ind := reflect.Indirect(reflect.ValueOf(md))
		var v2 interface{}
		if ind.Kind() != reflect.Struct {
			v2 = ind.Interface()
		} else {
			_, v2, exist = getExistPk(fi.RelModelInfo, ind)
			if !exist {
				panic(ErrMissPK)
			}
		}
		values = append(values, v1, v2)

	}
	names = append(names, otherNames...)
	values = append(values, otherValues...)
	return dbase.InsertValue(ctx, orm.db, mi, true, names, values)
}

// remove models following the origin model relationship
func (o *queryM2M) Remove(mds ...interface{}) (int64, error) {
	return o.RemoveWithCtx(context.Background(), mds...)
}

func (o *queryM2M) RemoveWithCtx(ctx context.Context, mds ...interface{}) (int64, error) {
	fi := o.fi
	qs := o.qs.Filter(fi.ReverseFieldInfo.Name, o.md)

	return qs.Filter(fi.ReverseFieldInfoTwo.Name+ExprSep+"in", mds).Delete()
}

// check model is existed in relationship of origin model
func (o *queryM2M) Exist(md interface{}) bool {
	return o.ExistWithCtx(context.Background(), md)
}

func (o *queryM2M) ExistWithCtx(ctx context.Context, md interface{}) bool {
	fi := o.fi
	return o.qs.Filter(fi.ReverseFieldInfo.Name, o.md).
		Filter(fi.ReverseFieldInfoTwo.Name, md).ExistWithCtx(ctx)
}

// clean all models in related of origin model
func (o *queryM2M) Clear() (int64, error) {
	return o.ClearWithCtx(context.Background())
}

func (o *queryM2M) ClearWithCtx(ctx context.Context) (int64, error) {
	fi := o.fi
	return o.qs.Filter(fi.ReverseFieldInfo.Name, o.md).DeleteWithCtx(ctx)
}

// count all related models of origin model
func (o *queryM2M) Count() (int64, error) {
	return o.CountWithCtx(context.Background())
}

func (o *queryM2M) CountWithCtx(ctx context.Context) (int64, error) {
	fi := o.fi
	return o.qs.Filter(fi.ReverseFieldInfo.Name, o.md).CountWithCtx(ctx)
}

var _ QueryM2Mer = new(queryM2M)

// create new M2M queryer.
func newQueryM2M(md interface{}, o *ormBase, mi *models.ModelInfo, fi *models.FieldInfo, ind reflect.Value) QueryM2Mer {
	qm2m := new(queryM2M)
	qm2m.md = md
	qm2m.mi = mi
	qm2m.fi = fi
	qm2m.ind = ind
	qm2m.qs = newQuerySet(o, fi.RelThroughModelInfo).(*querySet)
	return qm2m
}
