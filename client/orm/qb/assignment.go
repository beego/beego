package qb

// Assignable represents that something could be used alias "assignment" statement
type Assignable interface {
	assign()
}

type Assignment binaryExpr

func Assign(column string, value interface{}) Assignment {
	var expr Expression
	switch v := value.(type) {
	case Expression:
		expr = v
	default:
		expr = valueExpr{val: v}
	}
	return Assignment{left: C(column), op: opEqual, right: expr}
}

func (Assignment) assign() {}
