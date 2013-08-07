package orm

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
)

const defaultMaxIdle = 30

type DriverType int

const (
	_ DriverType = iota
	DR_MySQL
	DR_Sqlite
	DR_Oracle
	DR_Postgres
)

type driver string

func (d driver) Type() DriverType {
	a, _ := dataBaseCache.get(string(d))
	return a.Driver
}

func (d driver) Name() string {
	return string(d)
}

var _ Driver = new(driver)

var (
	dataBaseCache = &_dbCache{cache: make(map[string]*alias)}
	drivers       = map[string]DriverType{
		"mysql":    DR_MySQL,
		"postgres": DR_Postgres,
		"sqlite3":  DR_Sqlite,
	}
	dbBasers = map[DriverType]dbBaser{
		DR_MySQL:    newdbBaseMysql(),
		DR_Sqlite:   newdbBaseSqlite(),
		DR_Oracle:   newdbBaseMysql(),
		DR_Postgres: newdbBasePostgres(),
	}
)

type _dbCache struct {
	mux   sync.RWMutex
	cache map[string]*alias
}

func (ac *_dbCache) add(name string, al *alias) (added bool) {
	ac.mux.Lock()
	defer ac.mux.Unlock()
	if _, ok := ac.cache[name]; ok == false {
		ac.cache[name] = al
		added = true
	}
	return
}

func (ac *_dbCache) get(name string) (al *alias, ok bool) {
	ac.mux.RLock()
	defer ac.mux.RUnlock()
	al, ok = ac.cache[name]
	return
}

func (ac *_dbCache) getDefault() (al *alias) {
	al, _ = ac.get("default")
	return
}

type alias struct {
	Name       string
	Driver     DriverType
	DriverName string
	DataSource string
	MaxIdle    int
	DB         *sql.DB
	DbBaser    dbBaser
}

func RegisterDataBase(name, driverName, dataSource string, maxIdle int) {
	if maxIdle <= 0 {
		maxIdle = defaultMaxIdle
	}

	al := new(alias)
	al.Name = name
	al.DriverName = driverName
	al.DataSource = dataSource
	al.MaxIdle = maxIdle

	var (
		err error
	)

	if dr, ok := drivers[driverName]; ok {
		al.DbBaser = dbBasers[dr]
		al.Driver = dr
	} else {
		err = fmt.Errorf("driver name `%s` have not registered", driverName)
		goto end
	}

	if dataBaseCache.add(name, al) == false {
		err = fmt.Errorf("db name `%s` already registered, cannot reuse", name)
		goto end
	}

	al.DB, err = sql.Open(driverName, dataSource)
	if err != nil {
		err = fmt.Errorf("register db `%s`, %s", name, err.Error())
		goto end
	}

	err = al.DB.Ping()
	if err != nil {
		err = fmt.Errorf("register db `%s`, %s", name, err.Error())
		goto end
	}

end:
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

func RegisterDriver(name string, typ DriverType) {
	if t, ok := drivers[name]; ok == false {
		drivers[name] = typ
	} else {
		if t != typ {
			fmt.Println("name `%s` db driver already registered and is other type")
			os.Exit(2)
		}
	}
}
