package main

type Embeddable interface {
	Embed()
}

type Embedder interface {
	Embeddable
	Other()
}
