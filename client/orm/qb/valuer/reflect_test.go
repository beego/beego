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

package valuer

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/beego/beego/v2/client/orm/internal/models"
	"github.com/beego/beego/v2/client/orm/qb/errs"
	"github.com/beego/beego/v2/client/orm/qb/test"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestReflectValue_Field(t *testing.T) {
	testValueField(t, NewReflectValue)
	invalidCases := []valueFieldTestCase{
		{
			// 不存在的字段
			name:      "invalid field",
			field:     "UpdateTime",
			wantError: errs.NewErrUnknownField("UpdateTime"),
		},
	}
	t.Run("invalid cases", func(t *testing.T) {
		meta, err := models.DefaultModelCache.GetOrRegisterByMd(&test.SimpleStruct{})
		if err != nil {
			t.Fatal(err)
		}
		val := NewReflectValue(&test.SimpleStruct{}, meta)
		for _, tc := range invalidCases {
			t.Run(tc.name, func(t *testing.T) {
				v, err := val.Field(tc.field)
				assert.Equal(t, tc.wantError, err)
				if err != nil {
					return
				}
				assert.Equal(t, tc.wantVal, v.Interface())
			})
		}
	})
}

func Test_reflectValue_SetColumn(t *testing.T) {
	testSetColumn(t, NewReflectValue)
}

func FuzzReflectValue_Field(f *testing.F) {
	f.Fuzz(fuzzValueField(NewReflectValue))
}

func fuzzValueField(factory Creator) any {
	meta, _ := models.DefaultModelCache.GetOrRegisterByMd(&test.SimpleStruct{})
	return func(t *testing.T, b bool,
		i int, i8 int8, i16 int16, i32 int32, i64 int64,
		u uint, u8 uint8, u16 uint16, u32 uint32, u64 uint64,
		f32 float32, f64 float64, bt byte, bs []byte, s string) {
		cb := b
		entity := &test.SimpleStruct{
			Bool: b, BoolPtr: &cb,
			Int: i, IntPtr: &i,
			Int8: i8, Int8Ptr: &i8,
			Int16: i16, Int16Ptr: &i16,
			Int32: i32, Int32Ptr: &i32,
			Int64: i64, Int64Ptr: &i64,
			Uint: u, UintPtr: &u,
			Uint8: u8, Uint8Ptr: &u8,
			Uint16: u16, Uint16Ptr: &u16,
			Uint32: u32, Uint32Ptr: &u32,
			Uint64: u64, Uint64Ptr: &u64,
			Float32: f32, Float32Ptr: &f32,
			Float64: f64, Float64Ptr: &f64,
			String: s,
			//NullStringPtr:  &sql.NullString{String: s, Valid: b},
			//NullInt16Ptr:   &sql.NullInt16{Int16: i16, Valid: b},
			//NullInt32Ptr:   &sql.NullInt32{Int32: i32, Valid: b},
			//NullInt64Ptr:   &sql.NullInt64{Int64: i64, Valid: b},
			//NullBoolPtr:    &sql.NullBool{Bool: b, Valid: b},
			//NullFloat64Ptr: &sql.NullFloat64{Float64: f64, Valid: b},

			NullString:  sql.NullString{String: s, Valid: b},
			NullInt64:   sql.NullInt64{Int64: i64, Valid: b},
			NullBool:    sql.NullBool{Bool: b, Valid: b},
			NullFloat64: sql.NullFloat64{Float64: f64, Valid: b},
		}
		val := factory(entity, meta)
		cases := newValueFieldTestCases(entity)
		for _, c := range cases {
			v, err := val.Field(c.field)
			assert.Nil(t, err)
			assert.Equal(t, c.wantVal, v.Interface())
		}
	}
}

func BenchmarkReflectValue_Field(b *testing.B) {
	meta, _ := models.DefaultModelCache.GetOrRegisterByMd(&test.SimpleStruct{})
	ins := NewReflectValue(&test.SimpleStruct{Int64: 13}, meta)
	for i := 0; i < b.N; i++ {
		val, err := ins.Field("Int64")
		assert.Nil(b, err)
		assert.Equal(b, int64(13), val.Interface())
	}
}

