package main

type Embedder interface {
	Embeddable
	Other(string) string
}

type Embeddable interface {
	Embed(string) string
}
