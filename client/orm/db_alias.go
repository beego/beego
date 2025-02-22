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
	"database/sql"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

// DriverType database driver constant int.
type DriverType int

// Enum the Database driver
const (
	_          DriverType = iota // int enum type
	DRMySQL                      // mysql
	DRSqlite                     // sqlite
	DROracle                     // oracle
	DRPostgres                   // pgsql
	DRTiDB                       // TiDB
)

// database driver string.
type driver string

// Get type constant int of current driver..
func (d driver) Type() DriverType {
	a, _ := dataBaseCache.get(string(d))
	return a.Driver
}

// Get name of current driver
func (d driver) Name() string {
	return string(d)
}

// check driver iis implemented Driver interface or not.
var _ Driver = new(driver)

var (
	dataBaseCache = &_dbCache{cache: make(map[string]*DB)}
	drivers       = map[string]DriverType{
		"mysql":    DRMySQL,
		"postgres": DRPostgres,
		"sqlite3":  DRSqlite,
		"tidb":     DRTiDB,
		"oracle":   DROracle,
		"oci8":     DROracle, // github.com/mattn/go-oci8
		"ora":      DROracle, // https://github.com/rana/ora
	}
	dbBasers = map[DriverType]dbBaser{
		DRMySQL:    newdbBaseMysql(),
		DRSqlite:   newdbBaseSqlite(),
		DROracle:   newdbBaseOracle(),
		DRPostgres: newdbBasePostgres(),
		DRTiDB:     newdbBaseTidb(),
	}
)

// database alias cacher.
type _dbCache struct {
	mux   sync.RWMutex
	cache map[string]*DB
}

// add database db with original name.
func (ac *_dbCache) add(name string, db *DB) (added bool) {
	ac.mux.Lock()
	defer ac.mux.Unlock()
	if _, ok := ac.cache[name]; !ok {
		ac.cache[name] = db
		added = true
	}
	return
}

// get database alias if cached.
func (ac *_dbCache) get(name string) (db *DB, ok bool) {
	ac.mux.RLock()
	defer ac.mux.RUnlock()
	db, ok = ac.cache[name]
	return
}

func (ac *_dbCache) getORSet(aliasName, driverName string, db *sql.DB, params ...DBOption) (al *DB, err error) {
	ac.mux.RLock()
	d, ok := ac.cache[aliasName]
	ac.mux.RUnlock()
	if !ok {
		ac.mux.Lock()
		defer ac.mux.Unlock()
		al, err = newDB(aliasName, driverName, db, params...)
		if err != nil {
			return
		}
		ac.cache[aliasName] = d
	}
	return
}

// get default alias.
func (ac *_dbCache) getDefault() (db *DB) {
	db, _ = ac.get("default")
	return
}

type DB struct {
	*sync.RWMutex
	DB                  *sql.DB
	stmtDecorators      *lru.Cache
	stmtDecoratorsLimit int

	Name            string
	Driver          DriverType
	DriverName      string
	DataSource      string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdletime time.Duration
	StmtCacheSize   int
	//DB              *DB
	DbBaser dbBaser
	TZ      *time.Location
	Engine  string
}

var (
	_ dbQuerier = new(DB)
	_ txer      = new(DB)
)

func (d *DB) Begin() (*sql.Tx, error) {
	return d.DB.Begin()
}

func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.DB.BeginTx(ctx, opts)
}

// su must call release to release *sql.Stmt after using
func (d *DB) getStmtDecorator(query string) (*stmtDecorator, error) {
	d.RLock()
	c, ok := d.stmtDecorators.Get(query)
	if ok {
		c.(*stmtDecorator).acquire()
		d.RUnlock()
		return c.(*stmtDecorator), nil
	}
	d.RUnlock()

	d.Lock()
	c, ok = d.stmtDecorators.Get(query)
	if ok {
		c.(*stmtDecorator).acquire()
		d.Unlock()
		return c.(*stmtDecorator), nil
	}

	stmt, err := d.Prepare(query)
	if err != nil {
		d.Unlock()
		return nil, err
	}
	sd := newStmtDecorator(stmt)
	sd.acquire()
	d.stmtDecorators.Add(query, sd)
	d.Unlock()

	return sd, nil
}

func (d *DB) Prepare(query string) (*sql.Stmt, error) {
	return d.DB.Prepare(query)
}

