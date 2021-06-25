package order_clause

import (
	"testing"
)

func TestClause(t *testing.T) {
	column := `a`

	o := Clause(
		Column(column),
	)

	if o.GetColumn() != column {
		t.Error()
	}
}

func TestSortAscending(t *testing.T) {
	o := Clause(
		SortAscending(),
	)

	if o.GetSort() != Ascending {
		t.Error()
	}
}

func TestSortDescending(t *testing.T) {
	o := Clause(
		SortDescending(),
	)

	if o.GetSort() != Descending {
		t.Error()
	}
}

func TestSortNone(t *testing.T) {
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

func TestRaw(t *testing.T) {
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

func TestColumn(t *testing.T) {
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

func TestOrderGetColumn(t *testing.T) {
	o := Clause(
		Column(`user__id`),
	)
	if o.GetColumn() != `user.id` {
		t.Error()
	}
}

func TestSortString(t *testing.T) {
	template := "got: %s, want: %s"

	o1 := Clause(sort(Sort(1)))
	if o1.SortString() != "ASC" {
		t.Errorf(template, o1.SortString(), "ASC")
	}

	o2 := Clause(sort(Sort(2)))
	if o2.SortString() != "DESC" {
		t.Errorf(template, o2.SortString(), "DESC")
	}

	o3 := Clause(sort(Sort(3)))
	if o3.SortString() != `` {
		t.Errorf(template, o3.SortString(), ``)
	}
}
