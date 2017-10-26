package main

import "fmt"

type Qualifier interface {
	Qualify(fmt.Scanner) *fmt.Scanner
	NamedQualify(a, b, c fmt.Scanner) (d *fmt.Scanner)
}
