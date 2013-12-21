package orm

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"
)

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
	Name         string
	Driver       DriverType
	DriverName   string
	DataSource   string
	MaxIdleConns int
	MaxOpenConns int
	DB           *sql.DB
	DbBaser      dbBaser
	TZ           *time.Location
	Engine       string
}

// Setting the database connect params. Use the database driver self dataSource args.
func RegisterDataBase(aliasName, driverName, dataSource string, params ...int) {
	al := new(alias)
	al.Name = aliasName
	al.DriverName = driverName
	al.DataSource = dataSource

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

	if dataBaseCache.add(aliasName, al) == false {
		err = fmt.Errorf("db name `%s` already registered, cannot reuse", aliasName)
		goto end
	}

	al.DB, err = sql.Open(driverName, dataSource)
	if err != nil {
		err = fmt.Errorf("register db `%s`, %s", aliasName, err.Error())
		goto end
	}

	// orm timezone system match database
	// default use Local
	al.TZ = time.Local

	switch al.Driver {
	case DR_MySQL:
		row := al.DB.QueryRow("SELECT @@session.time_zone")
		var tz string
		row.Scan(&tz)
		if tz == "SYSTEM" {
			tz = ""
			row = al.DB.QueryRow("SELECT @@system_time_zone")
			row.Scan(&tz)
			t, err := time.Parse("MST", tz)
			if err == nil {
				al.TZ = t.Location()
			}
		} else {
			t, err := time.Parse("-07:00", tz)
			if err == nil {
				al.TZ = t.Location()
			}
		}

		// get default engine from current database
		row = al.DB.QueryRow("SELECT ENGINE, TRANSACTIONS FROM information_schema.engines WHERE SUPPORT = 'DEFAULT'")
		var engine string
		var tx bool
		row.Scan(&engine, &tx)

		if engine != "" {
			al.Engine = engine
		} else {
			engine = "INNODB"
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

	for i, v := range params {
		switch i {
		case 0:
			SetMaxIdleConns(al.Name, v)
		case 1:
			SetMaxOpenConns(al.Name, v)
		}
	}

	err = al.DB.Ping()
	if err != nil {
		err = fmt.Errorf("register db `%s`, %s", aliasName, err.Error())
		goto end
	}

end:
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

// Register a database driver use specify driver name, this can be definition the driver is which database type.
func RegisterDriver(driverName string, typ DriverType) {
	if t, ok := drivers[driverName]; ok == false {
		drivers[driverName] = typ
	} else {
		if t != typ {
			fmt.Sprintf("driverName `%s` db driver already registered and is other type\n", driverName)
			os.Exit(2)
		}
	}
}

// Change the database default used timezone
func SetDataBaseTZ(aliasName string, tz *time.Location) {
	if al, ok := dataBaseCache.get(aliasName); ok {
		al.TZ = tz
	} else {
		fmt.Sprintf("DataBase name `%s` not registered\n", aliasName)
		os.Exit(2)
	}
}

// Change the max idle conns for *sql.DB, use specify database alias name
func SetMaxIdleConns(aliasName string, maxIdleConns int) {
	al := getDbAlias(aliasName)
	al.MaxIdleConns = maxIdleConns
	al.DB.SetMaxIdleConns(maxIdleConns)
}

// Change the max open conns for *sql.DB, use specify database alias name
func SetMaxOpenConns(aliasName string, maxOpenConns int) {
	al := getDbAlias(aliasName)
	al.MaxOpenConns = maxOpenConns
	// for tip go 1.2
	if fun := reflect.ValueOf(al.DB).MethodByName("SetMaxOpenConns"); fun.IsValid() {
		fun.Call([]reflect.Value{reflect.ValueOf(maxOpenConns)})
	}
}
