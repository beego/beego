package order

import (
	"testing"
)

func TestOrderClause(t *testing.T) {
	var (
		column = `a`
	)

	o := Clause(
		Column(column),
	)

	if o.GetColumn() != column {
		t.Error()
	}
}

func TestOrderSortAscending(t *testing.T) {
	o := Clause(
		SortAscending(),
	)

	if o.GetSort() != Ascending {
		t.Error()
	}
}

func TestOrderSortDescending(t *testing.T) {
	o := Clause(
		SortDescending(),
	)

	if o.GetSort() != Descending {
		t.Error()
	}
}

func TestOrderSortNone(t *testing.T) {
	o1 := Clause(
		SortNone(),
	)

	if o1.GetSort() != None {
		t.Error()
	}

	o2 := Clause()

	if o2.GetSort() != None {
		t.Error()
	}
}

func TestOrderRaw(t *testing.T) {
	o1 := Clause()

	if o1.IsRaw() {
		t.Error()
	}

	o2 := Clause(
		Raw(),
	)

	if !o2.IsRaw() {
		t.Error()
	}
}

func TestOrderColumn(t *testing.T) {
	o1 := Clause(
		Column(`aaa`),
	)

	if o1.GetColumn() != `aaa` {
		t.Error()
	}
}

func TestParseOrder(t *testing.T) {
	orders := ParseOrder(
		`-user__status`,
		`status`,
		`user__status`,
	)

	t.Log(orders)

	if orders[0].GetSort() != Descending {
		t.Error()
	}

	if orders[0].GetColumn() != `user.status` {
		t.Error()
	}

	if orders[1].GetColumn() != `status` {
		t.Error()
	}

	if orders[1].GetSort() != Ascending {
		t.Error()
	}

	if orders[2].GetColumn() != `user.status` {
		t.Error()
	}

}