func BenchmarkReflectValue_fieldByIndexes_VS_FieldByName(b *testing.B) {
	meta, _ := models.DefaultModelCache.GetOrRegisterByMd(&test.SimpleStruct{})
	ins := NewReflectValue(&test.SimpleStruct{Int64: 13}, meta)
	in, ok := ins.(*reflectValue)
	assert.True(b, ok)
	fieldName, unknownFieldName := "Int64", "XXXX"
	unknownValue := int64(13)
	var fieldValue reflect.Value
	b.Run("fieldByIndex found", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			val, ok := in.fieldByIndex(fieldName)
			assert.True(b, ok)
			assert.Equal(b, fieldValue, val.Interface())
		}
	})
	b.Run("fieldByIndex not found", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			val, ok := in.fieldByIndex(unknownFieldName)
			assert.False(b, ok)
			assert.Equal(b, unknownValue, val)
		}
	})
	b.Run("FieldByName found", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			val := in.val.FieldByName(fieldName)
			assert.Equal(b, fieldValue, val.Interface())
		}
	})
	b.Run("fieldByIndex not found", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			val := in.val.FieldByName(unknownFieldName)
			assert.Equal(b, unknownValue, val)
		}
	})
}

func testValueField(t *testing.T, creator Creator) {
	var tm test.SimpleStruct
	meta, err := models.DefaultModelCache.GetOrRegisterByMd(&tm)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("zero value", func(t *testing.T) {
		entity := &test.SimpleStruct{}
		testCases := newValueFieldTestCases(entity)
		val := creator(entity, meta)
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				v, err := val.Field(tc.field)
				assert.Equal(t, tc.wantError, err)
				if err != nil {
					return
				}
				assert.Equal(t, tc.wantVal, v.Interface())
			})
		}
	})
}

