package orm

type dbBaseMysql struct {
	dbBase
}

func (d *dbBaseMysql) GetOperatorSql(mi *modelInfo, operator string, args []interface{}) (sql string, params []interface{}) {
	return d.dbBase.GetOperatorSql(mi, operator, args)
}

func newdbBaseMysql() dbBaser {
	b := new(dbBaseMysql)
	b.ins = b
	return b
}
