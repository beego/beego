package orm

type dbBaseOracle struct {
	dbBase
}

var _ dbBaser = new(dbBaseOracle)

func newdbBaseOracle() dbBaser {
	b := new(dbBaseOracle)
	b.ins = b
	return b
}