func testSetColumn(t *testing.T, creator Creator) {
	r := models.DefaultModelCache
	t.Run("types", func(t *testing.T) {
		testCases := []struct {
			name    string
			cs      map[string][]byte
			val     *test.SimpleStruct
			wantVal *test.SimpleStruct
			wantErr error
		}{
			{
				name: "normal value",
				cs: map[string][]byte{
					"id":          []byte("1"),
					"bool":        []byte("true"),
					"bool_ptr":    []byte("false"),
					"int":         []byte("12"),
					"int_ptr":     []byte("13"),
					"int8":        []byte("8"),
					"int8_ptr":    []byte("-8"),
					"int16":       []byte("16"),
					"int16_ptr":   []byte("-16"),
					"int32":       []byte("32"),
					"int32_ptr":   []byte("-32"),
					"int64":       []byte("64"),
					"int64_ptr":   []byte("-64"),
					"uint":        []byte("14"),
					"uint_ptr":    []byte("15"),
					"uint8":       []byte("8"),
					"uint8_ptr":   []byte("18"),
					"uint16":      []byte("16"),
					"uint16_ptr":  []byte("116"),
					"uint32":      []byte("32"),
					"uint32_ptr":  []byte("132"),
					"uint64":      []byte("64"),
					"uint64_ptr":  []byte("164"),
					"float32":     []byte("3.2"),
					"float32_ptr": []byte("-3.2"),
					"float64":     []byte("6.4"),
					"float64_ptr": []byte("-6.4"),
					"string":      []byte("world"),
					//"byte_array":       []byte("hello"),
					//"null_string_ptr":  []byte("null string"),
					//"null_int16_ptr":   []byte("16"),
					//"null_int32_ptr":   []byte("32"),
					//"null_int64_ptr":   []byte("64"),
					//"null_bool_ptr":    []byte("true"),
					//"null_float64_ptr": []byte("6.4"),
					//"json_column":      []byte(`{"name": "Tom"}`),

					"null_string":  []byte("null string"),
					"null_int64":   []byte("64"),
					"null_bool":    []byte("true"),
					"null_float64": []byte("6.4"),
				},
				val:     &test.SimpleStruct{},
				wantVal: test.NewSimpleStruct(1),
			},
			{
				name: "invalid field",
				cs: map[string][]byte{
					"invalid_column": nil,
				},
				wantErr: errs.NewErrUnknownField("invalid_column"),
			},
		}

		meta, err := r.GetOrRegisterByMd(&test.SimpleStruct{})
		if err != nil {
			t.Fatal(err)
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				defer func() { _ = db.Close() }()
				val := creator(tc.val, meta)
				cols := make([]string, 0, len(tc.cs))
				colVals := make([]driver.Value, 0, len(tc.cs))
				for k, v := range tc.cs {
					cols = append(cols, k)
					colVals = append(colVals, v)
				}
				mock.ExpectQuery("SELECT *").
					WillReturnRows(sqlmock.NewRows(cols).
						AddRow(colVals...))
				rows, _ := db.Query("SELECT *")
				rows.Next()
				err = val.SetColumns(rows)
				if err != nil {
					assert.Equal(t, tc.wantErr, err)
					return
				}
				if tc.wantErr != nil {
					t.Fatalf("期望得到错误，但是并没有得到 %v", tc.wantErr)
				}
				assert.Equal(t, tc.wantVal, tc.val)
			})
		}
	})

	type User struct {
		Name string
	}

	t.Run("invalid rows", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = db.Close() }()

		u := &User{}
		meta, err := r.GetOrRegisterByMd(u)
		if err != nil {
			t.Fatal(err)
		}
		val := creator(u, meta)
		// 多了一个列
		mock.ExpectQuery("SELECT *").
			WillReturnRows(sqlmock.NewRows([]string{"ID", "Name"}).
				AddRow(123, "Tom"))
		rows, _ := db.Query("SELECT *")
		rows.Next()
		err = val.SetColumns(rows)
		assert.Equal(t, errs.ErrTooManyColumns, err)

		// 读取列错误
		mock.ExpectQuery("SELECT *").
			WillReturnRows(sqlmock.NewRows([]string{"Name"}))
		rows, _ = db.Query("SELECT *")
		rows.Next()
		err = val.SetColumns(rows)
		assert.Equal(t, errors.New("sql: Rows are closed"), err)
	})

	type BaseEntity struct {
		Id         int64 `eorm:"primary_key"`
		CreateTime uint64
	}

	type CombinedUser struct {
		BaseEntity
		FirstName string
	}

	// 测试使用组合的场景
	t.Run("combination", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = db.Close() }()

		u := &CombinedUser{}
		meta, err := r.GetOrRegisterByMd(u)
		if err != nil {
			t.Fatal(err)
		}
		val := creator(u, meta)
		// 多了一个列
		mock.ExpectQuery("SELECT *").
			WillReturnRows(sqlmock.NewRows([]string{"id", "create_time", "first_name"}).
				AddRow(123, 100000, "Tom"))
		rows, _ := db.Query("SELECT *")
		rows.Next()
		err = val.SetColumns(rows)
		assert.NoError(t, err)
		wantUser := &CombinedUser{
			BaseEntity: BaseEntity{
				Id:         123,
				CreateTime: 100000,
			},
			FirstName: "Tom",
		}
		assert.Equal(t, wantUser, u)
	})
}

