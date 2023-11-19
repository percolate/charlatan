//go:build !android
// +build !android

package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

// This file contains a test that compiles and runs each program in testdata
// after generating the mocks for its interface. The rule is that for
// testdata/x.go we run `charlatan -dir=testdata X` and then compile
// and run the testdata/x.go program. The resulting binary panics if the mock
// structs are broken, including for error cases.

type endToEndTest struct {
	exe  string
	file string
}

func (e *endToEndTest) compileAndRun(t *testing.T) {
	t.Parallel()
	tempdir, err := ioutil.TempDir("", "charlatan")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	base := strings.TrimSuffix(path.Base(e.file), "_ete.go")
	interfaceName := strings.Title(base)

	sourceDef := filepath.Join(tempdir, base+"_def.go")
	err = copy(sourceDef, filepath.Join("testdata/"+base, base+"_def.go"))
	if err != nil {
		t.Fatalf("copying interface definition file to temporary directory: %s", err)
	}

	charlatanSource := filepath.Join(tempdir, interfaceName+"_charlatan.go")
	// Run charlatan in temporary directory.
	err = run(e.exe, "-dir", tempdir, "-output", charlatanSource, "-package", "main", interfaceName)
	if err != nil {
		t.Fatal(err)
	}

	source := filepath.Join(tempdir, path.Base(e.file))
	err = copy(source, e.file)
	if err != nil {
		t.Fatalf("copying end-to-end test file to temporary directory: %s", err)
	}

	// Run the binary in the temporary directory.
	err = run("go", "run", charlatanSource, sourceDef, source)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEndToEnd(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "charlatan")
	if err != nil {
		t.Fatal(err)
	}

	// Create charlatan in temporary directory.
	charlatan := filepath.Join(tempdir, "charlatan.exe")
	err = run("go", "build", "-o", charlatan)
	if err != nil {
		t.Fatalf("building charlatan: %s", err)
	}

	names, err := filepath.Glob("testdata/ete/*_ete.go")
	if err != nil {
		t.Fatalf("finding end-to-end test files: %s", err)
	}

	for _, name := range names {
		e2e := endToEndTest{charlatan, name}
		t.Run(path.Base(name), e2e.compileAndRun)
	}
}

// copy copies the from file to the to file.
func copy(to, from string) error {
	toFd, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFd.Close()
	fromFd, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFd.Close()
	_, err = io.Copy(toFd, fromFd)
	return err
}

// run runs a single command and returns an error if it does not succeed.
// os/exec should have this function, to be honest.
func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
