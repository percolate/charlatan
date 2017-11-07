package main

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

var golden = []string{
	"_",
	"Emptier",
	"Channeler",
	"Voider",
	"Namedvaluer",
	"Multireturner",
	"Variadic",
	"Pointer",
	"Interfacer",
	"Structer",
	"Embedder",
	"Qualifier",
}

func TestGolden(t *testing.T) {
	dmp := diffmatchpatch.New()

	for _, interfaceName := range golden {
		inputFilename := "./testdata/" + strings.ToLower(interfaceName) + "_def.go"
		outputFilename := "./testdata/" + strings.ToLower(interfaceName) + ".go"

		outputFile, err := ioutil.ReadFile(outputFilename)
		if err != nil {
			t.Fatalf("ReadFile error: %s", err)
		}

		g, err := LoadPackageFiles([]string{inputFilename})
		if err != nil {
			t.Fatalf("parsePackage error: %s", err)
		}
		got, err := g.Generate([]string{interfaceName})
		if err != nil {
			t.Fatalf("Generator.Generate error for %s: %s", interfaceName, err)
		}

		// reset gensyms
		symGen.Reset()
		identSymGen.Reset()

		readableOutput := string(outputFile)
		readableResult := string(got)

		// Only compare everything after the first line to avoid
		// comparing the generation commands
		outputStart := strings.Index(readableOutput, "\n")
		resultStart := strings.Index(readableResult, "\n")

		diffs := dmp.DiffMain(readableOutput, readableResult, false)

		assert.Equal(t, readableOutput[outputStart:], readableResult[resultStart:], dmp.DiffPrettyText(diffs))
	}
}
