package orm

type dbBaseOracle struct {
	dbBase
}

func newdbBaseOracle() dbBaser {
	b := new(dbBaseOracle)
	b.ins = b
	return b
}
