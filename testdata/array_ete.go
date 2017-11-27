package main

import "fmt"

var _ Array = &FakeArray{}

func main() {
	parameterHookCalled := false
	returnHookCalled := false
	s := []string{"one", "two", "three"}
	f := &FakeArray{
		ArrayParameterHook: func([]string) {
			parameterHookCalled = true
		},
		ArrayReturnHook: func() (a []string) {
			returnHookCalled = true
			return
		},
	}

	f.ArrayParameter(s)

	if len(f.ArrayParameterCalls) != 1 {
		panic(fmt.Sprintf("ArrayParameterCalls: %d", len(f.ArrayParameterCalls)))
	}
	if !parameterHookCalled {
		panic("ArrayParameterHook not called")
	}
	if !f.ArrayParameterCalled() {
		panic("ArrayParameterCalled: ArrayParameter not called")
	}
	if !f.ArrayParameterCalledOnce() {
		panic("ArrayParameterCalledOnce: ArrayParameter not called once")
	}
	if f.ArrayParameterNotCalled() {
		panic("ArrayParameterNotCalled: ArrayParameter not called")
	}
	if !f.ArrayParameterCalledN(1) {
		panic("ArrayParameterCalledN: ArrayParameter not called once")
	}
	if !f.ArrayParameterCalledWith(s) {
		panic(fmt.Sprintf("ArrayParameterCalledWith: ArrayParameter not called with %s", s))
	}
	if !f.ArrayParameterCalledOnceWith(s) {
		panic(fmt.Sprintf("ArrayParameterCalledOnceWith: ArrayParameter not called once with %s", s))
	}

	f.ArrayReturn()

	if len(f.ArrayReturnCalls) != 1 {
		panic(fmt.Sprintf("ArrayReturnCalls: %d", len(f.ArrayReturnCalls)))
	}
	if !returnHookCalled {
		panic("ArrayReturnHook not called")
	}
	if !f.ArrayReturnCalled() {
		panic("ArrayReturnCalled: ArrayReturn not called")
	}
	if !f.ArrayReturnCalledOnce() {
		panic("ArrayReturnCalledOnce: ArrayReturn not called once")
	}
	if f.ArrayReturnNotCalled() {
		panic("ArrayReturnNotCalled: ArrayReturn not called")
	}
	if !f.ArrayReturnCalledN(1) {
		panic("ArrayReturnCalledN: ArrayReturn not called once")
	}
}
