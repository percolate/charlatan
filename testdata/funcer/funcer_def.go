package main

type Funcer interface {
	FuncParameter(func(string) string)
	FuncReturn() func(string) string
}
