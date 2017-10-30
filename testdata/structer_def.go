package main

type Structer interface {
	Struct(struct {
		a string
		b string
	}) struct {
		c string
		d string
	}
	NamedStruct(a struct {
		a string
		b string
	}) (z struct {
		c string
		d string
	})
}
