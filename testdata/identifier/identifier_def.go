package main

type Identifier interface {
	// Issue #20 named return identifier conflicts with test context constructor parameter
	TestConstructor(val int64) (t string)
	// Issue #21 named return identifer conflicts with Set{Method}Invocation parameter(s)
	InvocationSetter(val int64) (call string, calls string, fallback string)
}
