package session

//CREATE TABLE `session` (
//  `session_key` char(64) NOT NULL,
//  `session_data` blob,
//  `session_expiry` int(11) unsigned NOT NULL,
//  PRIMARY KEY (`session_key`)
//) ENGINE=MyISAM DEFAULT CHARSET=utf8;

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

var mysqlpder = &MysqlProvider{}

type MysqlSessionStore struct {
	c      *sql.DB
	sid    string
	lock   sync.RWMutex
	values map[interface{}]interface{}
}

func (st *MysqlSessionStore) Set(key, value interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.values[key] = value
	st.updatemysql()
	return nil
}

func (st *MysqlSessionStore) Get(key interface{}) interface{} {
	st.lock.RLock()
	defer st.lock.RUnlock()
	if v, ok := st.values[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (st *MysqlSessionStore) Delete(key interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	delete(st.values, key)
	st.updatemysql()
	return nil
}

func (st *MysqlSessionStore) SessionID() string {
	return st.sid
}

func (st *MysqlSessionStore) updatemysql() {
	b, err := encodeGob(st.values)
	if err != nil {
		return
	}
	st.c.Exec("UPDATE session set `session_data`= ? where session_key=?", b, st.sid)
}

func (st *MysqlSessionStore) SessionRelease() {
	st.c.Close()
}

type MysqlProvider struct {
	maxlifetime int64
	savePath    string
}

func (mp *MysqlProvider) connectInit() *sql.DB {
	db, e := sql.Open("mysql", mp.savePath)
	if e != nil {
		return nil
	}
	return db
}

func (mp *MysqlProvider) SessionInit(maxlifetime int64, savePath string) error {
	mp.maxlifetime = maxlifetime
	mp.savePath = savePath
	return nil
}

func (mp *MysqlProvider) SessionRead(sid string) (SessionStore, error) {
	c := mp.connectInit()
	row := c.QueryRow("select session_data from session where session_key=?", sid)
	var sessiondata []byte
	err := row.Scan(&sessiondata)
	if err == sql.ErrNoRows {
		c.Exec("insert into session(`session_key`,`session_data`,`session_expiry`) values(?,?,?)", sid, "", time.Now().Unix())
	}
	var kv map[interface{}]interface{}
	if len(sessiondata) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = decodeGob(sessiondata)
		if err != nil {
			return nil, err
		}
	}
	rs := &MysqlSessionStore{c: c, sid: sid, values: kv}
	return rs, nil
}

func (mp *MysqlProvider) SessionDestroy(sid string) error {
	c := mp.connectInit()
	c.Exec("DELETE FROM session where session_key=?", sid)
	c.Close()
	return nil
}

func (mp *MysqlProvider) SessionGC() {
	c := mp.connectInit()
	c.Exec("DELETE from session where session_expiry < ?", time.Now().Unix()-mp.maxlifetime)
	c.Close()
	return
}

func init() {
	Register("mysql", mysqlpder)
}
