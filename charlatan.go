package main // import "github.com/percolate/charlatan"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type stringSliceValue []string

func (s *stringSliceValue) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s *stringSliceValue) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSliceValue) Get() interface{} {
	return []string(*s)
}

var (
	interfaces  stringSliceValue
	outputName  = flag.String("output", "", "output file name [default: charlatan.go]")
	packageName = flag.String("package", "", "output package name [default: \"<current package>\"]")
)

func init() {
	flag.Var(&interfaces, "interface", "name of interface to fake, may be repeated")
	log.SetFlags(0)
	log.SetPrefix("charlatan: ")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `charlatan
https://github.com/percolate/charlatan

Usage:
  charlatan [options] (--interface <I>)...
  charlatan -h | --help

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if len(interfaces) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var (
		dir string
		g   Generator
	)

	g.interfaceNames = interfaces
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
		g.parsePackageDir(args[0])
	} else {
		dir = filepath.Dir(args[0])
		g.parsePackageFiles(args)
	}

	g.generate()

	// format the output.
	src := g.format()

	// Write to file.
	outputName := *outputName
	if outputName == "" {
		outputName = filepath.Join(dir, "charlatan.go")
	}

	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// isDirectory returns true if the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
