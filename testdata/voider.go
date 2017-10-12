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
}
