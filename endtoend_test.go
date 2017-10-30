// +build !android

package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// This file contains a test that compiles and runs each program in testdata
// after generating the mocks for its interface. The rule is that for
// testdata/x.go we run `charlatan -file=testdata/x_def.go X` and then compile
// and run the testdata/x.go program. The resulting binary panics if the mock
// structs are broken, including for error cases.

func TestEndToEnd(t *testing.T) {
	dir, err := ioutil.TempDir("", "charlatan")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	// Create charlatan in temporary directory.
	charlatan := filepath.Join(dir, "charlatan.exe")
	err = run("go", "build", "-o", charlatan)
	if err != nil {
		t.Fatalf("building charlatan: %s", err)
	}
	// Read the testdata directory.
	fd, err := os.Open("testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()
	names, err := fd.Readdirnames(-1)
	if err != nil {
		t.Fatalf("Readdirnames: %s", err)
	}
	// Generate, compile, and run the test programs.
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			t.Errorf("%s is not a Go file", name)
			continue
		}

		if strings.HasSuffix(name, "_ete.go") {
			base := strings.TrimSuffix(name, "_ete.go")
			interfaceName := strings.Title(base)
			defName := base + "_def.go"

			charlatanCompileAndRun(t, dir, charlatan, interfaceName, defName, name)
		}
	}
}

// charlatanCompileAndRun runs charlatan for the named file and compiles and
// runs the target binary in directory dir. That binary will panic if the mock is broken.
func charlatanCompileAndRun(t *testing.T, dir, charlatan, interfaceName, defName, fileName string) {
	t.Logf("run: %s %s\n", fileName, interfaceName)
	source := filepath.Join(dir, fileName)
	sourceDef := filepath.Join(dir, defName)
	err := copy(source, filepath.Join("testdata", fileName))
	if err != nil {
		t.Fatalf("copying file to temporary directory: %s", err)
	}
	err = copy(sourceDef, filepath.Join("testdata", defName))
	if err != nil {
		t.Fatalf("copying file to temporary directory: %s", err)
	}
	charlatanSource := filepath.Join(dir, interfaceName+"_charlatan.go")
	// Run charlatan in temporary directory.
	err = run(charlatan, "-file", sourceDef, "-output", charlatanSource, "-package", "main", interfaceName)
	if err != nil {
		t.Fatal(err)
	}
	// Run the binary in the temporary directory.
	err = run("go", "run", charlatanSource, sourceDef, source)
	if err != nil {
		t.Fatal(err)
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
