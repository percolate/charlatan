package main

type Array interface {
	ArrayParameter([3]string)
	ArrayReturn() [3]string
	SliceParameter([]string)
	SliceReturn() []string
}
