package main

import (
	"flag"
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
		if err := generator.processImports(file, importer); err != nil {
			return nil, err
		}
		if err := generator.processInterfaces(file); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("error: no Go files found in %s", directory)
	}

	// N.B. - type check the package
	config := types.Config{Importer: importer, Error: func(err error) { fmt.Fprintln(os.Stderr, err) }}
	pkg, err := config.Check(directory, fileset, files, nil)
	if err != nil {
		return nil, fmt.Errorf("type check failed")
	}

	generator.packageName = pkg.Name()

	return generator, nil
}

// Generator holds the state of the analysis
type Generator struct {
	// PackageOverride can be set to control the package for the output file.  The default is the same package as the input interface(s).
	PackageOverride string
	packageName     string
	imports         *ImportSet
	interfaces      map[string]*Interface
}

func (g *Generator) processImports(file *ast.File, importer types.Importer) error {
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return err
		}
		pkg, err := importer.Import(path)
		if err != nil {
			return err
		}

		g.processImport(spec, pkg)
		if err := g.processImportInterfaces(pkg); err != nil {
			return err
		}

	}

	return nil
}

func (g *Generator) processImport(spec *ast.ImportSpec, pkg *types.Package) {
	decl := &Import{
		Name: pkg.Name(),
		Path: spec.Path.Value,
	}

	if spec.Name == nil {
		g.imports.Add(decl)
		return
	}

	switch spec.Name.Name {
	case "_":
		break
	case ".":
		decl.Required = true
		decl.Alias = "."
	default:
		decl.Alias = spec.Name.Name
	}

	g.imports.Add(decl)
}

func (g *Generator) processImportInterfaces(pkg *types.Package) error {
	for _, name := range pkg.Scope().Names() {
		obj := pkg.Scope().Lookup(name)

		qname := fmt.Sprintf("%s.%s", pkg.Name(), obj.Name())
		if _, exists := g.interfaces[qname]; exists {
			continue
		}

		if _, isType := obj.(*types.TypeName); !isType || !obj.Exported() || !types.IsInterface(obj.Type()) {
			continue
		}

		ifType := obj.Type().Underlying().(*types.Interface)
		decl := &Interface{
			Name: obj.Name(),
		}

		for i := 0; i < ifType.NumMethods(); i++ {
			m := ifType.Method(i)
			if !m.Exported() {
				continue
			}
			if err := decl.addMethodFromType(m, g.imports); err != nil {
				return err
			}
		}

		g.interfaces[qname] = decl
	}

	return nil
}

func (g *Generator) processInterfaces(file *ast.File) error {
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

		decl, err := g.processInterface(spec.Name.Name, ifType)
		if err != nil {
			return err
		}
		g.interfaces[spec.Name.Name] = decl
	}

	return nil
}

func (g *Generator) processInterface(name string, ifType *ast.InterfaceType) (*Interface, error) {
	decl := &Interface{
		Name: name,
	}

	for _, field := range ifType.Methods.List {
		switch f := field.Type.(type) {
		case *ast.BinaryExpr:
			// N.B. - type expression
			continue
		case *ast.FuncType:
			if err := decl.addMethodFromField(field, g.imports); err != nil {
				return nil, err
			}
		case *ast.Ident:
			// N.B. - embedded interface from current package
			decl.embeds = append(decl.embeds, f.Name)
		case *ast.SelectorExpr:
			// N.B. - embedded interface from imported package
			decl.embeds = append(decl.embeds, fmt.Sprintf("%s.%s", f.X.(*ast.Ident).String(), f.Sel.String()))
		default:
			return nil, fmt.Errorf("internal error: unsupported interface field: %t, %#v", f, field.Type)
		}
	}

	return decl, nil
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
		if len(decl.embeds) == 0 {
			continue
		}

		embeddedMethods := []*Method{}
		for _, embedName := range decl.embeds {
			embed, ok := g.interfaces[embedName]
			if !ok {
				return nil, fmt.Errorf("error: interface %q embedded in %s not found", embedName, name)
			}

			for _, m := range embed.Methods {
				c := *m
				c.Interface = decl.Name
				embeddedMethods = append(embeddedMethods, &c)
			}
		}
		decl.Methods = append(embeddedMethods, decl.Methods...)
	}

	if len(decls) == 0 {
		return nil, fmt.Errorf("error: no valid interface names provided")
	}

	packageName := g.packageName
	if g.PackageOverride != "" {
		packageName = g.PackageOverride
	}

	var argv strings.Builder
	argv.WriteString("charlatan")
	flag.Visit(func(f *flag.Flag) {
		fmt.Fprintf(&argv, " -%s=%s", f.Name, f.Value)
	})
	if flag.NArg() > 0 {
		argv.WriteByte(' ')
		argv.WriteString(strings.Join(flag.Args(), " "))
	}
	tmpl := charlatanTemplate{
		CommandLine: argv.String(),
		PackageName: packageName,
		Imports:     g.imports.GetRequired(),
		Interfaces:  decls,
	}

	return tmpl.execute()
}
