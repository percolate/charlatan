package main

import (
	"fmt"
)

var _ Interfacer = &FakeInterfacer{}

func main() {
	interfaceHookCalled := false
	namedInterfaceHookCalled := false

	a := "one"
	b := "two"

	f := &FakeInterfacer{
		InterfaceHook: func(interface{}) interface{} {
			interfaceHookCalled = true
			return a
		},
		NamedInterfaceHook: func(interface{}) interface{} {
			namedInterfaceHookCalled = true
			return b
		},
	}

	ra := f.Interface(b)

	if ra != a {
		panic(fmt.Sprintf("Unexpected results from Multireturn: %s (expected %s)", ra, a))
	}
	if len(f.InterfaceCalls) != 1 {
		panic(fmt.Sprintf("InterfaceCalls: %d", len(f.InterfaceCalls)))
	}
	if !interfaceHookCalled {
		panic("InterfaceHook not called")
	}
	if !f.InterfaceCalled() {
		panic("InterfaceCalled: Interface not called")
	}
	if !f.InterfaceCalledOnce() {
		panic("InterfaceCalledOnce: Interface not called once")
	}
	if f.InterfaceNotCalled() {
		panic("InterfaceNotCalled: Interface not called")
	}
	if !f.InterfaceCalledN(1) {
		panic("InterfaceCalledN: Interface not called once")
	}

	res, found := f.InterfaceResultsForCall(b)
	if res != a || found != true {
		panic(fmt.Sprintf("NamedQualifyResultsForCall: NamedQualify results for %s not %s, found: %s", b, a, res, found))
	}

	rb := f.NamedInterface(a)

	if rb != b {
		panic(fmt.Sprintf("Unexpected results from Multireturn: %s (expected %s)", rb, b))
	}
	if len(f.NamedInterfaceCalls) != 1 {
		panic(fmt.Sprintf("NamedInterfaceCalls: %d", len(f.NamedInterfaceCalls)))
	}
	if !namedInterfaceHookCalled {
		panic("NamedInterfaceHook not called")
	}
	if !f.NamedInterfaceCalled() {
		panic("NamedInterfaceCalled: NamedInterface not called")
	}
	if !f.NamedInterfaceCalledOnce() {
		panic("NamedInterfaceCalledOnce: NamedInterface not called once")
	}
	if f.NamedInterfaceNotCalled() {
		panic("NamedInterfaceNotCalled: NamedInterface not called")
	}
	if !f.NamedInterfaceCalledN(1) {
		panic("NamedInterfaceCalledN: NamedInterface not called once")
	}

	res, found = f.NamedInterfaceResultsForCall(a)
	if res != b || found != true {
		panic(fmt.Sprintf("NamedQualifyResultsForCall: NamedQualify results for %s not %s, found: %s", a, b, res, found))
	}
}
