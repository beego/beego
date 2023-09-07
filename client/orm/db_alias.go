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

	"github.com/beego/beego/v2/client/orm/internal/logs"

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
	dataBaseCache = &_dbCache{cache: make(map[string]*alias)}
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
	cache map[string]*alias
}

// add database alias with original name.
func (ac *_dbCache) add(name string, al *alias) (added bool) {
	ac.mux.Lock()
	defer ac.mux.Unlock()
	if _, ok := ac.cache[name]; !ok {
		ac.cache[name] = al
		added = true
	}
	return
}

// get database alias if cached.
func (ac *_dbCache) get(name string) (al *alias, ok bool) {
	ac.mux.RLock()
	defer ac.mux.RUnlock()
	al, ok = ac.cache[name]
	return
}

// get default alias.
func (ac *_dbCache) getDefault() (al *alias) {
	al, _ = ac.get("default")
	return
}

type DB struct {
	*sync.RWMutex
	DB                  *sql.DB
	stmtDecorators      *lru.Cache
	stmtDecoratorsLimit int
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

type alias struct {
	Name            string
	Driver          DriverType
	DriverName      string
	DataSource      string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	StmtCacheSize   int
	DB              *DB
	DbBaser         dbBaser
	TZ              *time.Location
	Engine          string
}

func detectTZ(al *alias) {
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
				logs.DebugLog.Printf("Detect DB timezone: %s %s\n", tz, err.Error())
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
			logs.DebugLog.Printf("Detect DB timezone: %s %s\n", tz, err.Error())
		}
	}
}

func addAliasWthDB(aliasName, driverName string, db *sql.DB, params ...DBOption) (*alias, error) {
	existErr := fmt.Errorf("DataBase alias name `%s` already registered, cannot reuse", aliasName)
	if _, ok := dataBaseCache.get(aliasName); ok {
		return nil, existErr
	}

	al, err := newAliasWithDb(aliasName, driverName, db, params...)
	if err != nil {
		return nil, err
	}

	if !dataBaseCache.add(aliasName, al) {
		return nil, existErr
	}

	return al, nil
}

func newAliasWithDb(aliasName, driverName string, db *sql.DB, params ...DBOption) (*alias, error) {
	al := &alias{}
	al.DB = &DB{
		RWMutex: new(sync.RWMutex),
		DB:      db,
	}

	for _, p := range params {
		p(al)
	}

	var stmtCache *lru.Cache
	var stmtCacheSize int

	if al.StmtCacheSize > 0 {
		_stmtCache, errC := newStmtDecoratorLruWithEvict(al.StmtCacheSize)
		if errC != nil {
			return nil, errC
		} else {
			stmtCache = _stmtCache
			stmtCacheSize = al.StmtCacheSize
		}
	}

	al.Name = aliasName
	al.DriverName = driverName
	al.DB.stmtDecorators = stmtCache
	al.DB.stmtDecoratorsLimit = stmtCacheSize

	if dr, ok := drivers[driverName]; ok {
		al.DbBaser = dbBasers[dr]
		al.Driver = dr
	} else {
		return nil, fmt.Errorf("driver name `%s` have not registered", driverName)
	}

	err := db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Register db Ping `%s`, %s", aliasName, err.Error())
	}

	detectTZ(al)

	return al, nil
}

// SetMaxIdleConns Change the max idle conns for *sql.DB, use specify database alias name
// Deprecated you should not use this, we will remove it in the future
func SetMaxIdleConns(aliasName string, maxIdleConns int) {
	al := getDbAlias(aliasName)
	al.SetMaxIdleConns(maxIdleConns)
}

// SetMaxOpenConns Change the max open conns for *sql.DB, use specify database alias name
// Deprecated you should not use this, we will remove it in the future
func SetMaxOpenConns(aliasName string, maxOpenConns int) {
	al := getDbAlias(aliasName)
	al.SetMaxOpenConns(maxOpenConns)
}

// SetMaxIdleConns Change the max idle conns for *sql.DB, use specify database alias name
func (al *alias) SetMaxIdleConns(maxIdleConns int) {
	al.MaxIdleConns = maxIdleConns
	al.DB.DB.SetMaxIdleConns(maxIdleConns)
}

// SetMaxOpenConns Change the max open conns for *sql.DB, use specify database alias name
func (al *alias) SetMaxOpenConns(maxOpenConns int) {
	al.MaxOpenConns = maxOpenConns
	al.DB.DB.SetMaxOpenConns(maxOpenConns)
}

func (al *alias) SetConnMaxLifetime(lifeTime time.Duration) {
	al.ConnMaxLifetime = lifeTime
	al.DB.DB.SetConnMaxLifetime(lifeTime)
}

// AddAliasWthDB add a aliasName for the drivename
func AddAliasWthDB(aliasName, driverName string, db *sql.DB, params ...DBOption) error {
	_, err := addAliasWthDB(aliasName, driverName, db, params...)
	return err
}

// RegisterDataBase Setting the database connect params. Use the database driver self dataSource args.
func RegisterDataBase(aliasName, driverName, dataSource string, params ...DBOption) error {
	var (
		err error
		db  *sql.DB
		al  *alias
	)

	db, err = sql.Open(driverName, dataSource)
	if err != nil {
		err = fmt.Errorf("Register db `%s`, %s", aliasName, err.Error())
		goto end
	}

	al, err = addAliasWthDB(aliasName, driverName, db, params...)
	if err != nil {
		goto end
	}

	al.DataSource = dataSource

end:
	if err != nil {
		if db != nil {
			db.Close()
		}
		logs.DebugLog.Println(err.Error())
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
func SetDataBaseTZ(aliasName string, tz *time.Location) error {
	if al, ok := dataBaseCache.get(aliasName); ok {
		al.TZ = tz
	} else {
		return fmt.Errorf("DataBase alias name `%s` not registered", aliasName)
	}
	return nil
}

// GetDB Get *sql.DB from registered database by db alias name.
// Use "default" as alias name if you not Set.
func GetDB(aliasNames ...string) (*sql.DB, error) {
	var name string
	if len(aliasNames) > 0 {
		name = aliasNames[0]
	} else {
		name = "default"
	}
	al, ok := dataBaseCache.get(name)
	if ok {
		return al.DB.DB, nil
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

type DBOption func(al *alias)

// MaxIdleConnections return a hint about MaxIdleConnections
func MaxIdleConnections(maxIdleConn int) DBOption {
	return func(al *alias) {
		al.SetMaxIdleConns(maxIdleConn)
	}
}

// MaxOpenConnections return a hint about MaxOpenConnections
func MaxOpenConnections(maxOpenConn int) DBOption {
	return func(al *alias) {
		al.SetMaxOpenConns(maxOpenConn)
	}
}

// ConnMaxLifetime return a hint about ConnMaxLifetime
func ConnMaxLifetime(v time.Duration) DBOption {
	return func(al *alias) {
		al.SetConnMaxLifetime(v)
	}
}

// MaxStmtCacheSize return a hint about MaxStmtCacheSize
func MaxStmtCacheSize(v int) DBOption {
	return func(al *alias) {
		al.StmtCacheSize = v
	}
}
