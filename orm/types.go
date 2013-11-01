package orm

import (
	"database/sql"
	"reflect"
	"time"
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
}

type Ormer interface {
	Read(interface{}, ...string) error
	Insert(interface{}) (int64, error)
	Update(interface{}, ...string) (int64, error)
	Delete(interface{}) (int64, error)
	LoadRelated(interface{}, string, ...interface{}) (int64, error)
	QueryM2M(interface{}, string) QueryM2Mer
	QueryTable(interface{}) QuerySeter
	Using(string) error
	Begin() error
	Commit() error
	Rollback() error
	Raw(string, ...interface{}) RawSeter
	Driver() Driver
}

type Inserter interface {
	Insert(interface{}) (int64, error)
	Close() error
}

type QuerySeter interface {
	Filter(string, ...interface{}) QuerySeter
	Exclude(string, ...interface{}) QuerySeter
	SetCond(*Condition) QuerySeter
	Limit(interface{}, ...interface{}) QuerySeter
	Offset(interface{}) QuerySeter
	OrderBy(...string) QuerySeter
	RelatedSel(...interface{}) QuerySeter
	Count() (int64, error)
	Exist() bool
	Update(Params) (int64, error)
	Delete() (int64, error)
	PrepareInsert() (Inserter, error)
	All(interface{}, ...string) (int64, error)
	One(interface{}, ...string) error
	Values(*[]Params, ...string) (int64, error)
	ValuesList(*[]ParamsList, ...string) (int64, error)
	ValuesFlat(*ParamsList, string) (int64, error)
}

type QueryM2Mer interface {
	Add(...interface{}) (int64, error)
	Remove(...interface{}) (int64, error)
	Exist(interface{}) bool
	Clear() (int64, error)
	Count() (int64, error)
}

type RawPreparer interface {
	Exec(...interface{}) (sql.Result, error)
	Close() error
}

type RawSeter interface {
	Exec() (sql.Result, error)
	QueryRow(...interface{}) error
	QueryRows(...interface{}) (int64, error)
	SetArgs(...interface{}) RawSeter
	Values(*[]Params) (int64, error)
	ValuesList(*[]ParamsList) (int64, error)
	ValuesFlat(*ParamsList) (int64, error)
	Prepare() (RawPreparer, error)
}

type stmtQuerier interface {
	Close() error
	Exec(args ...interface{}) (sql.Result, error)
	Query(args ...interface{}) (*sql.Rows, error)
	QueryRow(args ...interface{}) *sql.Row
}

type dbQuerier interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type txer interface {
	Begin() (*sql.Tx, error)
}

type txEnder interface {
	Commit() error
	Rollback() error
}

type dbBaser interface {
	Read(dbQuerier, *modelInfo, reflect.Value, *time.Location, []string) error
	Insert(dbQuerier, *modelInfo, reflect.Value, *time.Location) (int64, error)
	InsertValue(dbQuerier, *modelInfo, []string, []interface{}) (int64, error)
	InsertStmt(stmtQuerier, *modelInfo, reflect.Value, *time.Location) (int64, error)
	Update(dbQuerier, *modelInfo, reflect.Value, *time.Location, []string) (int64, error)
	Delete(dbQuerier, *modelInfo, reflect.Value, *time.Location) (int64, error)
	ReadBatch(dbQuerier, *querySet, *modelInfo, *Condition, interface{}, *time.Location, []string) (int64, error)
	SupportUpdateJoin() bool
	UpdateBatch(dbQuerier, *querySet, *modelInfo, *Condition, Params, *time.Location) (int64, error)
	DeleteBatch(dbQuerier, *querySet, *modelInfo, *Condition, *time.Location) (int64, error)
	Count(dbQuerier, *querySet, *modelInfo, *Condition, *time.Location) (int64, error)
	OperatorSql(string) string
	GenerateOperatorSql(*modelInfo, *fieldInfo, string, []interface{}, *time.Location) (string, []interface{})
	GenerateOperatorLeftCol(*fieldInfo, string, *string)
	PrepareInsert(dbQuerier, *modelInfo) (stmtQuerier, string, error)
	ReadValues(dbQuerier, *querySet, *modelInfo, *Condition, []string, interface{}, *time.Location) (int64, error)
	MaxLimit() uint64
	TableQuote() string
	ReplaceMarks(*string)
	HasReturningID(*modelInfo, *string) bool
	TimeFromDB(*time.Time, *time.Location)
	TimeToDB(*time.Time, *time.Location)
	DbTypes() map[string]string
	GetTables(dbQuerier) (map[string]bool, error)
	GetColumns(dbQuerier, string) (map[string][3]string, error)
	ShowTablesQuery() string
	ShowColumnsQuery(string) string
	IndexExists(dbQuerier, string, string) bool
	collectFieldValue(*modelInfo, *fieldInfo, reflect.Value, bool, *time.Location) (interface{}, error)
}
