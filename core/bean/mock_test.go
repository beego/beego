package bean

import (
	"fmt"
	"testing"
)

func TestMock(t *testing.T) {
	type MockSubSubObject struct {
		A int `mock:"20"`
	}
	type MockSubObjectAnoy struct {
		Anoy int `mock:"20"`
	}
	type MockSubObject struct {
		A bool `mock:"true"`
		B MockSubSubObject
	}
	type MockObject struct {
		A string           `mock:"aaaaa"`
		B int8             `mock:"10"`
		C []*MockSubObject `mock:"length:2"`
		D bool             `mock:"true"`
		E *MockSubObject
		F []int `mock:"length:3"`
		G InterfaceA
		H InterfaceA
		MockSubObjectAnoy
	}
	m := &MockObject{G: &ImplA{}}
	err := Mock(m)
	if err != nil {
		t.Fatalf("mock failed: %v", err)
	}
	if m.A != "aaaaa" || m.B != 10 || m.C[1].B.A != 20 ||
		!m.E.A || m.E.B.A != 20 || !m.D || len(m.F) != 3 {
		t.Fail()
	}
	_, ok := m.G.(*ImplA)
	if !ok {
		t.Fail()
	}
	_, ok = m.G.(*ImplB)
	if ok {
		t.Fail()
	}
	_, ok = m.H.(*ImplA)
	if ok {
		t.Fail()
	}
	if m.Anoy != 20 {
		t.Fail()
	}
}

type InterfaceA interface {
	Item()
}

type ImplA struct {
	A string `mock:"aaa"`
}

func (i *ImplA) Item() {
	fmt.Println("implA")
}

type ImplB struct {
	B string `mock:"bbb"`
}

func (i *ImplB) Item() {
	fmt.Println("implB")
}
