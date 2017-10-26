package main

type Variadic interface {
	SingleVariadic(a... string)
	MixedVariadic(a, b, c int, d... string)
}
