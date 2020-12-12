package cache

import (
	"fmt"
	"math"
)

func incr(originVal interface{}) (interface{}, error) {
	switch val := originVal.(type) {
	case int:
		tmp := val + 1
		if val > 0 && tmp < 0 {
			return nil, fmt.Errorf("increment would overflow")
		}
		return tmp, nil
	case int32:
		if val == math.MaxInt32 {
			return nil, fmt.Errorf("increment would overflow")
		}
		return val + 1, nil
	case int64:
		// if val == math.MaxInt64 {
		// 	return nil, fmt.Errorf("increment would overflow")
		// }
		return val + 1, nil
	case uint:
		tmp := val + 1
		if tmp < val {
			return nil, fmt.Errorf("increment would overflow")
		}
		return tmp, nil
	case uint32:
		if val == math.MaxUint32 {
			return nil, fmt.Errorf("increment would overflow")
		}
		return val + 1, nil
	case uint64:
		if val == math.MaxUint64 {
			return nil, fmt.Errorf("increment would overflow")
		}
		return val + 1, nil
	default:
		return nil, fmt.Errorf("item val is not (u)int (u)int32 (u)int64")
	}
}

func decr(originVal interface{}) (interface{}, error) {
	switch val := originVal.(type) {
	case int:
		tmp := val - 1
		if val < 0 && tmp > 0 {
			return nil, fmt.Errorf("decrement would overflow")
		}
		return tmp, nil
	case int32:
		if val == math.MinInt32 {
			return nil, fmt.Errorf("decrement would overflow")
		}
		return val - 1, nil
	case int64:
		if val == math.MinInt64 {
			return nil, fmt.Errorf("decrement would overflow")
		}
		return val - 1, nil
	case uint:
		if val == 0 {
			return nil, fmt.Errorf("decrement would overflow")
		}
		return val - 1, nil
	case uint32:
		if val == 0 {
			return nil, fmt.Errorf("increment would overflow")
		}
		return val - 1, nil
	case uint64:
		if val == 0 {
			return nil, fmt.Errorf("increment would overflow")
		}
		return val - 1, nil
	default:
		return nil, fmt.Errorf("item val is not (u)int (u)int32 (u)int64")
	}
}