func (d *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return d.DB.PrepareContext(ctx, query)
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.ExecContext(context.Background(), query, args...)
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if d.stmtDecorators == nil {
		return d.DB.ExecContext(ctx, query, args...)
	}

	sd, err := d.getStmtDecorator(query)
	if err != nil {
		return nil, err
	}
	stmt := sd.getStmt()
	defer sd.release()
	return stmt.ExecContext(ctx, args...)
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.QueryContext(context.Background(), query, args...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if d.stmtDecorators == nil {
		return d.DB.QueryContext(ctx, query, args...)
	}

	sd, err := d.getStmtDecorator(query)
	if err != nil {
		return nil, err
	}
	stmt := sd.getStmt()
	defer sd.release()
	return stmt.QueryContext(ctx, args...)
}

func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.QueryRowContext(context.Background(), query, args...)
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if d.stmtDecorators == nil {
		return d.DB.QueryRowContext(ctx, query, args...)
	}

	sd, err := d.getStmtDecorator(query)
	if err != nil {
		panic(err)
	}
	stmt := sd.getStmt()
	defer sd.release()
	return stmt.QueryRowContext(ctx, args...)
}

type TxDB struct {
	tx *sql.Tx
}

var (
	_ dbQuerier = new(TxDB)
	_ txEnder   = new(TxDB)
)

func (t *TxDB) Commit() error {
	return t.tx.Commit()
}

func (t *TxDB) Rollback() error {
	return t.tx.Rollback()
}

func (t *TxDB) RollbackUnlessCommit() error {
	err := t.tx.Rollback()
	if err != sql.ErrTxDone {
		return err
	}
	return nil
}

var (
	_ dbQuerier = new(TxDB)
	_ txEnder   = new(TxDB)
)

func (t *TxDB) Prepare(query string) (*sql.Stmt, error) {
	return t.PrepareContext(context.Background(), query)
}

func (t *TxDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return t.tx.PrepareContext(ctx, query)
}

func (t *TxDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.ExecContext(context.Background(), query, args...)
}

func (t *TxDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *TxDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.QueryContext(context.Background(), query, args...)
}

func (t *TxDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *TxDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.QueryRowContext(context.Background(), query, args...)
}

func (t *TxDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

//type alias struct {
//	Name            string
//	Driver          DriverType
//	DriverName      string
//	DataSource      string
//	MaxIdleConns    int
//	MaxOpenConns    int
//	ConnMaxLifetime time.Duration
//	ConnMaxIdletime time.Duration
//	StmtCacheSize   int
//	DB              *DB
//	DbBaser         dbBaser
//	TZ              *time.Location
//	Engine          string
//}

func detectTZ(al *DB) {
	// orm timezone system match database
	// default use Local
	al.TZ = DefaultTimeLoc

	if al.DriverName == "sphinx" {
		return
	}

	switch al.Driver {
	case DRMySQL:
		row := al.DB.QueryRow("SELECT TIMEDIFF(NOW(), UTC_TIMESTAMP)")
		var tz string
		row.Scan(&tz)
		if len(tz) >= 8 {
			if tz[0] != '-' {
				tz = "+" + tz
			}
			t, err := time.Parse("-07:00:00", tz)
			if err == nil {
				if t.Location().String() != "" {
					al.TZ = t.Location()
				}
			} else {
				DebugLog.Printf("Detect DB timezone: %s %s\n", tz, err.Error())
			}
		}

		// Get default engine from current database
		row = al.DB.QueryRow("SELECT ENGINE, TRANSACTIONS FROM information_schema.engines WHERE SUPPORT = 'DEFAULT'")
		var engine string
		var tx bool
		row.Scan(&engine, &tx)

		if engine != "" {
			al.Engine = engine
		} else {
			al.Engine = "INNODB"
		}

	case DRSqlite, DROracle:
		al.TZ = time.UTC

	case DRPostgres:
		row := al.DB.QueryRow("SELECT current_setting('TIMEZONE')")
		var tz string
		row.Scan(&tz)
		loc, err := time.LoadLocation(tz)
		if err == nil {
			al.TZ = loc
		} else {
			DebugLog.Printf("Detect DB timezone: %s %s\n", tz, err.Error())
		}
	}
}

func addDB(aliasName, driverName string, db *sql.DB, params ...DBOption) (*DB, error) {
	existErr := fmt.Errorf("DataBase alias name `%s` already registered, cannot reuse", aliasName)
	if _, ok := dataBaseCache.get(aliasName); ok {
		return nil, existErr
	}

	al, err := newDB(aliasName, driverName, db, params...)
	if err != nil {
		return nil, err
	}

	if !dataBaseCache.add(aliasName, al) {
		return nil, existErr
	}

	return al, nil
}

func getORSetDB(aliasName, driverName string, db *sql.DB, params ...DBOption) (*DB, error) {
	al, err := dataBaseCache.getORSet(aliasName, driverName, db, params...)
	if err != nil {
		return nil, err
	}

	return al, nil
}

func newDB(aliasName, driverName string, db *sql.DB, params ...DBOption) (*DB, error) {
	//al := &alias{}
	//al.DB = &DB{
	//	RWMutex: new(sync.RWMutex),
	//	DB:      db,
	//}

	res := &DB{
		RWMutex: new(sync.RWMutex),
		DB:      db,
	}
	for _, p := range params {
		p(res)
	}

	var stmtCache *lru.Cache
	var stmtCacheSize int

	if res.StmtCacheSize > 0 {
		_stmtCache, errC := newStmtDecoratorLruWithEvict(res.StmtCacheSize)
		if errC != nil {
			return nil, errC
		} else {
			stmtCache = _stmtCache
			stmtCacheSize = res.StmtCacheSize
		}
	}

	res.Name = aliasName
	res.DriverName = driverName
	res.stmtDecorators = stmtCache
	res.stmtDecoratorsLimit = stmtCacheSize

	if dr, ok := drivers[driverName]; ok {
		res.DbBaser = dbBasers[dr]
		res.Driver = dr
	} else {
		return nil, fmt.Errorf("driver name `%s` have not registered", driverName)
	}

	err := db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Register db Ping `%s`, %s", aliasName, err.Error())
	}

	detectTZ(res)

	return res, nil
}

