package main

type Pointer interface {
	Point(*string) *int
}
