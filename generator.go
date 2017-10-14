package main

import (
	"bytes"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	pkg           *Package // Package we are scanning.
	targetPackage string
	interfaces    []string
	imports       *ImportSet
}

// parsePackageDir parses the package residing in the directory.
func (g *Generator) parsePackageDir(directory string) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		log.Fatalf("cannot process directory %s: %s", directory, err)
	}
	var names []string

	names = append(names, pkg.GoFiles...)
	names = append(names, pkg.CgoFiles...)
	names = append(names, pkg.SFiles...)
	names = prefixDirectory(directory, names)

	g.parsePackage(directory, names, nil)
}

// parsePackageFiles parses the package occupying the named files.
func (g *Generator) parsePackageFiles(names []string) {
	g.parsePackage(".", names, nil)
}

// parsePackage analyzes the single package constructed from the named files.
// If text is non-nil, it is a string to be used instead of the content of the file,
// to be used for testing. parsePackage exits if there is an error.
func (g *Generator) parsePackage(directory string, names []string, text interface{}) {
	var files []*File
	var astFiles []*ast.File
	g.pkg = new(Package)
	fs := token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fs, name, text, 0)
		if err != nil {
			log.Fatalf("parsing package: %s: %s", name, err)
		}
		astFiles = append(astFiles, parsedFile)
		files = append(files, &File{
			file: parsedFile,
			pkg:  g.pkg,
		})
	}
	if len(astFiles) == 0 {
		log.Fatalf("%s: no buildable Go files", directory)
	}
	g.pkg.name = astFiles[0].Name.Name
	g.pkg.files = files
	g.pkg.dir = directory
	// Type check the package.
	g.pkg.check(fs, astFiles)
}

// generate produces the charlatan file for the named interface.
func (g *Generator) generate() []byte {
	interfacedecs := make([]*InterfaceDeclaration, 0, 100)
	g.imports = &ImportSet{
		imports: make([]*Import, 0),
	}
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.

		if file.file != nil {
			file.imports = g.imports
			file.interfaceNames = g.interfaces
			ast.Inspect(file.file, file.genDecl)

			interfacedecs = append(interfacedecs, file.interfaces...)
		}
	}

	if len(interfacedecs) == 0 {
		log.Fatalf("no interfaces named %s defined", g.interfaces)
	}

	if g.targetPackage == "" {
		g.targetPackage = g.pkg.name
	}

	requiredPackages := g.imports.GetRequired()
	imports := make([]string, len(requiredPackages))
	for i, pkg := range requiredPackages {
		imports[i] = pkg.Path
	}

	argv := []string{"charlatan"}
	tmpl := Template{
		CommandLine: strings.Join(append(argv, os.Args[1:]...), " "),
		PackageName: g.targetPackage,
		Imports:     imports,
		Interfaces:  interfacedecs,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf); err != nil {
		log.Fatal(err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return buf.Bytes()
	}

	return src
}

// prefixDirectory places the directory name on the beginning of each name in the list.
func prefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}
