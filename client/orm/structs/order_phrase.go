package structs

import "fmt"

type Sort int8

const (
	ASCENDING  Sort = 1
	DESCENDING Sort = 2
)

type OrderClause struct {
	Column string
	Sort   Sort
}

var _ fmt.Stringer = new(OrderClause)

func (o *OrderClause) String() string {
	sort := ``
	if o.Sort == ASCENDING {
		sort = `ASC`
	} else if o.Sort == DESCENDING {
		sort = `DESC`
	} else {
		return fmt.Sprintf("%s", o.Column)
	}
	return fmt.Sprintf("%s %s", o.Column, sort)
}

func ParseOrderClause(expressions ...string) []*OrderClause {
	var orders []*OrderClause
	for _, expression := range expressions {
		sort := ASCENDING
		column := expression
		if expression[0] == '-' {
			sort = DESCENDING
			column = expression[1:]
		}

		orders = append(orders, &OrderClause{
			Column: column,
			Sort:   sort,
		})
	}

	return orders
}
