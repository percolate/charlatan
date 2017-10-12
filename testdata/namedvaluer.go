package main

import "fmt"

var _ Namedvaluer = &FakeNamedvaluer{}

func main() {
	hookVal := false
	f := &FakeNamedvaluer{
		ManyNamedHook: func(a, b, c, d string, f, g, h int) bool {
			hookVal = true
			return hookVal
		},
	}

	b := f.ManyNamed("a", "b", "c", "d", 1, 2, 3)

	if len(f.ManyNamedCalls) != 1 {
		panic(fmt.Sprintf("NamedHookCalls: %d", len(f.ManyNamedCalls)))
	}
	if !hookVal {
		panic("NamedHook not called")
	}

	if !b {
		panic("Named didn't return `true` as expected")
	}
}
