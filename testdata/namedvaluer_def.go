package main

type Namedvaluer interface {
	ManyNamed(a, b, c, d string, f, g, h int) (ret bool)
	Named(a int, b string) (ret bool)
}