func newValueFieldTestCases(entity *test.SimpleStruct) []valueFieldTestCase {
	return []valueFieldTestCase{
		{
			name:    "bool",
			field:   "Bool",
			wantVal: entity.Bool,
		},
		{
			// bool 指针类型
			name:    "bool pointer",
			field:   "BoolPtr",
			wantVal: entity.BoolPtr,
		},
		{
			name:    "int",
			field:   "Int",
			wantVal: entity.Int,
		},
		{
			// int 指针类型
			name:    "int pointer",
			field:   "IntPtr",
			wantVal: entity.IntPtr,
		},
		{
			name:    "int8",
			field:   "Int8",
			wantVal: entity.Int8,
		},
		{
			name:    "int8 pointer",
			field:   "Int8Ptr",
			wantVal: entity.Int8Ptr,
		},
		{
			name:    "int16",
			field:   "Int16",
			wantVal: entity.Int16,
		},
		{
			name:    "int16 pointer",
			field:   "Int16Ptr",
			wantVal: entity.Int16Ptr,
		},
		{
			name:    "int32",
			field:   "Int32",
			wantVal: entity.Int32,
		},
		{
			name:    "int32 pointer",
			field:   "Int32Ptr",
			wantVal: entity.Int32Ptr,
		},
		{
			name:    "int64",
			field:   "Int64",
			wantVal: entity.Int64,
		},
		{
			name:    "int64 pointer",
			field:   "Int64Ptr",
			wantVal: entity.Int64Ptr,
		},
		{
			name:    "uint",
			field:   "Uint",
			wantVal: entity.Uint,
		},
		{
			name:    "uint pointer",
			field:   "UintPtr",
			wantVal: entity.UintPtr,
		},
		{
			name:    "uint8",
			field:   "Uint8",
			wantVal: entity.Uint8,
		},
		{
			name:    "uint8 pointer",
			field:   "Uint8Ptr",
			wantVal: entity.Uint8Ptr,
		},
		{
			name:    "uint16",
			field:   "Uint16",
			wantVal: entity.Uint16,
		},
		{
			name:    "uint16 pointer",
			field:   "Uint16Ptr",
			wantVal: entity.Uint16Ptr,
		},
		{
			name:    "uint32",
			field:   "Uint32",
			wantVal: entity.Uint32,
		},
		{
			name:    "uint32 pointer",
			field:   "Uint32Ptr",
			wantVal: entity.Uint32Ptr,
		},
		{
			name:    "uint64",
			field:   "Uint64",
			wantVal: entity.Uint64,
		},
		{
			name:    "uint64 pointer",
			field:   "Uint64Ptr",
			wantVal: entity.Uint64Ptr,
		},
		{
			name:    "float32",
			field:   "Float32",
			wantVal: entity.Float32,
		},
		{
			name:    "float32 pointer",
			field:   "Float32Ptr",
			wantVal: entity.Float32Ptr,
		},
		{
			name:    "float64",
			field:   "Float64",
			wantVal: entity.Float64,
		},
		{
			name:    "float64 pointer",
			field:   "Float64Ptr",
			wantVal: entity.Float64Ptr,
		},
		//{
		//	name:    "byte array",
		//	field:   "ByteArray",
		//	wantVal: entity.ByteArray,
		//},
		{
			name:    "string",
			field:   "String",
			wantVal: entity.String,
		},
		//{
		//	name:    "NullStringPtr",
		//	field:   "NullStringPtr",
		//	wantVal: entity.NullStringPtr,
		//},
		//{
		//	name:    "NullInt16Ptr",
		//	field:   "NullInt16Ptr",
		//	wantVal: entity.NullInt16Ptr,
		//},
		//{
		//	name:    "NullInt32Ptr",
		//	field:   "NullInt32Ptr",
		//	wantVal: entity.NullInt32Ptr,
		//},
		//{
		//	name:    "NullInt64Ptr",
		//	field:   "NullInt64Ptr",
		//	wantVal: entity.NullInt64Ptr,
		//},
		//{
		//	name:    "NullBoolPtr",
		//	field:   "NullBoolPtr",
		//	wantVal: entity.NullBoolPtr,
		//},
		//{
		//	name:    "NullFloat64Ptr",
		//	field:   "NullFloat64Ptr",
		//	wantVal: entity.NullFloat64Ptr,
		//},

		{
			name:    "NullString",
			field:   "NullString",
			wantVal: entity.NullString,
		},
		{
			name:    "NullInt64",
			field:   "NullInt64",
			wantVal: entity.NullInt64,
		},
		{
			name:    "NullBool",
			field:   "NullBool",
			wantVal: entity.NullBool,
		},
		{
			name:    "NullFloat64",
			field:   "NullFloat64",
			wantVal: entity.NullFloat64,
		},
	}
}

type valueFieldTestCase struct {
	name      string
	field     string
	wantVal   interface{}
	wantError error
}
