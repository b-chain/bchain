package wasmre

import (
	"reflect"
	"testing"
	"fmt"
)

type api struct {
	num *int
}

func (t *api) f() {
	fmt.Println(*t.num)
}


type Bridge struct {
	m *moudule
	num *int
	api *api
}

func (t *Bridge) load() {
	f := func() {
		t.api.f()
		fmt.Println(*t.api.num)
	}
	t.m.f = reflect.ValueOf(f)
}

func load(t *Bridge) interface{} {
	//a := api{t.num}
	//t.api = &a

	f := func() {
		//t.api.f()
		fmt.Println(*t.api.num)
	}
	//t.m.f = reflect.ValueOf(f)
	return f
}

type moudule struct {
	f reflect.Value
}


func TestClosePacket(t *testing.T) {
	ii := 1
	a := &api{&ii}
	m := &moudule{}
	b := Bridge{m, &ii,a}
	b.load()
	//f := load(&b)
	//b.m.f = reflect.ValueOf(f)

	//f := func() {
	//	fmt.Println(*a.num)
	//}
	//b.m.f = reflect.ValueOf(f)
	m.f.Call([]reflect.Value{})

	iii := 3
	a.num = &iii
	m.f.Call([]reflect.Value{})

	aa := &api{&iii}
	bb := Bridge{m, &iii,aa}
	bb.load()

	m.f.Call([]reflect.Value{})

}

func TestXxx(t *testing.T)  {
	xx := 1750 * DecaimasBase
	fmt.Println(xx)
}
