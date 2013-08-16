package orm

import (
	"fmt"
	"reflect"
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

func (o *querySet) setOffset(num interface{}) {
	val := reflect.ValueOf(num)
	switch num.(type) {
	case int, int8, int16, int32, int64:
		o.offset = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		o.offset = int64(val.Uint())
	default:
		panic(fmt.Errorf("<QuerySeter> offset value need numeric not `%T`", num))
	}
}

func (o querySet) Limit(limit int, args ...interface{}) QuerySeter {
	o.limit = limit
	if len(args) > 0 {
		o.setOffset(args[0])
	}
	return &o
}

func (o querySet) Offset(offset interface{}) QuerySeter {
	o.setOffset(offset)
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
	return o.orm.alias.DbBaser.Count(o.orm.db, o, o.mi, o.cond, o.orm.alias.TZ)
}

func (o *querySet) Update(values Params) (int64, error) {
	return o.orm.alias.DbBaser.UpdateBatch(o.orm.db, o, o.mi, o.cond, values, o.orm.alias.TZ)
}

func (o *querySet) Delete() (int64, error) {
	return o.orm.alias.DbBaser.DeleteBatch(o.orm.db, o, o.mi, o.cond, o.orm.alias.TZ)
}

func (o *querySet) PrepareInsert() (Inserter, error) {
	return newInsertSet(o.orm, o.mi)
}

func (o *querySet) All(container interface{}) (int64, error) {
	return o.orm.alias.DbBaser.ReadBatch(o.orm.db, o, o.mi, o.cond, container, o.orm.alias.TZ)
}

func (o *querySet) One(container interface{}) error {
	num, err := o.orm.alias.DbBaser.ReadBatch(o.orm.db, o, o.mi, o.cond, container, o.orm.alias.TZ)
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
	return o.orm.alias.DbBaser.ReadValues(o.orm.db, o, o.mi, o.cond, exprs, results, o.orm.alias.TZ)
}

func (o *querySet) ValuesList(results *[]ParamsList, exprs ...string) (int64, error) {
	return o.orm.alias.DbBaser.ReadValues(o.orm.db, o, o.mi, o.cond, exprs, results, o.orm.alias.TZ)
}

func (o *querySet) ValuesFlat(result *ParamsList, expr string) (int64, error) {
	return o.orm.alias.DbBaser.ReadValues(o.orm.db, o, o.mi, o.cond, []string{expr}, result, o.orm.alias.TZ)
}

func newQuerySet(orm *orm, mi *modelInfo) QuerySeter {
	o := new(querySet)
	o.mi = mi
	o.orm = orm
	return o
}
