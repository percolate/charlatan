// generated by "charlatan -dir=testdata/embedder -output=testdata/embedder/embedder.go Embedder".  DO NOT EDIT.

package main

import "reflect"

// EmbedderStringInvocation represents a single call of FakeEmbedder.String
type EmbedderStringInvocation struct {
	Results struct {
		Ident5 string
	}
}

// EmbedderEmbedInvocation represents a single call of FakeEmbedder.Embed
type EmbedderEmbedInvocation struct {
	Parameters struct {
		Ident1 string
	}
	Results struct {
		Ident2 string
	}
}

// EmbedderOtherInvocation represents a single call of FakeEmbedder.Other
type EmbedderOtherInvocation struct {
	Parameters struct {
		Ident1 string
	}
	Results struct {
		Ident2 string
	}
}

// EmbedderTestingT represents the methods of "testing".T used by charlatan Fakes.  It avoids importing the testing package.
type EmbedderTestingT interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Helper()
}

/*
FakeEmbedder is a mock implementation of Embedder for testing.
Use it in your tests as in this example:

	package example

	func TestWithEmbedder(t *testing.T) {
		f := &main.FakeEmbedder{
			StringHook: func() (ident5 string) {
				// ensure parameters meet expections, signal errors using t, etc
				return
			},
		}

		// test code goes here ...

		// assert state of FakeString ...
		f.AssertStringCalledOnce(t)
	}

Create anonymous function implementations for only those interface methods that
should be called in the code under test.  This will force a panic if any
unexpected calls are made to FakeString.
*/
type FakeEmbedder struct {
	StringHook func() string
	EmbedHook  func(string) string
	OtherHook  func(string) string

	StringCalls []*EmbedderStringInvocation
	EmbedCalls  []*EmbedderEmbedInvocation
	OtherCalls  []*EmbedderOtherInvocation
}

// NewFakeEmbedderDefaultPanic returns an instance of FakeEmbedder with all hooks configured to panic
func NewFakeEmbedderDefaultPanic() *FakeEmbedder {
	return &FakeEmbedder{
		StringHook: func() (ident5 string) {
			panic("Unexpected call to Embedder.String")
		},
		EmbedHook: func(string) (ident2 string) {
			panic("Unexpected call to Embedder.Embed")
		},
		OtherHook: func(string) (ident2 string) {
			panic("Unexpected call to Embedder.Other")
		},
	}
}

// NewFakeEmbedderDefaultFatal returns an instance of FakeEmbedder with all hooks configured to call t.Fatal
func NewFakeEmbedderDefaultFatal(t EmbedderTestingT) *FakeEmbedder {
	return &FakeEmbedder{
		StringHook: func() (ident5 string) {
			t.Fatal("Unexpected call to Embedder.String")
			return
		},
		EmbedHook: func(string) (ident2 string) {
			t.Fatal("Unexpected call to Embedder.Embed")
			return
		},
		OtherHook: func(string) (ident2 string) {
			t.Fatal("Unexpected call to Embedder.Other")
			return
		},
	}
}

// NewFakeEmbedderDefaultError returns an instance of FakeEmbedder with all hooks configured to call t.Error
func NewFakeEmbedderDefaultError(t EmbedderTestingT) *FakeEmbedder {
	return &FakeEmbedder{
		StringHook: func() (ident5 string) {
			t.Error("Unexpected call to Embedder.String")
			return
		},
		EmbedHook: func(string) (ident2 string) {
			t.Error("Unexpected call to Embedder.Embed")
			return
		},
		OtherHook: func(string) (ident2 string) {
			t.Error("Unexpected call to Embedder.Other")
			return
		},
	}
}

func (f *FakeEmbedder) Reset() {
	f.StringCalls = []*EmbedderStringInvocation{}
	f.EmbedCalls = []*EmbedderEmbedInvocation{}
	f.OtherCalls = []*EmbedderOtherInvocation{}
}

func (_f1 *FakeEmbedder) String() (ident5 string) {
	if _f1.StringHook == nil {
		panic("Embedder.String() called but FakeEmbedder.StringHook is nil")
	}

	invocation := new(EmbedderStringInvocation)
	_f1.StringCalls = append(_f1.StringCalls, invocation)

	ident5 = _f1.StringHook()

	invocation.Results.Ident5 = ident5

	return
}

// StringCalled returns true if FakeEmbedder.String was called
func (f *FakeEmbedder) StringCalled() bool {
	return len(f.StringCalls) != 0
}

// AssertStringCalled calls t.Error if FakeEmbedder.String was not called
func (f *FakeEmbedder) AssertStringCalled(t EmbedderTestingT) {
	t.Helper()
	if len(f.StringCalls) == 0 {
		t.Error("FakeEmbedder.String not called, expected at least one")
	}
}

// StringNotCalled returns true if FakeEmbedder.String was not called
func (f *FakeEmbedder) StringNotCalled() bool {
	return len(f.StringCalls) == 0
}

// AssertStringNotCalled calls t.Error if FakeEmbedder.String was called
func (f *FakeEmbedder) AssertStringNotCalled(t EmbedderTestingT) {
	t.Helper()
	if len(f.StringCalls) != 0 {
		t.Error("FakeEmbedder.String called, expected none")
	}
}

