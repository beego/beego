package orm

import (
	"reflect"
	"testing"
)

type GetAllFieldT1 struct {
	GetAllFieldT3
	*GetAllFieldT4
	B int
}

type GetAllFieldT2 struct {
	A int
	B int
	C int
}

type GetAllFieldT3 struct {
	A int
	B int
	GetAllFieldT2
}
type GetAllFieldT4 struct {
	A int
	D int
}

type GetAllFieldT5 struct {
	GetAllFieldT6
	A int
}
type GetAllFieldT6 int

func TestStructGetAllField(t *testing.T) {
	t1 := reflect.TypeOf(&GetAllFieldT1{})
	ret := structGetAllField(t1)
	throwFail(t, AssertIs(len(ret), 7))
	throwFail(t, AssertIs(ret[0].Name, "GetAllFieldT3"))
	throwFail(t, AssertIs(ret[1].Name, "GetAllFieldT4"))
	throwFail(t, AssertIs(ret[2].Name, "B"))
	throwFail(t, AssertIs(ret[2].Index, []int{2}))
	throwFail(t, AssertIs(ret[3].Name, "A"))
	throwFail(t, AssertIs(ret[3].Index, []int{0, 0}))
	throwFail(t, AssertIs(ret[4].Name, "GetAllFieldT2"))
	throwFail(t, AssertIs(ret[5].Name, "C"))
	throwFail(t, AssertIs(ret[5].Index, []int{0, 2, 2}))
	throwFail(t, AssertIs(ret[6].Name, "D"))
	throwFail(t, AssertIs(ret[6].Index, []int{1, 1}))

	ret = structGetAllField(reflect.TypeOf(&GetAllFieldT5{}))
	throwFail(t, AssertIs(len(ret), 2))

}
