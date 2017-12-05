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
	outputPath    = flag.String("output", "", "output file path [default: ./charlatan.go]")
	outputPackage = flag.String("package", "", "output package name [default: \"<current package>\"]")
	dirName       = flag.String("dir", "", "input package directory [default: current package directory]")
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("charlatan: ")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, usageFormat)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Print("interface parameters are required")
		flag.Usage()
		os.Exit(1)
	}

	if *outputPath != "" && !strings.HasSuffix(*outputPath, ".go") {
		log.Print("output path must be a Go source file name")
		flag.Usage()
		os.Exit(1)
	}

	packageDirectory := "."
	if *dirName != "" {
		packageDirectory = *dirName
	}

	g, err := LoadPackageDir(packageDirectory)
	if err != nil {
		log.Fatal(err)
	}

	g.PackageOverride = *outputPackage

	src, err := g.Generate(flag.Args())
	if err != nil {
		log.Print(err)
	}
	if src == nil {
		os.Exit(1)
	}

	if *outputPath == "" {
		*outputPath = "charlatan.go"
	}

	if err := os.MkdirAll(filepath.Dir(*outputPath), 0755); err != nil {
		log.Fatalf("error writing output: %s", err)
	}

	if err := ioutil.WriteFile(*outputPath, src, 0644); err != nil {
		log.Fatalf("error writing output: %s", err)
	}

	log.Printf("wrote %s\n", *outputPath)
}
