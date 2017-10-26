package main

import "fmt"

var _ Variadic = &FakeVariadic{}

func main() {
	singleHookCalled := false
	mixedHookCalled := false
	s := []string{"one", "two", "three"}
	f := &FakeVariadic{
		SingleVariadicHook: func(z... string) {
			singleHookCalled = true
		},
		MixedVariadicHook: func(int, int, int, ...string) {
			mixedHookCalled = true
		},
	}

	f.SingleVariadic(s...)

	if len(f.SingleVariadicCalls) != 1 {
		panic(fmt.Sprintf("SingleVariadicCalls: %d", len(f.SingleVariadicCalls)))
	}
	if !singleHookCalled {
		panic("SingleVariadicHook not called")
	}
	if !f.SingleVariadicCalled() {
		panic("SingleVariadicCalled: SingleVariadic not called")
	}
	if !f.SingleVariadicCalledOnce() {
		panic("SingleVariadicCalledOnce: SingleVariadic not called once")
	}
	if f.SingleVariadicNotCalled() {
		panic("SingleVariadicNotCalled: SingleVariadic not called")
	}
	if !f.SingleVariadicCalledN(1) {
		panic("SingleVariadicCalledN: SingleVariadic not called once")
	}
	if !f.SingleVariadicCalledWith(s...) {
		panic(fmt.Sprintf("SingleVariadicCalledWith: SingleVariadic not called with %s", s))
	}
	if !f.SingleVariadicCalledOnceWith(s...) {
		panic(fmt.Sprintf("SingleVariadicCalledOnceWith: SingleVariadic not called once with %s", s))
	}

	f.MixedVariadic(1, 2, 3, s...)

	if len(f.MixedVariadicCalls) != 1 {
		panic(fmt.Sprintf("MixedVariadicCalls: %d", len(f.MixedVariadicCalls)))
	}
	if !mixedHookCalled {
		panic("MixedVariadicHook not called")
	}
	if !f.MixedVariadicCalled() {
		panic("MixedVariadicCalled: MixedVariadic not called")
	}
	if !f.MixedVariadicCalledOnce() {
		panic("MixedVariadicCalledOnce: MixedVariadic not called once")
	}
	if f.MixedVariadicNotCalled() {
		panic("MixedVariadicNotCalled: MixedVariadic not called")
	}
	if !f.MixedVariadicCalledN(1) {
		panic("MixedVariadicCalledN: MixedVariadic not called once")
	}
	if !f.MixedVariadicCalledWith(1, 2, 3, s...) {
		panic(fmt.Sprintf("MixedVariadicCalledWith: MixedVariadic not called once with 1, 2, 3, %s", s))
	}
	if !f.MixedVariadicCalledOnceWith(1, 2, 3, s...) {
		panic(fmt.Sprintf("MixedVariadicCalledOnceWith: MixedVariadic not called once with 1, 2, 3, %s", s))
	}
}
