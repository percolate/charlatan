package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

var (
	golden = []string{
		"Array",
		"Channeler",
		"Embedder",
		"Funcer",
		"Identifier",
		"Interfacer",
		"Importer",
		"Mapper",
		"Multireturner",
		"Namedvaluer",
		"Pointer",
		"Qualifier",
		"Structer",
		"Variadic",
		"Voider",
	}
	unsupported = []string{
		"_",
		"Emptier",
	}
	dmp = diffmatchpatch.New()
)

func TestUnsupported(t *testing.T) {
	for _, interfaceName := range unsupported {
		t.Run(interfaceName, CheckOneUnsupported)
	}
}

func CheckOneUnsupported(t *testing.T) {
	name := path.Base(t.Name())
	lname := strings.ToLower(name)
	inputFilename := fmt.Sprintf("testdata/%s/%s_def.go", lname, lname)

	g, err := parsePackage("testdata/"+lname, []string{inputFilename})
	if err != nil {
		t.Fatalf("parsePackage error: %s", err)
	}
	got, err := g.Generate([]string{name})

	if err == nil {
		t.Fatalf("expected error for %s", name)
	}
	if got != nil {
		t.Fatalf("Unexpected output result for %s", name)
	}
}

func TestGolden(t *testing.T) {
	for _, interfaceName := range golden {
		t.Run(interfaceName, CheckOneGolden)
	}
}

func CheckOneGolden(t *testing.T) {
	name := path.Base(t.Name())
	lname := strings.ToLower(name)

	symGen.reset()
	identSymGen.reset()

	inputFilename := fmt.Sprintf("./testdata/%s/%s_def.go", lname, lname)
	outputFilename := fmt.Sprintf("./testdata/%s/%s.go", lname, lname)

	outputFile, err := ioutil.ReadFile(outputFilename)
	if err != nil {
		outputFile = []byte{}
	}

	g, err := parsePackage("testdata/"+lname, []string{inputFilename})
	if err != nil {
		t.Fatalf("parsePackage error: %s", err)
	}
	got, err := g.Generate([]string{name})
	if err != nil {
		t.Fatalf("Generator.Generate error for %s: %s", name, err)
	}

	if len(got) == 0 {
		t.Fatalf("%q resulted in an empty file when the contents of %q were expected", name, outputFilename)
	}

	readableOutput := string(outputFile)
	readableResult := string(got)

	// Only compare everything after the first line to avoid
	// comparing the generation commands
	outputStart := strings.Index(readableOutput, "\n")
	resultStart := strings.Index(readableResult, "\n")

	diffs := dmp.DiffMain(readableOutput, readableResult, false)

	assert.Equal(t, readableOutput[outputStart:], readableResult[resultStart:], dmp.DiffPrettyText(diffs))
}
