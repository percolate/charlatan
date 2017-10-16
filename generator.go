package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Generator holds the state of the analysis
type Generator struct {
	// PackageOverride can be set to control the package for the output file.  The default is the same package as the input interfaces.
	PackageOverride string
	pkg             *Package
	imports         *ImportSet
}

// LoadPackageDir parses the package residing in the given directory.
func LoadPackageDir(directory string) (*Generator, error) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot process directory %s: %s", directory, err)
	}
	names := make([]string, 0, len(pkg.GoFiles)+len(pkg.CgoFiles))
	names = append(names, pkg.GoFiles...)
	names = append(names, pkg.CgoFiles...)
	if directory != "." {
		for i, name := range names {
			names[i] = filepath.Join(directory, name)
		}
	}

	return parsePackage(directory, names)
}

// LoadPackageFiles parses the package using only the given files.
func LoadPackageFiles(names []string) (*Generator, error) {
	return parsePackage(".", names)
}

// parsePackage analyzes the single package constructed from the named files.
func parsePackage(directory string, names []string) (*Generator, error) {
	var files []*File
	var astFiles []*ast.File
	g := new(Generator)
	g.pkg = new(Package)
	fs := token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fs, name, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parsing package: %s: %s", name, err)
		}
		astFiles = append(astFiles, parsedFile)
		files = append(files, &File{
			file: parsedFile,
			pkg:  g.pkg,
		})
	}
	if len(astFiles) == 0 {
		return nil, fmt.Errorf("%s: no buildable Go files", directory)
	}
	g.pkg.name = astFiles[0].Name.Name
	g.pkg.files = files
	g.pkg.dir = directory
	// Type check the package.
	g.pkg.check(fs, astFiles)

	return g, nil
}

// generate produces the charlatan file for the named interface.
func (g *Generator) Generate(interfaces []string) ([]byte, error) {
	interfacedecs := make([]*InterfaceDeclaration, 0, 100)
	g.imports = &ImportSet{
		imports: make([]*Import, 0),
	}
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.

		if file.file != nil {
			file.imports = g.imports
			file.interfaceNames = interfaces
			ast.Inspect(file.file, file.genDecl)

			interfacedecs = append(interfacedecs, file.interfaces...)
		}
	}

	if len(interfacedecs) == 0 {
		return nil, fmt.Errorf("no interfaces named %s defined", interfaces)
	}

	packageName := g.pkg.name
	if g.PackageOverride != "" {
		packageName = g.PackageOverride
	}

	requiredPackages := g.imports.GetRequired()
	imports := make([]string, len(requiredPackages))
	for i, pkg := range requiredPackages {
		imports[i] = pkg.Path
	}

	argv := []string{"charlatan"}
	tmpl := Template{
		CommandLine: strings.Join(append(argv, os.Args[1:]...), " "),
		PackageName: packageName,
		Imports:     imports,
		Interfaces:  interfacedecs,
	}

	return tmpl.Execute()
}
