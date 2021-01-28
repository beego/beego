package order_clause

import (
	"strings"

	"github.com/beego/beego/v2/client/orm/clauses"
)

type Sort int8

const (
	None       Sort = 0
	Ascending  Sort = 1
	Descending Sort = 2
)

type Option func(order *Order)

type Order struct {
	column string
	sort   Sort
	isRaw  bool
}

func Clause(options ...Option) *Order {
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
	case Ascending:
		return "ASC"
	case Descending:
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
		sort := Ascending
		column := strings.ReplaceAll(expression, clauses.ExprSep, clauses.ExprDot)
		if column[0] == '-' {
			sort = Descending
			column = column[1:]
		}

		orders = append(orders, &Order{
			column: column,
			sort:   sort,
		})
	}

	return orders
}

func Column(column string) Option {
	return func(order *Order) {
		order.column = strings.ReplaceAll(column, clauses.ExprSep, clauses.ExprDot)
	}
}

func sort(sort Sort) Option {
	return func(order *Order) {
		order.sort = sort
	}
}

func SortAscending() Option {
	return sort(Ascending)
}

func SortDescending() Option {
	return sort(Descending)
}

func SortNone() Option {
	return sort(None)
}

func Raw() Option {
	return func(order *Order) {
		order.isRaw = true
	}
}
