package orm

type dbBaseOracle struct {
	dbBase
}

var _ dbBaser = new(dbBaseOracle)

func (d *dbBase) OperatorSql(operator string) string {
	return ""
}

func newdbBaseOracle() dbBaser {
	b := new(dbBaseOracle)
	b.ins = b
	return b
}
