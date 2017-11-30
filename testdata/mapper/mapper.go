// generated by "charlatan -dir=testdata/mapper -output=testdata/mapper/mapper.go Mapper".  DO NOT EDIT.

package main

import (
	"reflect"
	"testing"
)

// MapParameterInvocation represents a single call of FakeMapper.MapParameter
type MapParameterInvocation struct {
	Parameters struct {
		Ident1 map[string]string
	}
}

// MapReturnInvocation represents a single call of FakeMapper.MapReturn
type MapReturnInvocation struct {
	Results struct {
		Ident2 map[string]string
	}
}

/*
FakeMapper is a mock implementation of Mapper for testing.
Use it in your tests as in this example:

	package example

	func TestWithMapper(t *testing.T) {
		f := &main.FakeMapper{
			MapParameterHook: func(ident1 map[string]string) () {
				// ensure parameters meet expections, signal errors using t, etc
				return
			},
		}

		// test code goes here ...

		// assert state of FakeMapParameter ...
		f.AssertMapParameterCalledOnce(t)
	}

Create anonymous function implementations for only those interface methods that
should be called in the code under test.  This will force a painc if any
unexpected calls are made to FakeMapParameter.
*/
type FakeMapper struct {
	MapParameterHook func(map[string]string)
	MapReturnHook    func() map[string]string

	MapParameterCalls []*MapParameterInvocation
	MapReturnCalls    []*MapReturnInvocation
}

// NewFakeMapperDefaultPanic returns an instance of FakeMapper with all hooks configured to panic
func NewFakeMapperDefaultPanic() *FakeMapper {
	return &FakeMapper{
		MapParameterHook: func(map[string]string) {
			panic("Unexpected call to Mapper.MapParameter")
			return
		},
		MapReturnHook: func() (ident2 map[string]string) {
			panic("Unexpected call to Mapper.MapReturn")
			return
		},
	}
}

// NewFakeMapperDefaultFatal returns an instance of FakeMapper with all hooks configured to call t.Fatal
func NewFakeMapperDefaultFatal(t *testing.T) *FakeMapper {
	return &FakeMapper{
		MapParameterHook: func(map[string]string) {
			t.Fatal("Unexpected call to Mapper.MapParameter")
			return
		},
		MapReturnHook: func() (ident2 map[string]string) {
			t.Fatal("Unexpected call to Mapper.MapReturn")
			return
		},
	}
}

// NewFakeMapperDefaultError returns an instance of FakeMapper with all hooks configured to call t.Error
func NewFakeMapperDefaultError(t *testing.T) *FakeMapper {
	return &FakeMapper{
		MapParameterHook: func(map[string]string) {
			t.Error("Unexpected call to Mapper.MapParameter")
			return
		},
		MapReturnHook: func() (ident2 map[string]string) {
			t.Error("Unexpected call to Mapper.MapReturn")
			return
		},
	}
}

func (_f1 *FakeMapper) MapParameter(ident1 map[string]string) {
	invocation := new(MapParameterInvocation)

	invocation.Parameters.Ident1 = ident1

	_f1.MapParameterHook(ident1)

	_f1.MapParameterCalls = append(_f1.MapParameterCalls, invocation)

	return
}

// MapParameterCalled returns true if FakeMapper.MapParameter was called
func (f *FakeMapper) MapParameterCalled() bool {
	return len(f.MapParameterCalls) != 0
}

// AssertMapParameterCalled calls t.Error if FakeMapper.MapParameter was not called
func (f *FakeMapper) AssertMapParameterCalled(t *testing.T) {
	t.Helper()
	if len(f.MapParameterCalls) == 0 {
		t.Error("FakeMapper.MapParameter not called, expected at least one")
	}
}

// MapParameterNotCalled returns true if FakeMapper.MapParameter was not called
func (f *FakeMapper) MapParameterNotCalled() bool {
	return len(f.MapParameterCalls) == 0
}

// AssertMapParameterNotCalled calls t.Error if FakeMapper.MapParameter was called
func (f *FakeMapper) AssertMapParameterNotCalled(t *testing.T) {
	t.Helper()
	if len(f.MapParameterCalls) != 0 {
		t.Error("FakeMapper.MapParameter called, expected none")
	}
}

// MapParameterCalledOnce returns true if FakeMapper.MapParameter was called exactly once
func (f *FakeMapper) MapParameterCalledOnce() bool {
	return len(f.MapParameterCalls) == 1
}

// AssertMapParameterCalledOnce calls t.Error if FakeMapper.MapParameter was not called exactly once
func (f *FakeMapper) AssertMapParameterCalledOnce(t *testing.T) {
	t.Helper()
	if len(f.MapParameterCalls) != 1 {
		t.Errorf("FakeMapper.MapParameter called %d times, expected 1", len(f.MapParameterCalls))
	}
}

