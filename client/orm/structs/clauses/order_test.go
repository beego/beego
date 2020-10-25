package clauses

import "testing"

func TestOrderClause(t *testing.T) {
	var (
		column = `a`
		sort   = SortDescending
		raw    = true
	)

	o := OrderClause(
		OrderColumn(column),
		OrderSort(sort),
		OrderRaw(raw),
	)

	if o.GetColumn() != column {
		t.Error()
	}

	if o.GetSort() != sort {
		t.Error()
	}

	if o.IsRaw() != raw {
		t.Error()
	}
}
