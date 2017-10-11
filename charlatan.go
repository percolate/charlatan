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

const (
	usageFormat = `charlatan
https://github.com/percolate/charlatan

Usage:
  charlatan [options] <interface> ...
  charlatan -h | --help

Options:
`
)

var (
	outputName  = flag.String("output", "", "output file name [default: charlatan.go]")
	packageName = flag.String("package", "", "output package name [default: \"<current package>\"]")
	dirName     = flag.String("dir", "", "input package directory [default: current package directory]")
	fileNames   stringSliceValue
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("charlatan: ")
	flag.Usage = usage
	flag.Var(&fileNames, "file", "name of input file, may be repeated, ignored if -dir is present")
}

func usage() {
	fmt.Fprintf(os.Stderr, usageFormat)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	var (
		dir string
		g   Generator
	)

	g.interfaces = flag.Args()

	if *dirName != "" {
		dir = *dirName
		g.parsePackageDir(dir)
	} else if len(fileNames) != 0 {
		dir = filepath.Dir(fileNames[0])
		for _, name := range fileNames[1:] {
			if dir != filepath.Dir(name) {
				log.Fatal("all input source files must be in the same package directory")
			}
		}
		g.parsePackageFiles(fileNames)
	} else {
		// process the package in current directory.
		dir = *dirName
		g.parsePackageDir(".")
	}

	src := g.generate()

	if *outputName == "" {
		*outputName = "charlatan.go"
	}

	output := filepath.Join(dir, *outputName)
	if err := ioutil.WriteFile(output, src, 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
