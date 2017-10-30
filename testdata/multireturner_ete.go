package main

import (
	"fmt"
)

var _ Multireturner = &FakeMultireturner{}

func main() {
	multiReturnHookCalled := false
	namedHookCalled := false

	s, i := "one", 2
	a, b, c, d := 1, 2, 3, 4

	f := &FakeMultireturner{
		MultiReturnHook: func() (string, int) {
			multiReturnHookCalled = true
			return s, i
		},
		NamedReturnHook: func() (int, int, int, int) {
			namedHookCalled = true
			return a, b, c, d
		},
	}

	rs, ri := f.MultiReturn()

	if rs != s || ri != i {
		panic(fmt.Sprintf("Unexpected results from Multireturn: %s, %s (expected %s, %s)", rs, ri, s, i))
	}
	if len(f.MultiReturnCalls) != 1 {
		panic(fmt.Sprintf("MultiReturnCalls: %d", len(f.MultiReturnCalls)))
	}
	if !multiReturnHookCalled {
		panic("MultiReturnHook not called")
	}
	if !f.MultiReturnCalled() {
		panic("MultiReturnCalled: MultiReturn not called")
	}
	if !f.MultiReturnCalledOnce() {
		panic("MultiReturnCalledOnce: MultiReturn not called once")
	}
	if f.MultiReturnNotCalled() {
		panic("MultiReturnNotCalled: MultiReturn not called")
	}
	if !f.MultiReturnCalledN(1) {
		panic("MultiReturnCalledN: MultiReturn not called once")
	}

	ra, rb, rc, rd := f.NamedReturn()

	if ra != a ||  rb != b || rc != c || rd != d {
		panic(fmt.Sprintf("Unexpected results from Multireturn: %s, %s, %s, %s (expected %s, %s, %s, %s)", ra, rb, rc, rd, a, b, c, d))
	}
	if len(f.NamedReturnCalls) != 1 {
		panic(fmt.Sprintf("NamedReturnCalls: %d", len(f.NamedReturnCalls)))
	}
	if !namedHookCalled {
		panic("NamedReturnHook not called")
	}
	if !f.NamedReturnCalled() {
		panic("NamedReturnCalled: NamedReturn not called")
	}
	if !f.NamedReturnCalledOnce() {
		panic("NamedReturnCalledOnce: NamedReturn not called once")
	}
	if f.NamedReturnNotCalled() {
		panic("NamedReturnNotCalled: NamedReturn not called")
	}
	if !f.NamedReturnCalledN(1) {
		panic("NamedReturnCalledN: NamedReturn not called once")
	}
}
