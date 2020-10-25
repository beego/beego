package clauses

import "strings"

type Sort int8

const (
	SortNone       Sort = 0
	SortAscending  Sort = 1
	SortDescending Sort = 2
)

type OrderOption func(order *Order)

type Order struct {
	column string
	sort   Sort
	isRaw  bool
}

func OrderClause(options ...OrderOption) *Order {
	o := &Order{}
	for _, option := range options {
		option(o)
	}

	return o
}

func (o *Order) GetColumn() string {
	return o.column
}

func (o *Order) GetSort() Sort {
	return o.sort
}

func (o *Order) SortString() string {
	switch o.GetSort() {
	case SortAscending:
		return "ASC"
	case SortDescending:
		return "DESC"
	}

	return ``
}


func (o *Order) IsRaw() bool {
	return o.isRaw
}

func ParseOrder(expressions ...string) []*Order {
	var orders []*Order
	for _, expression := range expressions {
		sort := SortAscending
		column := strings.ReplaceAll(expression, ExprSep, ExprDot)
		if column[0] == '-' {
			sort = SortDescending
			column = column[1:]
		}

		orders = append(orders, &Order{
			column: column,
			sort:   sort,
		})
	}

	return orders
}

func OrderColumn(column string) OrderOption {
	return func(order *Order) {
		order.column = column
	}
}

func OrderSort(sort Sort) OrderOption {
	return func(order *Order) {
		order.sort = sort
	}
}

func OrderRaw(isRaw bool) OrderOption {
	return func(order *Order) {
		order.isRaw = isRaw
	}
}
