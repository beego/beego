package orm

type dbBasePostgres struct {
	dbBase
}

func newdbBasePostgres() dbBaser {
	b := new(dbBasePostgres)
	b.ins = b
	return b
}