// SetMaxIdleConns Change the max idle conns for *sql.DB, use specify database alias name
// Deprecated you should not use this, we will remove it in the future
func SetMaxIdleConns(aliasName string, maxIdleConns int) {
	d := getDB(aliasName)
	d.SetMaxIdleConns(maxIdleConns)
}

// SetMaxOpenConns Change the max open conns for *sql.DB, use specify database alias name
// Deprecated you should not use this, we will remove it in the future
func SetMaxOpenConns(aliasName string, maxOpenConns int) {
	d := getDB(aliasName)
	d.SetMaxOpenConns(maxOpenConns)
}

// SetMaxIdleConns Change the max idle conns for *sql.DB, use specify database alias name
func (d *DB) SetMaxIdleConns(maxIdleConns int) {
	d.MaxIdleConns = maxIdleConns
	d.DB.SetMaxIdleConns(maxIdleConns)
}

// SetMaxOpenConns Change the max open conns for *sql.DB, use specify database alias name
func (d *DB) SetMaxOpenConns(maxOpenConns int) {
	d.MaxOpenConns = maxOpenConns
	d.DB.SetMaxOpenConns(maxOpenConns)
}

func (d *DB) SetConnMaxLifetime(lifeTime time.Duration) {
	d.ConnMaxLifetime = lifeTime
	d.DB.SetConnMaxLifetime(lifeTime)
}

func (d *DB) SetConnMaxIdleTime(idleTime time.Duration) {
	d.ConnMaxIdletime = idleTime
	d.DB.SetConnMaxIdleTime(idleTime)
}

// AddDB add a aliasName for the drivename
func AddDB(aliasName, driverName string, db *sql.DB, params ...DBOption) error {
	_, err := addDB(aliasName, driverName, db, params...)
	return err
}

// RegisterDataBase Setting the database connect params. Use the database driver self dataSource args.
func RegisterDataBase(aliasName, driverName, dataSource string, params ...DBOption) error {
	var (
		err error
		db  *sql.DB
		res *DB
	)
	db, err = sql.Open(driverName, dataSource)
	if err != nil {
		err = fmt.Errorf("register db `%s`, %s", aliasName, err.Error())
		goto end
	}

	res, err = addDB(aliasName, driverName, db, params...)
	if err != nil {
		goto end
	}

	res.DataSource = dataSource

end:
	if err != nil {
		if db != nil {
			db.Close()
		}
		DebugLog.Println(err.Error())
	}

	return err
}

