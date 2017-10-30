package main

import "fmt"

var _ Pointer = &FakePointer{}

func main() {
	pointHookCalled := false

	s := "hello"
	ps := &s

	h := 123
	ph := &h

	f := &FakePointer{
		PointHook: func(a *string) (int) {
			pointHookCalled = true
			x := len(*a)
			return x
		},
	}

	*ph = f.Point(ps)

	if len(f.PointCalls) != 1 {
		panic(fmt.Sprintf("PointCalls: %d", len(f.PointCalls)))
	}
	if !pointHookCalled {
		panic("PointHook not called")
	}
	if !f.PointCalled() {
		panic("PointCalled: Point not called")
	}
	if !f.PointCalledOnce() {
		panic("PointCalledOnce: Point not called once")
	}
	if f.PointNotCalled() {
		panic("PointNotCalled: Point not called")
	}
	if !f.PointCalledN(1) {
		panic("PointCalledN: Point not called once")
	}
	if !f.PointCalledWith(ps) {
		panic(fmt.Sprintf("PointCalledWith: Point not called with %s", ps))
	}
	if !f.PointCalledOnceWith(ps) {
		panic(fmt.Sprintf("PointCalledOnceWith: Point not called once with %s", ps))
	}
}
