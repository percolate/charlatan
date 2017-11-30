package main

import (
	"fmt"
)

var _ Embedder = &FakeEmbedder{}

func main() {
	embedHookCalled := false
	otherHookCalled := false
	stringHookCalled := false

	answer := "42"

	f := &FakeEmbedder{
		EmbedHook: func(v string) string {
			embedHookCalled = true
			return v
		},
		OtherHook: func(v string) string {
			otherHookCalled = true
			return v
		},
		StringHook: func() string {
			stringHookCalled = true
			return answer
		},
	}

	e := f.Embed(answer)

	if e != answer {
		panic("unexpected result from Embed method")
	}
	if len(f.EmbedCalls) != 1 {
		panic(fmt.Sprintf("EmbedCalls: %d", len(f.EmbedCalls)))
	}
	if !embedHookCalled {
		panic("EmbedHook not called")
	}
	if !f.EmbedCalled() {
		panic("EmbedCalled: Embed not called")
	}
	if !f.EmbedCalledOnce() {
		panic("EmbedCalledOnce: Embed not called once")
	}
	if f.EmbedNotCalled() {
		panic("EmbedNotCalled: Embed not called")
	}
	if !f.EmbedCalledN(1) {
		panic("EmbedCalledN: Embed not called once")
	}
	if !f.EmbedCalledWith(answer) {
		panic(fmt.Sprintf("EmbedCalledWith: Embed not called with %s", answer))
	}
	if !f.EmbedCalledOnceWith(answer) {
		panic(fmt.Sprintf("EmbedCalledOnceWith: Embed not called once with %s", answer))
	}

	o := f.Other(answer)

	if o != answer {
		panic("unexpected result from Other method")
	}
	if len(f.OtherCalls) != 1 {
		panic(fmt.Sprintf("OtherCalls: %d", len(f.OtherCalls)))
	}
	if !otherHookCalled {
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
	if !f.OtherCalledWith(answer) {
		panic(fmt.Sprintf("OtherCalledWith: Other not called with %s", answer))
	}
	if !f.OtherCalledOnceWith(answer) {
		panic(fmt.Sprintf("OtherCalledOnceWith: Other not called once with %s", answer))
	}

	s := f.String()

	if s != answer {
		panic("unexpected result from String method")
	}
	if len(f.StringCalls) != 1 {
		panic(fmt.Sprintf("StringCalls: %d", len(f.StringCalls)))
	}
	if !stringHookCalled {
		panic("StringHook not called")
	}
	if !f.StringCalled() {
		panic("StringCalled: String not called")
	}
	if !f.StringCalledOnce() {
		panic("StringCalledOnce: String not called once")
	}
	if f.StringNotCalled() {
		panic("StringNotCalled: String not called")
	}
	if !f.StringCalledN(1) {
		panic("StringCalledN: String not called once")
	}
}