func RegisterDB(aliasName, driverName, dataSource string, params ...DBOption) error {
	var (
		err error
		db  *sql.DB
		res *DB
	)
	db, err = sql.Open(driverName, dataSource)
	if err != nil {
		err = fmt.Errorf("Register db `%s`, %s", aliasName, err.Error())
		goto end
	}

	res, err = getORSetDB(aliasName, driverName, db, params...)
	if err != nil {
		goto end
	}

	if res.DataSource != dataSource {
		res.DataSource = dataSource
	}

end:
	if err != nil {
		if db != nil {
			_ = db.Close()
		}
		DebugLog.Println(err.Error())
	}

	return err
}

// RegisterDriver Register a database driver use specify driver name, this can be definition the driver is which database type.
func RegisterDriver(driverName string, typ DriverType) error {
	if t, ok := drivers[driverName]; !ok {
		drivers[driverName] = typ
	} else {
		if t != typ {
			return fmt.Errorf("driverName `%s` db driver already registered and is other type", driverName)
		}
	}
	return nil
}

// SetDataBaseTZ Change the database default used timezone
func SetDataBaseTZ(dbName string, tz *time.Location) error {
	if db, ok := dataBaseCache.get(dbName); ok {
		db.TZ = tz
	} else {
		return fmt.Errorf("DataBase alias name `%s` not registered", dbName)
	}
	return nil
}

// GetSqlDB Get *sql.DB from registered database by db alias name.
// Use "default" as alias name if you not Set.
func GetSqlDB(dbNames ...string) (*sql.DB, error) {
	var name string
	if len(dbNames) > 0 {
		name = dbNames[0]
	} else {
		name = "default"
	}
	al, ok := dataBaseCache.get(name)
	if ok {
		return al.DB, nil
	}
	return nil, fmt.Errorf("DataBase of alias name `%s` not found", name)
}

type stmtDecorator struct {
	wg   sync.WaitGroup
	stmt *sql.Stmt
}

func (s *stmtDecorator) getStmt() *sql.Stmt {
	return s.stmt
}

// acquire will add one
// since this method will be used inside read lock scope,
// so we can not do more things here
// we should think about refactor this
func (s *stmtDecorator) acquire() {
	s.wg.Add(1)
}

func (s *stmtDecorator) release() {
	s.wg.Done()
}

// garbage recycle for stmt
func (s *stmtDecorator) destroy() {
	go func() {
		s.wg.Wait()
		_ = s.stmt.Close()
	}()
}

func newStmtDecorator(sqlStmt *sql.Stmt) *stmtDecorator {
	return &stmtDecorator{
		stmt: sqlStmt,
	}
}

func newStmtDecoratorLruWithEvict(cacheSize int) (*lru.Cache, error) {
	cache, err := lru.NewWithEvict(cacheSize, func(key interface{}, value interface{}) {
		value.(*stmtDecorator).destroy()
	})
	if err != nil {
		return nil, err
	}
	return cache, nil
}

type DBOption func(d *DB)

// MaxIdleConnections return a hint about MaxIdleConnections
func MaxIdleConnections(maxIdleConn int) DBOption {
	return func(d *DB) {
		d.SetMaxIdleConns(maxIdleConn)
	}
}

// MaxOpenConnections return a hint about MaxOpenConnections
func MaxOpenConnections(maxOpenConn int) DBOption {
	return func(d *DB) {
		d.SetMaxOpenConns(maxOpenConn)
	}
}

// ConnMaxLifetime return a hint about ConnMaxLifetime
func ConnMaxLifetime(v time.Duration) DBOption {
	return func(d *DB) {
		d.SetConnMaxLifetime(v)
	}
}

// ConnMaxIdletime return a hint about ConnMaxIdletime
func ConnMaxIdletime(v time.Duration) DBOption {
	return func(d *DB) {
		d.SetConnMaxIdleTime(v)
	}
}

// MaxStmtCacheSize return a hint about MaxStmtCacheSize
func MaxStmtCacheSize(v int) DBOption {
	return func(d *DB) {
		d.StmtCacheSize = v
	}
}
