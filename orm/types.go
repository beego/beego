package orm

import (
	"database/sql"
	"reflect"
)

type Driver interface {
	Name() string
	Type() DriverType
}

type Fielder interface {
	String() string
	FieldType() int
	SetRaw(interface{}) error
	RawValue() interface{}
	Clean() error
}

type Modeler interface {
	Init(Modeler, ...interface{}) Modeler
	IsInited() bool
	Clean() IFieldErrors
	CleanFields(string) IFieldErrors
	GetTableName() string
}

type Ormer interface {
	Read(Modeler) error
	Insert(Modeler) (int64, error)
	Update(Modeler) (int64, error)
	Delete(Modeler) (int64, error)
	M2mAdd(Modeler, string, ...interface{}) (int64, error)
	M2mDel(Modeler, string, ...interface{}) (int64, error)
	LoadRel(Modeler, string) (int64, error)
	QueryTable(interface{}) QuerySeter
	Using(string) error
	Begin() error
	Commit() error
	Rollback() error
	Raw(string, ...interface{}) RawSeter
	Driver() Driver
}

type Inserter interface {
	Insert(Modeler) (int64, error)
	Close() error
}

type QuerySeter interface {
	Filter(string, ...interface{}) QuerySeter
	Exclude(string, ...interface{}) QuerySeter
	SetCond(*Condition) QuerySeter
	Limit(int, ...int64) QuerySeter
	Offset(int64) QuerySeter
	OrderBy(...string) QuerySeter
	RelatedSel(...interface{}) QuerySeter
	Count() (int64, error)
	Update(Params) (int64, error)
	Delete() (int64, error)
	PrepareInsert() (Inserter, error)
	All(interface{}) (int64, error)
	One(Modeler) error
	Values(*[]Params, ...string) (int64, error)
	ValuesList(*[]ParamsList, ...string) (int64, error)
	ValuesFlat(*ParamsList, string) (int64, error)
}

type RawPreparer interface {
	Exec(...interface{}) (int64, error)
	Close() error
}

type RawSeter interface {
	Exec() (int64, error)
	QueryRow(...interface{}) error
	QueryRows(...interface{}) (int64, error)
	SetArgs(...interface{}) RawSeter
	Values(*[]Params) (int64, error)
	ValuesList(*[]ParamsList) (int64, error)
	ValuesFlat(*ParamsList) (int64, error)
	Prepare() (RawPreparer, error)
}

type IFieldError interface {
	Name() string
	Error() error
}

type IFieldErrors interface {
	Get(string) IFieldError
	Set(string, IFieldError)
	List() []IFieldError
}

type dbQuerier interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type dbBaser interface {
	Read(dbQuerier, *modelInfo, reflect.Value) error
	Insert(dbQuerier, *modelInfo, reflect.Value) (int64, error)
	InsertStmt(*sql.Stmt, *modelInfo, reflect.Value) (int64, error)
	Update(dbQuerier, *modelInfo, reflect.Value) (int64, error)
	Delete(dbQuerier, *modelInfo, reflect.Value) (int64, error)
	ReadBatch(dbQuerier, *querySet, *modelInfo, *Condition, interface{}) (int64, error)
	UpdateBatch(dbQuerier, *querySet, *modelInfo, *Condition, Params) (int64, error)
	DeleteBatch(dbQuerier, *querySet, *modelInfo, *Condition) (int64, error)
	Count(dbQuerier, *querySet, *modelInfo, *Condition) (int64, error)
	GetOperatorSql(*modelInfo, string, []interface{}) (string, []interface{})
	PrepareInsert(dbQuerier, *modelInfo) (*sql.Stmt, error)
	ReadValues(dbQuerier, *querySet, *modelInfo, *Condition, []string, interface{}) (int64, error)
}
