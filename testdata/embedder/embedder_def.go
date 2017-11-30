package main

import (
	"fmt"
)

type Embedder interface {
	fmt.Stringer
	Embeddable
	Other(string) string
}

type Embeddable interface {
	Embed(string) string
}
