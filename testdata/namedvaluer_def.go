package main

type Namedvaluer interface {
	ManyNamed(a, b string, f, g int) (ret bool)
	Named(a int, b string) (ret bool)
}
