package main

import "fmt"

var _ Funcer = &FakeFuncer{}

func oracle(s string) string {
	if s == "answer" {
		return "42"
	}

	return ""
}

func main() {
	funcParameterHookCalled := false
	funcReturnHookCalled := false

	f := &FakeFuncer{
		FuncParameterHook: func(func(string) string) {
			funcParameterHookCalled = true
		},
		FuncReturnHook: func() func(string) string {
			funcReturnHookCalled = true
			return oracle
		},
	}

	f.FuncParameter(oracle)

	if len(f.FuncParameterCalls) != 1 {
		panic(fmt.Sprintf("FuncParameterCalls: %d", len(f.FuncParameterCalls)))
	}
	if !funcParameterHookCalled {
		panic("FuncParameterHook not called")
	}
	if !f.FuncParameterCalled() {
		panic("FuncParameterCalled: FuncParameter not called")
	}
	if !f.FuncParameterCalledOnce() {
		panic("FuncParameterCalledOnce: FuncParameter not called once")
	}
	if f.FuncParameterNotCalled() {
		panic("FuncParameterNotCalled: FuncParameter not called")
	}
	if !f.FuncParameterCalledN(1) {
		panic("FuncParameterCalledN: FuncParameter not called once")
	}
	// N.B. function references can't be compared

	fr := f.FuncReturn()

	if fr == nil {
		panic("FuncReturn result is nil")
	}
	r := fr("answer")
	if r != "42" {
		panic("unexpected result from function result")
	}
	if len(f.FuncReturnCalls) != 1 {
		panic(fmt.Sprintf("FuncReturnCalls: %d", len(f.FuncReturnCalls)))
	}
	if !funcReturnHookCalled {
		panic("FuncReturnHook not called")
	}
	if !f.FuncReturnCalled() {
		panic("FuncReturnCalled: FuncReturn not called")
	}
	if !f.FuncReturnCalledOnce() {
		panic("FuncReturnCalledOnce: FuncReturn not called once")
	}
	if f.FuncReturnNotCalled() {
		panic("FuncReturnNotCalled: FuncReturn not called")
	}
	if !f.FuncReturnCalledN(1) {
		panic("FuncReturnCalledN: FuncReturn not called once")
	}
}
