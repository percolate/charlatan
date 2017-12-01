package main

type Multireturner interface {
	MultiReturn() (string, int)
	NamedReturn() (a, b, c, d int)
}
