package main

import "fmt"

var _ Array = &FakeArray{}

func main() {
	arrayParameterHookCalled := false
	arrayReturnHookCalled := false
	sliceParameterHookCalled := false
	sliceReturnHookCalled := false

	a := [...]string{"red", "green", "blue"}
	s := []string{"one", "two", "three"}

	f := &FakeArray{
		ArrayParameterHook: func([3]string) {
			arrayParameterHookCalled = true
		},
		ArrayReturnHook: func() [3]string {
			arrayReturnHookCalled = true
			return a
		},
		SliceParameterHook: func([]string) {
			sliceParameterHookCalled = true
		},
		SliceReturnHook: func() []string {
			sliceReturnHookCalled = true
			return s
		},
	}

	f.ArrayParameter(a)

	if len(f.ArrayParameterCalls) != 1 {
		panic(fmt.Sprintf("ArrayParameterCalls: %d", len(f.ArrayParameterCalls)))
	}
	if !arrayParameterHookCalled {
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
	if !f.ArrayParameterCalledWith(a) {
		panic(fmt.Sprintf("ArrayParameterCalledWith: ArrayParameter not called with %s", a))
	}
	if !f.ArrayParameterCalledOnceWith(a) {
		panic(fmt.Sprintf("ArrayParameterCalledOnceWith: ArrayParameter not called once with %s", a))
	}

	ar := f.ArrayReturn()

	if len(ar) != 3 {
		panic("ArrayReturn result has unexpected size")
	}
	if len(f.ArrayReturnCalls) != 1 {
		panic(fmt.Sprintf("ArrayReturnCalls: %d", len(f.ArrayReturnCalls)))
	}
	if !arrayReturnHookCalled {
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

	f.SliceParameter(s)

	if len(f.SliceParameterCalls) != 1 {
		panic(fmt.Sprintf("SliceParameterCalls: %d", len(f.SliceParameterCalls)))
	}
	if !sliceParameterHookCalled {
		panic("SliceParameterHook not called")
	}
	if !f.SliceParameterCalled() {
		panic("SliceParameterCalled: SliceParameter not called")
	}
	if !f.SliceParameterCalledOnce() {
		panic("SliceParameterCalledOnce: SliceParameter not called once")
	}
	if f.SliceParameterNotCalled() {
		panic("SliceParameterNotCalled: SliceParameter not called")
	}
	if !f.SliceParameterCalledN(1) {
		panic("SliceParameterCalledN: SliceParameter not called once")
	}
	if !f.SliceParameterCalledWith(s) {
		panic(fmt.Sprintf("SliceParameterCalledWith: SliceParameter not called with %s", s))
	}
	if !f.SliceParameterCalledOnceWith(s) {
		panic(fmt.Sprintf("SliceParameterCalledOnceWith: SliceParameter not called once with %s", s))
	}

	sr := f.SliceReturn()

	if sr == nil {
		panic("SliceReturn result was nil")
	}
	if len(f.SliceReturnCalls) != 1 {
		panic(fmt.Sprintf("SliceReturnCalls: %d", len(f.SliceReturnCalls)))
	}
	if !sliceReturnHookCalled {
		panic("SliceReturnHook not called")
	}
	if !f.SliceReturnCalled() {
		panic("SliceReturnCalled: SliceReturn not called")
	}
	if !f.SliceReturnCalledOnce() {
		panic("SliceReturnCalledOnce: SliceReturn not called once")
	}
	if f.SliceReturnNotCalled() {
		panic("SliceReturnNotCalled: SliceReturn not called")
	}
	if !f.SliceReturnCalledN(1) {
		panic("SliceReturnCalledN: SliceReturn not called once")
	}
}
