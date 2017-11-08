package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	importer := defaultImporter()
	for _, filename := range filenames {
		if !strings.HasSuffix(filename, ".go") {
			continue
		}
		file, err := parser.ParseFile(fileset, filename, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("syntax error: %s", err)
		}
		if err := generator.extractImports(file, importer); err != nil {
			return nil, err
		}
		if err := generator.extractInterfaces(file); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("error: no Go files found in %s", directory)
	}
	generator.packageName = files[0].Name.Name

	// Type check the package.
	config := types.Config{Importer: importer, Error: func(err error) { fmt.Fprintln(os.Stderr, err) }}
	if _, err := config.Check(directory, fileset, files, nil); err != nil {
		return nil, fmt.Errorf("type check failed")
	}

	return generator, nil
}

func (g *Generator) extractImports(file *ast.File, importer types.Importer) error {
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return err
		}
		pkg, err := importer.Import(path)
		if err != nil {
			return err
		}
		decl := &Import{
			Name: pkg.Name(),
			Path: spec.Path.Value,
		}

		if spec.Name == nil {
			g.imports.Add(decl)
			continue
		}

		switch spec.Name.Name {
		case "_":
			continue
		case ".":
			decl.Required = true
			decl.Alias = "."
		default:
			decl.Alias = spec.Name.Name
		}
		g.imports.Add(decl)
	}

	return nil
}

func (g *Generator) extractInterfaces(file *ast.File) (err error) {
	for _, node := range file.Decls {
		gen, ok := node.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		spec := gen.Specs[0].(*ast.TypeSpec)
		ifType, ok := spec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		decl := &Interface{
			Name: spec.Name.Name,
		}
		g.interfaces[spec.Name.Name] = decl

		for _, method := range ifType.Methods.List {
			if _, ok := method.Type.(*ast.FuncType); ok {
				err = decl.addMethod(method, g.imports)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

// Generate produces the charlatan source file data for the named interfaces.
func (g *Generator) Generate(interfaceNames []string) ([]byte, error) {
	decls := make([]*Interface, 0, len(interfaceNames))
	for _, name := range interfaceNames {
		decl, ok := g.interfaces[name]
		if !ok {
			return nil, fmt.Errorf("error: interface %q not found", name)
		}
		if len(decl.Methods) == 0 {
			log.Printf("warning: ignoring empty interface %q\n", decl.Name)
			continue
		}
		if decl.Name == "_" {
			log.Println(`warning: ignorning interface named "_"`)
			continue
		}
		decls = append(decls, decl)
	}

	if len(decls) == 0 {
		return nil, fmt.Errorf("error: no valid interface names provided")
	}

	packageName := g.packageName
	if g.PackageOverride != "" {
		packageName = g.PackageOverride
	}

	argv := []string{"charlatan"}
	tmpl := Template{
		CommandLine: strings.Join(append(argv, os.Args[1:]...), " "),
		PackageName: packageName,
		Imports:     g.imports.GetRequired(),
		Interfaces:  decls,
	}

	return tmpl.Execute()
}