// MapParameterCalledN returns true if FakeMapper.MapParameter was called at least n times
func (f *FakeMapper) MapParameterCalledN(n int) bool {
	return len(f.MapParameterCalls) >= n
}

// AssertMapParameterCalledN calls t.Error if FakeMapper.MapParameter was called less than n times
func (f *FakeMapper) AssertMapParameterCalledN(t *testing.T, n int) {
	t.Helper()
	if len(f.MapParameterCalls) < n {
		t.Errorf("FakeMapper.MapParameter called %d times, expected >= %d", len(f.MapParameterCalls), n)
	}
}

// MapParameterCalledWith returns true if FakeMapper.MapParameter was called with the given values
func (_f2 *FakeMapper) MapParameterCalledWith(ident1 map[string]string) (found bool) {
	for _, call := range _f2.MapParameterCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			found = true
			break
		}
	}

	return
}

// AssertMapParameterCalledWith calls t.Error if FakeMapper.MapParameter was not called with the given values
func (_f3 *FakeMapper) AssertMapParameterCalledWith(t *testing.T, ident1 map[string]string) {
	t.Helper()
	var found bool
	for _, call := range _f3.MapParameterCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			found = true
			break
		}
	}

	if !found {
		t.Error("FakeMapper.MapParameter not called with expected parameters")
	}
}

// MapParameterCalledOnceWith returns true if FakeMapper.MapParameter was called exactly once with the given values
func (_f4 *FakeMapper) MapParameterCalledOnceWith(ident1 map[string]string) bool {
	var count int
	for _, call := range _f4.MapParameterCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			count++
		}
	}

	return count == 1
}

// AssertMapParameterCalledOnceWith calls t.Error if FakeMapper.MapParameter was not called exactly once with the given values
func (_f5 *FakeMapper) AssertMapParameterCalledOnceWith(t *testing.T, ident1 map[string]string) {
	t.Helper()
	var count int
	for _, call := range _f5.MapParameterCalls {
		if reflect.DeepEqual(call.Parameters.Ident1, ident1) {
			count++
		}
	}

	if count != 1 {
		t.Errorf("FakeMapper.MapParameter called %d times with expected parameters, expected one", count)
	}
}

func (_f6 *FakeMapper) MapReturn() (ident2 map[string]string) {
	invocation := new(MapReturnInvocation)

	ident2 = _f6.MapReturnHook()

	invocation.Results.Ident2 = ident2

	_f6.MapReturnCalls = append(_f6.MapReturnCalls, invocation)

	return
}

// MapReturnCalled returns true if FakeMapper.MapReturn was called
func (f *FakeMapper) MapReturnCalled() bool {
	return len(f.MapReturnCalls) != 0
}

// AssertMapReturnCalled calls t.Error if FakeMapper.MapReturn was not called
func (f *FakeMapper) AssertMapReturnCalled(t *testing.T) {
	t.Helper()
	if len(f.MapReturnCalls) == 0 {
		t.Error("FakeMapper.MapReturn not called, expected at least one")
	}
}

// MapReturnNotCalled returns true if FakeMapper.MapReturn was not called
func (f *FakeMapper) MapReturnNotCalled() bool {
	return len(f.MapReturnCalls) == 0
}

// AssertMapReturnNotCalled calls t.Error if FakeMapper.MapReturn was called
func (f *FakeMapper) AssertMapReturnNotCalled(t *testing.T) {
	t.Helper()
	if len(f.MapReturnCalls) != 0 {
		t.Error("FakeMapper.MapReturn called, expected none")
	}
}

// MapReturnCalledOnce returns true if FakeMapper.MapReturn was called exactly once
func (f *FakeMapper) MapReturnCalledOnce() bool {
	return len(f.MapReturnCalls) == 1
}

// AssertMapReturnCalledOnce calls t.Error if FakeMapper.MapReturn was not called exactly once
func (f *FakeMapper) AssertMapReturnCalledOnce(t *testing.T) {
	t.Helper()
	if len(f.MapReturnCalls) != 1 {
		t.Errorf("FakeMapper.MapReturn called %d times, expected 1", len(f.MapReturnCalls))
	}
}

// MapReturnCalledN returns true if FakeMapper.MapReturn was called at least n times
func (f *FakeMapper) MapReturnCalledN(n int) bool {
	return len(f.MapReturnCalls) >= n
}

// AssertMapReturnCalledN calls t.Error if FakeMapper.MapReturn was called less than n times
func (f *FakeMapper) AssertMapReturnCalledN(t *testing.T, n int) {
	t.Helper()
	if len(f.MapReturnCalls) < n {
		t.Errorf("FakeMapper.MapReturn called %d times, expected >= %d", len(f.MapReturnCalls), n)
	}
}
