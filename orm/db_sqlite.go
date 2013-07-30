package orm

type dbBaseSqlite struct {
	dbBase
}

func newdbBaseSqlite() dbBaser {
	b := new(dbBaseSqlite)
	b.ins = b
	return b
}
