package main

import (
	. "fmt"
	z "strings"
	_ "testing"
)

type Importer interface {
	Scan(*Scanner) z.Reader
}
