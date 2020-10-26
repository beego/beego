package clauses

import "testing"

func TestOrderClause(t *testing.T) {
	var (
		column = `a`
	)

	o := OrderClause(
		OrderColumn(column),
	)

	if o.GetColumn() != column {
		t.Error()
	}
}

func TestOrderSortAscending(t *testing.T) {
	o := OrderClause(
		OrderSortAscending(),
	)

	if o.GetSort() != SortAscending {
		t.Error()
	}
}

func TestOrderSortDescending(t *testing.T) {
	o := OrderClause(
		OrderSortDescending(),
	)

	if o.GetSort() != SortDescending {
		t.Error()
	}
}

func TestOrderSortNone(t *testing.T) {
	o1 := OrderClause(
		OrderSortNone(),
	)

	if o1.GetSort() != SortNone {
		t.Error()
	}

	o2 := OrderClause()

	if o2.GetSort() != SortNone {
		t.Error()
	}
}

func TestOrderRaw(t *testing.T) {
	o1 := OrderClause()

	if o1.IsRaw() {
		t.Error()
	}

	o2 := OrderClause(
		OrderRaw(),
	)

	if !o2.IsRaw() {
		t.Error()
	}
}

func TestOrderColumn(t *testing.T) {
	o1 := OrderClause(
		OrderColumn(`aaa`),
	)

	if o1.GetColumn() != `aaa` {
		t.Error()
	}
}

