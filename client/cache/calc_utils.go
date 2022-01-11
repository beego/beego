package cache

import (
	"math"

	"github.com/beego/beego/v2/core/berror"
)

var (
	ErrIncrementOverflow = berror.Error(IncrementOverflow, "this incr invocation will overflow.")
	ErrDecrementOverflow = berror.Error(DecrementOverflow, "this decr invocation will overflow.")
	ErrNotIntegerType    = berror.Error(NotIntegerType, "item val is not (u)int (u)int32 (u)int64")
)

const (
	MinUint32 uint32 = 0
	MinUint64 uint64 = 0
)

func incr(originVal interface{}) (interface{}, error) {
	switch val := originVal.(type) {
	case int:
		tmp := val + 1
		if val > 0 && tmp < 0 {
			return nil, ErrIncrementOverflow
		}
		return tmp, nil
	case int32:
		if val == math.MaxInt32 {
			return nil, ErrIncrementOverflow
		}
		return val + 1, nil
	case int64:
		if val == math.MaxInt64 {
			return nil, ErrIncrementOverflow
		}
		return val + 1, nil
	case uint:
		tmp := val + 1
		if tmp < val {
			return nil, ErrIncrementOverflow
		}
		return tmp, nil
	case uint32:
		if val == math.MaxUint32 {
			return nil, ErrIncrementOverflow
		}
		return val + 1, nil
	case uint64:
		if val == math.MaxUint64 {
			return nil, ErrIncrementOverflow
		}
		return val + 1, nil
	default:
		return nil, ErrNotIntegerType
	}
}

func decr(originVal interface{}) (interface{}, error) {
	switch val := originVal.(type) {
	case int:
		tmp := val - 1
		if val < 0 && tmp > 0 {
			return nil, ErrDecrementOverflow
		}
		return tmp, nil
	case int32:
		if val == math.MinInt32 {
			return nil, ErrDecrementOverflow
		}
		return val - 1, nil
	case int64:
		if val == math.MinInt64 {
			return nil, ErrDecrementOverflow
		}
		return val - 1, nil
	case uint:
		if val == 0 {
			return nil, ErrDecrementOverflow
		}
		return val - 1, nil
	case uint32:
		if val == MinUint32 {
			return nil, ErrDecrementOverflow
		}
		return val - 1, nil
	case uint64:
		if val == MinUint64 {
			return nil, ErrDecrementOverflow
		}
		return val - 1, nil
	default:
		return nil, ErrNotIntegerType
	}
}
