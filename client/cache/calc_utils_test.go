package cache

import (
	"math"
	"strconv"
	"testing"
)

func TestIncr(t *testing.T) {
	// int
	var originVal interface{} = int(1)
	var updateVal interface{} = int(2)
	val, err := incr(originVal)
	if err != nil {
		t.Errorf("incr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("incr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = incr(int(1 << (strconv.IntSize - 1) - 1))
	if err == nil {
		t.Error("incr failed")
		return
	}

	// int32
	originVal = int32(1)
	updateVal = int32(2)
	val, err = incr(originVal)
	if err != nil {
		t.Errorf("incr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("incr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = incr(int32(math.MaxInt32))
	if err == nil {
		t.Error("incr failed")
		return
	}

	// int64
	originVal = int64(1)
	updateVal = int64(2)
	val, err = incr(originVal)
	if err != nil {
		t.Errorf("incr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("incr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = incr(int64(math.MaxInt64))
	if err == nil {
		t.Error("incr failed")
		return
	}

	// uint
	originVal = uint(1)
	updateVal = uint(2)
	val, err = incr(originVal)
	if err != nil {
		t.Errorf("incr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("incr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = incr(uint(1 << (strconv.IntSize) - 1))
	if err == nil {
		t.Error("incr failed")
		return
	}

	// uint32
	originVal = uint32(1)
	updateVal = uint32(2)
	val, err = incr(originVal)
	if err != nil {
		t.Errorf("incr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("incr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = incr(uint32(math.MaxUint32))
	if err == nil {
		t.Error("incr failed")
		return
	}

	// uint64
	originVal = uint64(1)
	updateVal = uint64(2)
	val, err = incr(originVal)
	if err != nil {
		t.Errorf("incr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("incr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = incr(uint64(math.MaxUint64))
	if err == nil {
		t.Error("incr failed")
		return
	}

	// other type
	_, err = incr("string")
	if err == nil {
		t.Error("incr failed")
		return
	}
}

func TestDecr(t *testing.T) {
	// int
	var originVal interface{} = int(2)
	var updateVal interface{} = int(1)
	val, err := decr(originVal)
	if err != nil {
		t.Errorf("decr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("decr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = decr(int(-1 << (strconv.IntSize - 1)))
	if err == nil {
		t.Error("decr failed")
		return
	}

	// int32
	originVal = int32(2)
	updateVal = int32(1)
	val, err = decr(originVal)
	if err != nil {
		t.Errorf("decr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("decr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = decr(int32(math.MinInt32))
	if err == nil {
		t.Error("decr failed")
		return
	}

	// int64
	originVal = int64(2)
	updateVal = int64(1)
	val, err = decr(originVal)
	if err != nil {
		t.Errorf("decr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("decr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = decr(int64(math.MinInt64))
	if err == nil {
		t.Error("decr failed")
		return
	}

	// uint
	originVal = uint(2)
	updateVal = uint(1)
	val, err = decr(originVal)
	if err != nil {
		t.Errorf("decr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("decr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = decr(uint(0))
	if err == nil {
		t.Error("decr failed")
		return
	}

	// uint32
	originVal = uint32(2)
	updateVal = uint32(1)
	val, err = decr(originVal)
	if err != nil {
		t.Errorf("decr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("decr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = decr(uint32(0))
	if err == nil {
		t.Error("decr failed")
		return
	}

	// uint64
	originVal = uint64(2)
	updateVal = uint64(1)
	val, err = decr(originVal)
	if err != nil {
		t.Errorf("decr failed, err: %s", err.Error())
		return
	}
	if val != updateVal {
		t.Errorf("decr failed, expect %v, but %v actually", updateVal, val)
		return
	}
	_, err = decr(uint64(0))
	if err == nil {
		t.Error("decr failed")
		return
	}

	// other type
	_, err = decr("string")
	if err == nil {
		t.Error("decr failed")
		return
	}
}