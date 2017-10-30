package main

import "fmt"

var _ Namedvaluer = &FakeNamedvaluer{}

func main() {
	namedHookCalled := false
	manyNamedHookCalled := false

	a := "one"
	b := "two"
	c := 3
	d := 4

	f := &FakeNamedvaluer{
		ManyNamedHook: func(a, b string, c, d int) bool {
			manyNamedHookCalled = true
			if c == 3 {
				return true
			}
			return false
		},
		NamedHook: func(a int, b string) bool {
			namedHookCalled = true
			if b == "one" {
				return true
			}
			return false
		},
	}

	res := f.Named(c, a)

	if len(f.NamedCalls) != 1 {
		panic(fmt.Sprintf("NamedCalls: %d", len(f.NamedCalls)))
	}
	if !namedHookCalled {
		panic("NamedHook not called")
	}
	if !f.NamedCalled() {
		panic("NamedCalled: Named not called")
	}
	if !f.NamedCalledOnce() {
		panic("NamedCalledOnce: Named not called once")
	}
	if f.NamedNotCalled() {
		panic("NamedNotCalled: Named not called")
	}
	if !f.NamedCalledN(1) {
		panic("NamedCalledN: Named not called once")
	}
	if !f.NamedCalledWith(c, a) {
		panic(fmt.Sprintf("NamedCalledWith: Named not called with %s, %s", c, a))
	}
	if !f.NamedCalledOnceWith(c, a) {
		panic(fmt.Sprintf("NamedCalledOnceWith: Named not called once with %s, %s", c, a))
	}

	res = f.ManyNamed(a, b, c, d)

	if len(f.ManyNamedCalls) != 1 {
		panic(fmt.Sprintf("ManyNamedCalls: %d", len(f.ManyNamedCalls)))
	}
	if !manyNamedHookCalled {
		panic("ManyNamedHook not called")
	}
	if !f.ManyNamedCalled() {
		panic("ManyNamedCalled: ManyNamed not called")
	}
	if !f.ManyNamedCalledOnce() {
		panic("ManyNamedCalledOnce: ManyNamed not called once")
	}
	if f.ManyNamedNotCalled() {
		panic("ManyNamedNotCalled: ManyNamed not called")
	}
	if !f.ManyNamedCalledN(1) {
		panic("ManyNamedCalledN: ManyNamed not called once")
	}
	if !f.ManyNamedCalledWith(a, b, c, d) {
		panic(fmt.Sprintf("ManyNamedCalledWith: ManyNamed not called once with %s, %s, %s", a, b, c, d))
	}
	if !f.ManyNamedCalledOnceWith(a, b, c, d) {
		panic(fmt.Sprintf("ManyNamedCalledOnceWith: ManyNamed not called once with %s, %s, %s", a, b, c, d))
	}

	res, found := f.ManyNamedResultsForCall(a, b, c, d)
	if res != true || found != true {
		panic(fmt.Sprintf("ManyNamedResultsForCall: ManyNamed results for %s, %s, %s not %s, found: %s", a, b, c, d, true, found))
	}

	res, found = f.ManyNamedResultsForCall(b, a, d, c)
	if found != false {
		panic(fmt.Sprintf("ManyNamedResultsForCall: ManyNamed results for %s, %s, %s, %s found", b, a, d, c))
	}

	f.ManyNamed(b, a, d, c)

	if len(f.ManyNamedCalls) != 2 {
		panic(fmt.Sprintf("ManyNamedCalls: %d", len(f.ManyNamedCalls)))
	}

	if !f.ManyNamedCalledN(2) {
		panic("ManyNamedCalledN: ManyNamed not called twice")
	}

	res, found = f.ManyNamedResultsForCall(b, a, d, c)
	if res != false || found != true {
		panic(fmt.Sprintf("ManyNamedResultsForCall: ManyNamed results for %s, %s, %s not %s, found: %s", b, a, d, c, false, found))
	}
}
