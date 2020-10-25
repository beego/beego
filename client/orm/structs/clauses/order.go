package clauses

type Sort int8

const (
	ASCENDING  Sort = 1
	DESCENDING Sort = 2
)

type Order struct {
	column string
	sort   Sort
}

func (o *Order) GetColumn() string {
	return o.column
}

func (o *Order) GetSort() Sort {
	return o.sort
}

func ParseOrder(expressions ...string) []*Order {
	var orders []*Order
	for _, expression := range expressions {
		sort := ASCENDING
		column := expression
		if expression[0] == '-' {
			sort = DESCENDING
			column = expression[1:]
		}

		orders = append(orders, &Order{
			column: column,
			sort:   sort,
		})
	}

	return orders
}
