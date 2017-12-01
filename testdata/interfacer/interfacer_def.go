package main

type Interfacer interface {
	Interface(interface{}) interface{}
	NamedInterface(a interface{}) (z interface{})
}
