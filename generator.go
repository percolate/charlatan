package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
)

// Generator holds the state of the analysis
type Generator struct {
	// PackageOverride can be set to control the package for the output file.  The default is the same package as the input interface(s).
	PackageOverride string
	packageName     string
	imports         *ImportSet
	interfaces      map[string]*Interface
}

// LoadPackageDir parses a package in the given directory.
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

// LoadPackageFiles parses a package using only the given files.
func LoadPackageFiles(names []string) (*Generator, error) {
	return parsePackage(".", names)
}

func parsePackage(directory string, filenames []string) (*Generator, error) {
	generator := &Generator{
		imports:    new(ImportSet),
		interfaces: make(map[string]*Interface),
	}
	files := make([]*ast.File, 0, len(filenames))
	fileset := token.NewFileSet()
	for _, filename := range filenames {
		if !strings.HasSuffix(filename, ".go") {
			continue
		}
		file, err := parser.ParseFile(fileset, filename, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("syntax error: %s", err)
		}
		generator.extractImports(file)
		generator.extractInterfaces(file)
		files = append(files, file)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("error: no Go files found in %s", directory)
	}
	generator.packageName = files[0].Name.Name

	// Type check the package.
	config := types.Config{Importer: defaultImporter(), Error: func(err error) { fmt.Fprintln(os.Stderr, err) }}
	if _, err := config.Check(directory, fileset, files, nil); err != nil {
		return nil, fmt.Errorf("type check failed")
	}

	return generator, nil
}

func (g *Generator) extractImports(file *ast.File) {
	for _, spec := range file.Imports {
		if spec.Name != nil {
			// Only add un-named imports for now
			// xxx - why?
			continue
		}
		parts := strings.Split(spec.Path.Value, "/")
		g.imports.Add(&Import{
			Name: strings.Replace(parts[len(parts)-1], "\"", "", -1),
			Path: spec.Path.Value,
		})
	}
}

func (g *Generator) extractInterfaces(file *ast.File) {
	for _, node := range file.Decls {
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			continue
		}
		spec := decl.Specs[0].(*ast.TypeSpec)
		ifType, ok := spec.Type.(*ast.InterfaceType)
		if !ok || len(ifType.Methods.List) == 0 || spec.Name.Name == "_" {
			continue
		}

		ifDecl := &Interface{
			Name: spec.Name.Name,
		}
		g.interfaces[spec.Name.Name] = ifDecl

		for _, method := range ifType.Methods.List {
			if _, ok := method.Type.(*ast.FuncType); ok {
				ifDecl.addMethod(method, g.imports)
			}
		}
	}
}

// Generate produces the charlatan source file data for the named interfaces.
func (g *Generator) Generate(interfaceNames []string) ([]byte, error) {
	decls := make([]*Interface, 0, len(interfaceNames))
	for _, name := range interfaceNames {
		decl, ok := g.interfaces[name]
		if !ok {
			return nil, fmt.Errorf("error: interface %q not found", name)
		}
		decls = append(decls, decl)
	}

	packageName := g.packageName
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
		Interfaces:  decls,
	}

	return tmpl.Execute()
}
