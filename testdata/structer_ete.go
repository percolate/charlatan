package main

import "fmt"

var _ Structer = &FakeStructer{}

func main() {
	structHookCalled := false
	namedStructHookCalled := false

	s := struct {
		a string
		b string
	}{"e", "f"}
	t := struct {
		a string
		b string
	}{"g", "h"}
	u := struct {
		c string
		d string
	}{"i", "j"}
	v := struct {
		c string
		d string
	}{"k", "l"}

	f := &FakeStructer{
		StructHook: func(a struct {
			a string
			b string
		}) struct {
			c string
			d string
		} {
			structHookCalled = true
			return v
		},
		NamedStructHook: func(z struct {
			a string
			b string
		}) (a struct {
			c string
			d string
		}) {
			namedStructHookCalled = true
			if z == s {
				return u
			}
			return v
		},
	}

	f.Struct(s)

	if len(f.StructCalls) != 1 {
		panic(fmt.Sprintf("StructCalls: %d", len(f.StructCalls)))
	}
	if !structHookCalled {
		panic("StructHook not called")
	}
	if !f.StructCalled() {
		panic("StructCalled: Struct not called")
	}
	if !f.StructCalledOnce() {
		panic("StructCalledOnce: Struct not called once")
	}
	if f.StructNotCalled() {
		panic("StructNotCalled: Struct not called")
	}
	if !f.StructCalledN(1) {
		panic("StructCalledN: Struct not called once")
	}
	if !f.StructCalledWith(s) {
		panic(fmt.Sprintf("StructCalledWith: Struct not called with %s", s))
	}
	if !f.StructCalledOnceWith(s) {
		panic(fmt.Sprintf("StructCalledOnceWith: Struct not called once with %s", s))
	}

	f.NamedStruct(s)

	if len(f.NamedStructCalls) != 1 {
		panic(fmt.Sprintf("NamedStructCalls: %d", len(f.NamedStructCalls)))
	}
	if !namedStructHookCalled {
		panic("NamedStructHook not called")
	}
	if !f.NamedStructCalled() {
		panic("NamedStructCalled: NamedStruct not called")
	}
	if !f.NamedStructCalledOnce() {
		panic("NamedStructCalledOnce: NamedStruct not called once")
	}
	if f.NamedStructNotCalled() {
		panic("NamedStructNotCalled: NamedStruct not called")
	}
	if !f.NamedStructCalledN(1) {
		panic("NamedStructCalledN: NamedStruct not called once")
	}
	if !f.NamedStructCalledWith(s) {
		panic(fmt.Sprintf("NamedStructCalledWith: NamedStruct not called once with %s", s))
	}
	if !f.NamedStructCalledOnceWith(s) {
		panic(fmt.Sprintf("NamedStructCalledOnceWith: NamedStruct not called once with %s", s))
	}

	res, found := f.NamedStructResultsForCall(s)
	if res != u || found != true {
		panic(fmt.Sprintf("NamedStructResultsForCall: NamedStruct results for %s not %s, found: %s", s, u, found))
	}

	res, found = f.NamedStructResultsForCall(t)
	if found != false {
		panic(fmt.Sprintf("NamedStructResultsForCall: NamedStruct results for %s found", t))
	}

	f.NamedStruct(t)

	if len(f.NamedStructCalls) != 2 {
		panic(fmt.Sprintf("NamedStructCalls: %d", len(f.NamedStructCalls)))
	}

	if !f.NamedStructCalledN(2) {
		panic("NamedStructCalledN: NamedStruct not called twice")
	}

	res, found = f.NamedStructResultsForCall(t)
	if res != v || found != true {
		panic(fmt.Sprintf("NamedStructResultsForCall: NamedStruct results for %s not %s, found: %s", t, v, found))
	}
}
