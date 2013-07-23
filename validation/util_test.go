package validation

import (
	"reflect"
	"testing"
)

type user struct {
	Id   int
	Tag  string `valid:"Maxx(aa)"`
	Name string `valid:"Required"`
	Age  int    `valid:"Required;Range(1, 140)"`
}

func TestGetValidFuncs(t *testing.T) {
	u := user{Name: "test", Age: 1}
	tf := reflect.TypeOf(u)
	var vfs []ValidFunc
	var err error

	f, _ := tf.FieldByName("Id")
	if vfs, err = getValidFuncs(f); err != nil {
		t.Fatal(err)
	}
	if len(vfs) != 0 {
		t.Fatal("should get none ValidFunc")
	}

	f, _ = tf.FieldByName("Tag")
	if vfs, err = getValidFuncs(f); err.Error() != "doesn't exsits Maxx valid function" {
		t.Fatal(err)
	}

	f, _ = tf.FieldByName("Name")
	if vfs, err = getValidFuncs(f); err != nil {
		t.Fatal(err)
	}
	if len(vfs) != 1 {
		t.Fatal("should get 1 ValidFunc")
	}
	if vfs[0].Name != "Required" && len(vfs[0].Params) != 0 {
		t.Error("Required funcs should be got")
	}

	f, _ = tf.FieldByName("Age")
	if vfs, err = getValidFuncs(f); err != nil {
		t.Fatal(err)
	}
	if len(vfs) != 2 {
		t.Fatal("should get 2 ValidFunc")
	}
	if vfs[0].Name != "Required" && len(vfs[0].Params) != 0 {
		t.Error("Required funcs should be got")
	}
	if vfs[1].Name != "Range" && len(vfs[1].Params) != 2 {
		t.Error("Range funcs should be got")
	}
}
