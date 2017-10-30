package main

import (
	"fmt"
)

var _ Embedder = &FakeEmbedder{}

func (e *FakeEmbedder) Embed() {
	// NOOP
}

func main() {
	interfaceHookCalled := false

	f := &FakeEmbedder{
		OtherHook: func() {
			interfaceHookCalled = true
		},
	}

	f.Other()

	if len(f.OtherCalls) != 1 {
		panic(fmt.Sprintf("OtherCalls: %d", len(f.OtherCalls)))
	}
	if !interfaceHookCalled {
		panic("OtherHook not called")
	}
	if !f.OtherCalled() {
		panic("OtherCalled: Other not called")
	}
	if !f.OtherCalledOnce() {
		panic("OtherCalledOnce: Other not called once")
	}
	if f.OtherNotCalled() {
		panic("OtherNotCalled: Other not called")
	}
	if !f.OtherCalledN(1) {
		panic("OtherCalledN: Other not called once")
	}
}
