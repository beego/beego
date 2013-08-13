package orm

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"
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
	TZ         *time.Location
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

	al.DB.SetMaxIdleConns(al.MaxIdle)

	// orm timezone system match database
	// default use Local
	al.TZ = time.Local

	switch al.Driver {
	case DR_MySQL:
		row := al.DB.QueryRow("SELECT @@session.time_zone")
		var tz string
		row.Scan(&tz)
		if tz != "SYSTEM" {
			t, err := time.Parse("-07:00", tz)
			if err == nil {
				al.TZ = t.Location()
			}
		}
	case DR_Sqlite:
		al.TZ = time.UTC
	case DR_Postgres:
		row := al.DB.QueryRow("SELECT current_setting('TIMEZONE')")
		var tz string
		row.Scan(&tz)
		loc, err := time.LoadLocation(tz)
		if err == nil {
			al.TZ = loc
		}
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

func RegisterDriver(driverName string, typ DriverType) {
	if t, ok := drivers[driverName]; ok == false {
		drivers[driverName] = typ
	} else {
		if t != typ {
			fmt.Println("driverName `%s` db driver already registered and is other type")
			os.Exit(2)
		}
	}
}

func SetDataBaseTZ(name string, tz *time.Location) {
	if al, ok := dataBaseCache.get(name); ok {
		al.TZ = tz
	} else {
		err := fmt.Errorf("DataBase name `%s` not registered", name)
		fmt.Println(err)
	}
}