// StringCalledOnce returns true if FakeEmbedder.String was called exactly once
func (f *FakeEmbedder) StringCalledOnce() bool {
	return len(f.StringCalls) == 1
}

// AssertStringCalledOnce calls t.Error if FakeEmbedder.String was not called exactly once
func (f *FakeEmbedder) AssertStringCalledOnce(t EmbedderTestingT) {
	t.Helper()
	if len(f.StringCalls) != 1 {
		t.Errorf("FakeEmbedder.String called %d times, expected 1", len(f.StringCalls))
	}
}

// StringCalledN returns true if FakeEmbedder.String was called at least n times
func (f *FakeEmbedder) StringCalledN(n int) bool {
	return len(f.StringCalls) >= n
}

// AssertStringCalledN calls t.Error if FakeEmbedder.String was called less than n times
func (f *FakeEmbedder) AssertStringCalledN(t EmbedderTestingT, n int) {
	t.Helper()
	if len(f.StringCalls) < n {
		t.Errorf("FakeEmbedder.String called %d times, expected >= %d", len(f.StringCalls), n)
	}
}

func (_f2 *FakeEmbedder) Embed(ident1 string) (ident2 string) {
	if _f2.EmbedHook == nil {
		panic("Embedder.Embed() called but FakeEmbedder.EmbedHook is nil")
	}

	invocation := new(EmbedderEmbedInvocation)
	_f2.EmbedCalls = append(_f2.EmbedCalls, invocation)

	invocation.Parameters.Ident1 = ident1

	ident2 = _f2.EmbedHook(ident1)

	invocation.Results.Ident2 = ident2

	return
}

// EmbedCalled returns true if FakeEmbedder.Embed was called
func (f *FakeEmbedder) EmbedCalled() bool {
	return len(f.EmbedCalls) != 0
}

// AssertEmbedCalled calls t.Error if FakeEmbedder.Embed was not called
func (f *FakeEmbedder) AssertEmbedCalled(t EmbedderTestingT) {
	t.Helper()
	if len(f.EmbedCalls) == 0 {
		t.Error("FakeEmbedder.Embed not called, expected at least one")
	}
}

// EmbedNotCalled returns true if FakeEmbedder.Embed was not called
func (f *FakeEmbedder) EmbedNotCalled() bool {
	return len(f.EmbedCalls) == 0
}

// AssertEmbedNotCalled calls t.Error if FakeEmbedder.Embed was called
func (f *FakeEmbedder) AssertEmbedNotCalled(t EmbedderTestingT) {
	t.Helper()
	if len(f.EmbedCalls) != 0 {
		t.Error("FakeEmbedder.Embed called, expected none")
	}
}

// EmbedCalledOnce returns true if FakeEmbedder.Embed was called exactly once
func (f *FakeEmbedder) EmbedCalledOnce() bool {
	return len(f.EmbedCalls) == 1
}

// AssertEmbedCalledOnce calls t.Error if FakeEmbedder.Embed was not called exactly once
func (f *FakeEmbedder) AssertEmbedCalledOnce(t EmbedderTestingT) {
	t.Helper()
	if len(f.EmbedCalls) != 1 {
		t.Errorf("FakeEmbedder.Embed called %d times, expected 1", len(f.EmbedCalls))
	}
}

// EmbedCalledN returns true if FakeEmbedder.Embed was called at least n times
func (f *FakeEmbedder) EmbedCalledN(n int) bool {
	return len(f.EmbedCalls) >= n
}

// AssertEmbedCalledN calls t.Error if FakeEmbedder.Embed was called less than n times
func (f *FakeEmbedder) AssertEmbedCalledN(t EmbedderTestingT, n int) {
	t.Helper()
	if len(f.EmbedCalls) < n {
		t.Errorf("FakeEmbedder.Embed called %d times, expected >= %d", len(f.EmbedCalls), n)
	}
}

// EmbedCalledWith returns true if FakeEmbedder.Embed was called with the given values
func (_f3 *FakeEmbedder) EmbedCalledWith(ident1 string) (found bool) {
	for _, call := range _f3.EmbedCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			found = true
			break
		}
	}

	return
}

// AssertEmbedCalledWith calls t.Error if FakeEmbedder.Embed was not called with the given values
func (_f4 *FakeEmbedder) AssertEmbedCalledWith(t EmbedderTestingT, ident1 string) {
	t.Helper()
	var found bool
	for _, call := range _f4.EmbedCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			found = true
			break
		}
	}

	if !found {
		t.Error("FakeEmbedder.Embed not called with expected parameters")
	}
}

// EmbedCalledOnceWith returns true if FakeEmbedder.Embed was called exactly once with the given values
func (_f5 *FakeEmbedder) EmbedCalledOnceWith(ident1 string) bool {
	var count int
	for _, call := range _f5.EmbedCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			count++
		}
	}

	return count == 1
}

