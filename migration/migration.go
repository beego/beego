// Beego (http://beego.me/)
//
// @description beego is an open-source, high-performance web framework for the Go programming language.
//
// @link        http://github.com/astaxie/beego for the canonical source repository
//
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
//
// @authors     astaxie
//
//	CREATE TABLE `migrations` (
//		`id_migration` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'surrogate key',
//		`name` varchar(255) DEFAULT NULL COMMENT 'migration name, unique',
//		`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'date migrated or rolled back',
//		`statements` longtext COMMENT 'SQL statements for this migration',
//		`rollback_statements` longtext,
//		`status` enum('update','rollback') DEFAULT NULL COMMENT 'update indicates it is a normal migration while rollback means this migration is rolled back',
//		PRIMARY KEY (`id_migration`),
//		UNIQUE KEY `name` (`name`)
//	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
package migration

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// const the data format for the bee generate migration datatype
const M_DATE_FORMAT = "20060102_150405"

// Migrationer is an interface for all Migration struct
type Migrationer interface {
	Up()
	Down()
	Exec(name, status string) error
	GetCreated() int64
}

var migrationMap map[string]Migrationer

func init() {
	migrationMap = make(map[string]Migrationer)
}

// the basic type which will implement the basic type
type Migration struct {
	sqls    []string
	Created string
}

// implement in the Inheritance struct for upgrade
func (m *Migration) Up() {

}

// implement in the Inheritance struct for down
func (m *Migration) Down() {

}

// add sql want to execute
func (m *Migration) Sql(sql string) {
	m.sqls = append(m.sqls, sql)
}

// execute the sql already add in the sql
func (m *Migration) Exec(name, status string) error {
	o := orm.NewOrm()
	for _, s := range m.sqls {
		beego.Info("exec sql:", s)
		r := o.Raw(s)
		_, err := r.Exec()
		if err != nil {
			return err
		}
	}
	return m.addOrUpdateRecord(name, status)
}

func (m *Migration) addOrUpdateRecord(name, status string) error {
	o := orm.NewOrm()
	if status == "down" {
		status = "rollback"
		p, err := o.Raw("update migrations set status = ?, rollback_statements = ? where name = ?").Prepare()
		if err != nil {
			return nil
		}
		_, err = p.Exec(status, strings.Join(m.sqls, "; "), name)
		return err
	} else {
		status = "update"
		p, err := o.Raw("insert into migrations(`name`, `created_at`, `statements`, `status`) values(?,?,?,?)").Prepare()
		if err != nil {
			return err
		}
		_, err = p.Exec(name, strconv.FormatInt(m.GetCreated(), 10), strings.Join(m.sqls, "; "), status)
		return err
	}
}

// get the unixtime from the Created
func (m *Migration) GetCreated() int64 {
	t, err := time.Parse(M_DATE_FORMAT, m.Created)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// register the Migration in the map
func Register(name string, m Migrationer) error {
	if _, ok := migrationMap[name]; ok {
		return errors.New("already exist name:" + name)
	}
	migrationMap[name] = m
	return nil
}

// upgrate the migration from lasttime
func Upgrade(lasttime int64) error {
	sm := sortMap(migrationMap)
	i := 0
	for _, v := range sm {
		if v.created > lasttime {
			beego.Info("start upgrade", v.name)
			v.m.Up()
			err := v.m.Exec(v.name, "up")
			if err != nil {
				beego.Error("execute error:", err)
				return err
			}
			beego.Info("end upgrade:", v.name)
			i++
		}
	}
	beego.Info("total success upgrade:", i, " migration")
	return nil
}

//rollback the migration by the name
func Rollback(name string) error {
	if v, ok := migrationMap[name]; ok {
		beego.Info("start rollback")
		v.Down()
		err := v.Exec(name, "down")
		if err != nil {
			beego.Error("execute error:", err)
			return err
		}
		beego.Info("end rollback")
		return nil
	} else {
		return errors.New("not exist the migrationMap name:" + name)
	}
}

// reset all migration
// run all migration's down function
func Reset() error {
	i := 0
	for k, v := range migrationMap {
		beego.Info("start reset:", k)
		v.Down()
		err := v.Exec(k, "down")
		if err != nil {
			beego.Error("execute error:", err)
			return err
		}
		beego.Info("end reset:", k)
	}
	beego.Info("total success reset:", i, " migration")
	return nil
}

// first Reset, then Upgrade
func Refresh() error {
	err := Reset()
	if err != nil {
		beego.Error("execute error:", err)
		return err
	}
	return Upgrade(0)
}

type dataSlice []data

type data struct {
	created int64
	name    string
	m       Migrationer
}

// Len is part of sort.Interface.
func (d dataSlice) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d dataSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d dataSlice) Less(i, j int) bool {
	return d[i].created < d[j].created
}

func sortMap(m map[string]Migrationer) dataSlice {
	s := make(dataSlice, 0, len(m))
	for k, v := range m {
		d := data{}
		d.created = v.GetCreated()
		d.name = k
		d.m = v
		s = append(s, d)
	}
	sort.Sort(s)
	return s
}
