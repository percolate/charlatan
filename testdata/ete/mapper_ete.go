package main

import "fmt"

var _ Mapper = &FakeMapper{}

func main() {
	mapParameterHookCalled := false
	mapReturnHookCalled := false

	m := map[string]string{"answer": "42"}

	f := &FakeMapper{
		MapParameterHook: func(map[string]string) {
			mapParameterHookCalled = true
		},
		MapReturnHook: func() map[string]string {
			mapReturnHookCalled = true
			return m
		},
	}

	f.MapParameter(m)

	if len(f.MapParameterCalls) != 1 {
		panic(fmt.Sprintf("MapParameterCalls: %d", len(f.MapParameterCalls)))
	}
	if !mapParameterHookCalled {
		panic("MapParameterHook not called")
	}
	if !f.MapParameterCalled() {
		panic("MapParameterCalled: MapParameter not called")
	}
	if !f.MapParameterCalledOnce() {
		panic("MapParameterCalledOnce: MapParameter not called once")
	}
	if f.MapParameterNotCalled() {
		panic("MapParameterNotCalled: MapParameter not called")
	}
	if !f.MapParameterCalledN(1) {
		panic("MapParameterCalledN: MapParameter not called once")
	}
	if !f.MapParameterCalledWith(m) {
		panic(fmt.Sprintf("MapParameterCalledWith: MapParameter not called with %s", m))
	}
	if !f.MapParameterCalledOnceWith(m) {
		panic(fmt.Sprintf("MapParameterCalledOnceWith: MapParameter not called once with %s", m))
	}

	mr := f.MapReturn()

	if len(mr) != 1 {
		panic("MapReturn result has unexpected size")
	}
	if len(f.MapReturnCalls) != 1 {
		panic(fmt.Sprintf("MapReturnCalls: %d", len(f.MapReturnCalls)))
	}
	if !mapReturnHookCalled {
		panic("MapReturnHook not called")
	}
	if !f.MapReturnCalled() {
		panic("MapReturnCalled: MapReturn not called")
	}
	if !f.MapReturnCalledOnce() {
		panic("MapReturnCalledOnce: MapReturn not called once")
	}
	if f.MapReturnNotCalled() {
		panic("MapReturnNotCalled: MapReturn not called")
	}
	if !f.MapReturnCalledN(1) {
		panic("MapReturnCalledN: MapReturn not called once")
	}
}