// AssertEmbedCalledOnceWith calls t.Error if FakeEmbedder.Embed was not called exactly once with the given values
func (_f6 *FakeEmbedder) AssertEmbedCalledOnceWith(t EmbedderTestingT, ident1 string) {
	t.Helper()
	var count int
	for _, call := range _f6.EmbedCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			count++
		}
	}

	if count != 1 {
		t.Errorf("FakeEmbedder.Embed called %d times with expected parameters, expected one", count)
	}
}

// EmbedResultsForCall returns the result values for the first call to FakeEmbedder.Embed with the given values
func (_f7 *FakeEmbedder) EmbedResultsForCall(ident1 string) (ident2 string, found bool) {
	for _, call := range _f7.EmbedCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			ident2 = call.Results.Ident2
			found = true
			break
		}
	}

	return
}

func (_f8 *FakeEmbedder) Other(ident1 string) (ident2 string) {
	if _f8.OtherHook == nil {
		panic("Embedder.Other() called but FakeEmbedder.OtherHook is nil")
	}

	invocation := new(EmbedderOtherInvocation)
	_f8.OtherCalls = append(_f8.OtherCalls, invocation)

	invocation.Parameters.Ident1 = ident1

	ident2 = _f8.OtherHook(ident1)

	invocation.Results.Ident2 = ident2

	return
}

// OtherCalled returns true if FakeEmbedder.Other was called
func (f *FakeEmbedder) OtherCalled() bool {
	return len(f.OtherCalls) != 0
}

// AssertOtherCalled calls t.Error if FakeEmbedder.Other was not called
func (f *FakeEmbedder) AssertOtherCalled(t EmbedderTestingT) {
	t.Helper()
	if len(f.OtherCalls) == 0 {
		t.Error("FakeEmbedder.Other not called, expected at least one")
	}
}

// OtherNotCalled returns true if FakeEmbedder.Other was not called
func (f *FakeEmbedder) OtherNotCalled() bool {
	return len(f.OtherCalls) == 0
}

// AssertOtherNotCalled calls t.Error if FakeEmbedder.Other was called
func (f *FakeEmbedder) AssertOtherNotCalled(t EmbedderTestingT) {
	t.Helper()
	if len(f.OtherCalls) != 0 {
		t.Error("FakeEmbedder.Other called, expected none")
	}
}

// OtherCalledOnce returns true if FakeEmbedder.Other was called exactly once
func (f *FakeEmbedder) OtherCalledOnce() bool {
	return len(f.OtherCalls) == 1
}

// AssertOtherCalledOnce calls t.Error if FakeEmbedder.Other was not called exactly once
func (f *FakeEmbedder) AssertOtherCalledOnce(t EmbedderTestingT) {
	t.Helper()
	if len(f.OtherCalls) != 1 {
		t.Errorf("FakeEmbedder.Other called %d times, expected 1", len(f.OtherCalls))
	}
}

// OtherCalledN returns true if FakeEmbedder.Other was called at least n times
func (f *FakeEmbedder) OtherCalledN(n int) bool {
	return len(f.OtherCalls) >= n
}

// AssertOtherCalledN calls t.Error if FakeEmbedder.Other was called less than n times
func (f *FakeEmbedder) AssertOtherCalledN(t EmbedderTestingT, n int) {
	t.Helper()
	if len(f.OtherCalls) < n {
		t.Errorf("FakeEmbedder.Other called %d times, expected >= %d", len(f.OtherCalls), n)
	}
}

// OtherCalledWith returns true if FakeEmbedder.Other was called with the given values
func (_f9 *FakeEmbedder) OtherCalledWith(ident1 string) (found bool) {
	for _, call := range _f9.OtherCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			found = true
			break
		}
	}

	return
}

// AssertOtherCalledWith calls t.Error if FakeEmbedder.Other was not called with the given values
func (_f10 *FakeEmbedder) AssertOtherCalledWith(t EmbedderTestingT, ident1 string) {
	t.Helper()
	var found bool
	for _, call := range _f10.OtherCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			found = true
			break
		}
	}

	if !found {
		t.Error("FakeEmbedder.Other not called with expected parameters")
	}
}

// OtherCalledOnceWith returns true if FakeEmbedder.Other was called exactly once with the given values
func (_f11 *FakeEmbedder) OtherCalledOnceWith(ident1 string) bool {
	var count int
	for _, call := range _f11.OtherCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			count++
		}
	}

	return count == 1
}

// AssertOtherCalledOnceWith calls t.Error if FakeEmbedder.Other was not called exactly once with the given values
func (_f12 *FakeEmbedder) AssertOtherCalledOnceWith(t EmbedderTestingT, ident1 string) {
	t.Helper()
	var count int
	for _, call := range _f12.OtherCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			count++
		}
	}

	if count != 1 {
		t.Errorf("FakeEmbedder.Other called %d times with expected parameters, expected one", count)
	}
}

// OtherResultsForCall returns the result values for the first call to FakeEmbedder.Other with the given values
func (_f13 *FakeEmbedder) OtherResultsForCall(ident1 string) (ident2 string, found bool) {
	for _, call := range _f13.OtherCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			ident2 = call.Results.Ident2
			found = true
			break
		}
	}

	return
}
