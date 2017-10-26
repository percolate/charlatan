package main

type Structer interface {
	Struct(struct{}) struct{}
	NamedStruct(a struct{}) (z struct{})
}
