package orm

import (
	"fmt"
)

type querySet struct {
	mi       *modelInfo
	cond     *Condition
	related  []string
	relDepth int
	limit    int
	offset   int64
	orders   []string
	orm      *orm
}

var _ QuerySeter = new(querySet)

func (o querySet) Filter(expr string, args ...interface{}) QuerySeter {
	if o.cond == nil {
		o.cond = NewCondition()
	}
	o.cond = o.cond.And(expr, args...)
	return &o
}

func (o querySet) Exclude(expr string, args ...interface{}) QuerySeter {
	if o.cond == nil {
		o.cond = NewCondition()
	}
	o.cond = o.cond.AndNot(expr, args...)
	return &o
}

func (o querySet) Limit(limit int, args ...int64) QuerySeter {
	o.limit = limit
	if len(args) > 0 {
		o.offset = args[0]
	}
	return &o
}

func (o querySet) Offset(offset int64) QuerySeter {
	o.offset = offset
	return &o
}

func (o querySet) OrderBy(exprs ...string) QuerySeter {
	o.orders = exprs
	return &o
}

func (o querySet) RelatedSel(params ...interface{}) QuerySeter {
	var related []string
	if len(params) == 0 {
		o.relDepth = DefaultRelsDepth
	} else {
		for _, p := range params {
			switch val := p.(type) {
			case string:
				related = append(o.related, val)
			case int:
				o.relDepth = val
			default:
				panic(fmt.Sprintf("<QuerySeter.RelatedSel> wrong param kind: %v", val))
			}
		}
	}
	o.related = related
	return &o
}

func (o querySet) SetCond(cond *Condition) QuerySeter {
	o.cond = cond
	return &o
}

func (o *querySet) Count() (int64, error) {
	return o.orm.alias.DbBaser.Count(o.orm.db, o, o.mi, o.cond)
}

func (o *querySet) Update(values Params) (int64, error) {
	return o.orm.alias.DbBaser.UpdateBatch(o.orm.db, o, o.mi, o.cond, values)
}

func (o *querySet) Delete() (int64, error) {
	return o.orm.alias.DbBaser.DeleteBatch(o.orm.db, o, o.mi, o.cond)
}

func (o *querySet) PrepareInsert() (Inserter, error) {
	return newInsertSet(o.orm, o.mi)
}

func (o *querySet) All(container interface{}) (int64, error) {
	return o.orm.alias.DbBaser.ReadBatch(o.orm.db, o, o.mi, o.cond, container)
}

func (o *querySet) One(container interface{}) error {
	num, err := o.orm.alias.DbBaser.ReadBatch(o.orm.db, o, o.mi, o.cond, container)
	if err != nil {
		return err
	}
	if num > 1 {
		return ErrMultiRows
	}
	if num == 0 {
		return ErrNoRows
	}
	return nil
}

func (o *querySet) Values(results *[]Params, exprs ...string) (int64, error) {
	return o.orm.alias.DbBaser.ReadValues(o.orm.db, o, o.mi, o.cond, exprs, results)
}

func (o *querySet) ValuesList(results *[]ParamsList, exprs ...string) (int64, error) {
	return o.orm.alias.DbBaser.ReadValues(o.orm.db, o, o.mi, o.cond, exprs, results)
}

func (o *querySet) ValuesFlat(result *ParamsList, expr string) (int64, error) {
	return o.orm.alias.DbBaser.ReadValues(o.orm.db, o, o.mi, o.cond, []string{expr}, result)
}

func newQuerySet(orm *orm, mi *modelInfo) QuerySeter {
	o := new(querySet)
	o.mi = mi
	o.orm = orm
	return o
}
