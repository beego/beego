package cache

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncr(t *testing.T) {
	// int
	var originVal interface{} = int(1)
	var updateVal interface{} = int(2)
	val, err := incr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = incr(int(1<<(strconv.IntSize-1) - 1))
	assert.Equal(t, ErrIncrementOverflow, err)

	// int32
	originVal = int32(1)
	updateVal = int32(2)
	val, err = incr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = incr(int32(math.MaxInt32))
	assert.Equal(t, ErrIncrementOverflow, err)

	// int64
	originVal = int64(1)
	updateVal = int64(2)
	val, err = incr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = incr(int64(math.MaxInt64))
	assert.Equal(t, ErrIncrementOverflow, err)

	// uint
	originVal = uint(1)
	updateVal = uint(2)
	val, err = incr(originVal)

	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = incr(uint(1<<(strconv.IntSize) - 1))
	assert.Equal(t, ErrIncrementOverflow, err)

	// uint32
	originVal = uint32(1)
	updateVal = uint32(2)
	val, err = incr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = incr(uint32(math.MaxUint32))
	assert.Equal(t, ErrIncrementOverflow, err)

	// uint64
	originVal = uint64(1)
	updateVal = uint64(2)
	val, err = incr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = incr(uint64(math.MaxUint64))
	assert.Equal(t, ErrIncrementOverflow, err)
	// other type
	_, err = incr("string")
	assert.Equal(t, ErrNotIntegerType, err)
}

func TestDecr(t *testing.T) {
	// int
	var originVal interface{} = int(2)
	var updateVal interface{} = int(1)
	val, err := decr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = decr(int(-1 << (strconv.IntSize - 1)))
	assert.Equal(t, ErrDecrementOverflow, err)
	// int32
	originVal = int32(2)
	updateVal = int32(1)
	val, err = decr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = decr(int32(math.MinInt32))
	assert.Equal(t, ErrDecrementOverflow, err)

	// int64
	originVal = int64(2)
	updateVal = int64(1)
	val, err = decr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = decr(int64(math.MinInt64))
	assert.Equal(t, ErrDecrementOverflow, err)

	// uint
	originVal = uint(2)
	updateVal = uint(1)
	val, err = decr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = decr(uint(0))
	assert.Equal(t, ErrDecrementOverflow, err)

	// uint32
	originVal = uint32(2)
	updateVal = uint32(1)
	val, err = decr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = decr(uint32(0))
	assert.Equal(t, ErrDecrementOverflow, err)

	// uint64
	originVal = uint64(2)
	updateVal = uint64(1)
	val, err = decr(originVal)
	assert.Nil(t, err)
	assert.Equal(t, val, updateVal)

	_, err = decr(uint64(0))
	assert.Equal(t, ErrDecrementOverflow, err)

	// other type
	_, err = decr("string")
	assert.Equal(t, ErrNotIntegerType, err)
}
