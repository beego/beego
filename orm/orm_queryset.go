package orm

import (
	"fmt"
)

type colValue struct {
	value int64
	opt   operator
}

type operator int

const (
	Col_Add operator = iota
	Col_Minus
	Col_Multiply
	Col_Except
)

func ColValue(opt operator, value interface{}) interface{} {
	switch opt {
	case Col_Add, Col_Minus, Col_Multiply, Col_Except:
	default:
		panic(fmt.Errorf("orm.ColValue wrong operator"))
	}
	v, err := StrTo(ToStr(value)).Int64()
	if err != nil {
		panic(fmt.Errorf("orm.ColValue doesn't support non string/numeric type, %s", err))
	}
	var val colValue
	val.value = v
	val.opt = opt
	return val
}

type querySet struct {
	mi       *modelInfo
	cond     *Condition
	related  []string
	relDepth int
	limit    int64
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
	o.offset = ToInt64(num)
}

func (o querySet) Limit(limit interface{}, args ...interface{}) QuerySeter {
	o.limit = ToInt64(limit)
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
				panic(fmt.Errorf("<QuerySeter.RelatedSel> wrong param kind: %v", val))
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

func (o *querySet) Exist() bool {
	cnt, _ := o.orm.alias.DbBaser.Count(o.orm.db, o, o.mi, o.cond, o.orm.alias.TZ)
	return cnt > 0
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

func (o *querySet) All(container interface{}, cols ...string) (int64, error) {
	return o.orm.alias.DbBaser.ReadBatch(o.orm.db, o, o.mi, o.cond, container, o.orm.alias.TZ, cols)
}

func (o *querySet) One(container interface{}, cols ...string) error {
	num, err := o.orm.alias.DbBaser.ReadBatch(o.orm.db, o, o.mi, o.cond, container, o.orm.alias.TZ, cols)
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
