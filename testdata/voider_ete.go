package main

import "fmt"

var _ Voider = &FakeVoider{}

func main() {
	hookCalled := false
	f := &FakeVoider{
		VoidMethodHook: func() {
			hookCalled = true
			return
		},
	}

	f.VoidMethod()

	if len(f.VoidMethodCalls) != 1 {
		panic(fmt.Sprintf("VoidMethodCalls: %d", len(f.VoidMethodCalls)))
	}
	if !hookCalled {
		panic("VoidMethodHook not called")
	}
	if !f.VoidMethodCalled() {
		panic("VoidMethodCalled: VoidMethod not called")
	}
	if !f.VoidMethodCalledOnce() {
		panic("VoidMethodCalledOnce: VoidMethod not called once")
	}
	if f.VoidMethodNotCalled() {
		panic("VoidMethodNotCalled: VoidMethod not called")
	}
	if !f.VoidMethodCalledN(1) {
		panic("VoidMethodCalledN: VoidMethod not called once")
	}
}
