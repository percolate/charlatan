package main

import "fmt"

var _ Qualifier = &FakeQualifier{}

type NiceScanner struct {
	name string
}

func (n *NiceScanner) Scan(state fmt.ScanState, verb rune) error {
	return nil
}

func main() {
	qualifiyHookCalled := false
	namedQualifyHookCalled := false

	scannerOne := &NiceScanner{"one"}
	scannerTwo := &NiceScanner{"two"}
	scannerThree := &NiceScanner{"three"}

	f := &FakeQualifier{
		QualifyHook: func(a fmt.Scanner) fmt.Scanner {
			qualifiyHookCalled = true
			return a
		},
		NamedQualifyHook: func(a, b, c fmt.Scanner) (z fmt.Scanner) {
			namedQualifyHookCalled = true
			return c
		},
	}

	f.Qualify(scannerOne)

	if len(f.QualifyCalls) != 1 {
		panic(fmt.Sprintf("QualifyCalls: %d", len(f.QualifyCalls)))
	}
	if !qualifiyHookCalled {
		panic("QualifyHook not called")
	}
	if !f.QualifyCalled() {
		panic("QualifyCalled: Qualify not called")
	}
	if !f.QualifyCalledOnce() {
		panic("QualifyCalledOnce: Qualify not called once")
	}
	if f.QualifyNotCalled() {
		panic("QualifyNotCalled: Qualify not called")
	}
	if !f.QualifyCalledN(1) {
		panic("QualifyCalledN: Qualify not called once")
	}
	if !f.QualifyCalledWith(scannerOne) {
		panic(fmt.Sprintf("QualifyCalledWith: Qualify not called with %s", scannerOne))
	}
	if !f.QualifyCalledOnceWith(scannerOne) {
		panic(fmt.Sprintf("QualifyCalledOnceWith: Qualify not called once with %s", scannerOne))
	}

	f.NamedQualify(scannerOne, scannerTwo, scannerThree)

	if len(f.NamedQualifyCalls) != 1 {
		panic(fmt.Sprintf("NamedQualifyCalls: %d", len(f.NamedQualifyCalls)))
	}
	if !namedQualifyHookCalled {
		panic("NamedQualifyHook not called")
	}
	if !f.NamedQualifyCalled() {
		panic("NamedQualifyCalled: NamedQualify not called")
	}
	if !f.NamedQualifyCalledOnce() {
		panic("NamedQualifyCalledOnce: NamedQualify not called once")
	}
	if f.NamedQualifyNotCalled() {
		panic("NamedQualifyNotCalled: NamedQualify not called")
	}
	if !f.NamedQualifyCalledN(1) {
		panic("NamedQualifyCalledN: NamedQualify not called once")
	}
	if !f.NamedQualifyCalledWith(scannerOne, scannerTwo, scannerThree) {
		panic(fmt.Sprintf("NamedQualifyCalledWith: NamedQualify not called once with %s, %s, %s", scannerOne, scannerTwo, scannerThree))
	}
	if !f.NamedQualifyCalledOnceWith(scannerOne, scannerTwo, scannerThree) {
		panic(fmt.Sprintf("NamedQualifyCalledOnceWith: NamedQualify not called once with %s, %s, %s", scannerOne, scannerTwo, scannerThree))
	}

	res, found := f.NamedQualifyResultsForCall(scannerOne, scannerTwo, scannerThree)
	if res != scannerThree || found != true {
		panic(fmt.Sprintf("NamedQualifyResultsForCall: NamedQualify results for %s, %s, %s not %s, found: %s", scannerOne, scannerTwo, scannerThree, scannerThree, found))
	}

	res, found = f.NamedQualifyResultsForCall(scannerTwo, scannerThree, scannerOne)
	if found != false {
		panic(fmt.Sprintf("NamedQualifyResultsForCall: NamedQualify results for %s found", scannerTwo, scannerThree, scannerOne))
	}

	f.NamedQualify(scannerTwo, scannerThree, scannerOne)

	if len(f.NamedQualifyCalls) != 2 {
		panic(fmt.Sprintf("NamedQualifyCalls: %d", len(f.NamedQualifyCalls)))
	}

	if !f.NamedQualifyCalledN(2) {
		panic("NamedQualifyCalledN: NamedQualify not called twice")
	}

	res, found = f.NamedQualifyResultsForCall(scannerTwo, scannerThree, scannerOne)
	if res != scannerOne || found != true {
		panic(fmt.Sprintf("NamedQualifyResultsForCall: NamedQualify results for %s, %s, %s not %s, found: %s", scannerTwo, scannerThree, scannerOne, scannerOne, found))
	}
}
