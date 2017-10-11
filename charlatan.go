package main // import "github.com/percolate/charlatan"

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
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

var (
	interfaces  stringSliceValue
	outputName  = flag.String("output", "", "output file name [default: charlatan.go]")
	packageName = flag.String("package", "", "output package name [default: \"<current package>\"]")
)

func init() {
	flag.Var(&interfaces, "interface", "name of interface to fake, may be repeated")
	log.SetFlags(0)
	log.SetPrefix("charlatan: ")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `charlatan
https://github.com/percolate/charlatan

Usage:
  charlatan [options] (--interface <I>)...
  charlatan -h | --help

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if len(interfaces) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var (
		dir string
		g   Generator
	)

	g.interfaceNames = interfaces
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
		g.parsePackageDir(args[0])
	} else {
		dir = filepath.Dir(args[0])
		g.parsePackageFiles(args)
	}

	g.generate()

	// format the output.
	src := g.format()

	// Write to file.
	outputName := *outputName
	if outputName == "" {
		outputName = filepath.Join(dir, "charlatan.go")
	}

	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf            bytes.Buffer // Accumulated output.
	pkg            *Package     // Package we are scanning.
	interfaceNames []string
	imports        *ImportSet
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// File holds a single parsed file and associated data.
type File struct {
	pkg            *Package                // Package to which this file belongs.
	file           *ast.File               // Parsed AST.
	interfaces     []*InterfaceDeclaration // The interface declarations.
	interfaceNames []string
	imports        *ImportSet
}

type Package struct {
	dir      string
	name     string
	defs     map[*ast.Ident]types.Object
	files    []*File
	typesPkg *types.Package
}

// An Import is a struct used to track qualified types across files
// so that we can import their packages in the charlatan file
type Import struct {
	Name     string // the package's name
	Path     string // import path for the package
	Required bool   // denotes whether the import is required by charlatan
}

type ImportSet struct {
	// possibly makes more sense to use a `map` underneath?
	imports []*Import
}

func (r *ImportSet) Add(ri *Import) {
	if !r.Contains(ri) {
		r.imports = append(r.imports, ri)
	}
}

func (r *ImportSet) Remove(ri *Import) error {
	for index, i := range r.imports {
		found := i.Name == ri.Name && i.Path == ri.Path
		if found {
			r.imports = append(r.imports[:index], r.imports[index+1:]...)
			return nil
		}
	}
	return errors.New("Item not present in set, cannot remove")
}

func (r *ImportSet) Contains(ri *Import) bool {
	for _, i := range r.imports {
		found := i.Name == ri.Name && i.Path == ri.Path
		if found {
			return found
		}
	}
	return false
}

func (r *ImportSet) GetAll() []*Import {
	return r.imports
}

func (r *ImportSet) GetRequired() []*Import {
	reqimps := make([]*Import, 0)
	for _, imp := range r.imports {
		if imp.Required {
			reqimps = append(reqimps, imp)
		}
	}
	return reqimps
}

func (r *ImportSet) RequireByName(s string) {
	for i, imp := range r.imports {
		if imp.Name == s {
			r.imports[i].Required = true
		}
	}
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

// check type-checks the package. The package must be OK to proceed.
func (pkg *Package) check(fs *token.FileSet, astFiles []*ast.File) {
	pkg.defs = make(map[*ast.Ident]types.Object)
	config := types.Config{Importer: defaultImporter(), FakeImportC: true}
	info := &types.Info{
		Defs: pkg.defs,
	}
	typesPkg, err := config.Check(pkg.dir, fs, astFiles, info)
	if err != nil {
		log.Fatalf("checking package: %s", err)
	}
	pkg.typesPkg = typesPkg
}

// generate produces the charlatan file for the named interface.
func (g *Generator) generate() {
	interfacedecs := make([]*InterfaceDeclaration, 0, 100)
	g.imports = &ImportSet{
		imports: make([]*Import, 0),
	}
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.

		if file.file != nil {
			file.imports = g.imports
			file.interfaceNames = g.interfaceNames
			ast.Inspect(file.file, file.genDecl)

			interfacedecs = append(interfacedecs, file.interfaces...)
		}
	}

	if len(interfacedecs) == 0 {
		log.Fatalf("no interfaces named %s defined", g.interfaceNames)
	}

	packageName := *packageName
	if packageName == "" {
		packageName = fmt.Sprintf("%stest", g.pkg.name)
	}

	g.Printf("// generated by \"charlatan %s\"; DO NOT EDIT.\n\n", strings.Join(os.Args[1:], " "))
	g.Printf("package %s\n\n", packageName)

	allimps := g.imports.GetRequired()

	if len(allimps) == 1 {
		g.Printf("import %s", allimps[0].Path)
	} else if len(allimps) > 1 {
		g.Printf("import (\n")
		for _, imp := range allimps {
			g.Printf("\t%s\n", imp.Path)
		}
		g.Printf(")\n")
	}

	g.Printf("\n")

	for _, i := range interfacedecs {
		title := fmt.Sprintf("// This is a mock for %s //", i.Name)
		bar := strings.Repeat("/", len(title))

		g.Printf("%s\n", bar)
		g.Printf("%s\n", title)
		g.Printf("%s\n", bar)

		for _, m := range i.Methods {
			if err := invocationTempl.Execute(&g.buf, m); err != nil {
				panic(err)
			}
		}

		if err := fakeTempl.Execute(&g.buf, i); err != nil {
			panic(err)
		}

		for _, m := range i.Methods {
			if err := methodTempl.Execute(&g.buf, m); err != nil {
				panic(err)
			}
		}
	}

}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

// InterfaceDeclaration represents a declared interface.
type InterfaceDeclaration struct {
	Name    string
	Methods []*Method
}

// Method represents a method in an interface's method set
type Method struct {
	InterfaceName string
	Name          string
	Params        []*Value
	Results       []*Value
}

// Value represents a Parameter or Result of a Method
type Value struct {
	Name       string // name if a named parameter/result, else null string
	Type       string //
	Qualifier  string
	Pointer    bool
	Elliptical bool
	ChanDir    int
}

func (m Method) FormatParamsDeclaration() string {
	var f []string
	for _, v := range m.Params {
		f = append(f, v.functionDeclarationFormat())
	}
	return strings.Join(f, ", ")
}

func (m Method) FormatParamsCall() string {
	var f []string
	for _, v := range m.Params {
		f = append(f, v.argumentFormat())
	}
	return strings.Join(f, ", ")
}

func (m Method) FormatResultsDeclaration() string {
	var f []string
	for _, v := range m.Results {
		f = append(f, v.functionDeclarationFormat())
	}
	return strings.Join(f, ", ")
}

func (m Method) FormatResultsCall() string {
	var f []string
	for _, v := range m.Results {
		f = append(f, v.argumentFormat())
	}
	return strings.Join(f, ", ")
}

func (v Value) functionDeclarationFormat() string {
	formatted := ""
	switch v.ChanDir {
	case 1:
		formatted = "chan<- "
	case 2:
		formatted = "<-chan "
	case 3:
		formatted = "chan "
	}
	if v.Elliptical {
		formatted = fmt.Sprintf("%s...", formatted)
	}
	if v.Pointer {
		formatted = fmt.Sprintf("%s*", formatted)
	}
	if len(v.Qualifier) > 0 {
		formatted = fmt.Sprintf("%s%s.", formatted, v.Qualifier)
	}
	return fmt.Sprintf("%s %s%s", v.Name, formatted, v.Type)
}

func (v Value) argumentFormat() string {
	formatted := ""
	if v.Elliptical {
		formatted = "..."
	}
	return fmt.Sprintf("%s%s", v.Name, formatted)
}

func (v Value) StructDef() string {
	formatted := ""
	switch v.ChanDir {
	case 1:
		formatted = "chan<- "
	case 2:
		formatted = "<-chan "
	case 3:
		formatted = "chan "
	}
	if v.Elliptical {
		formatted = fmt.Sprintf("%s[]", formatted)
	}
	if v.Pointer {
		formatted = fmt.Sprintf("%s*", formatted)
	}
	if len(v.Qualifier) > 0 {
		formatted = fmt.Sprintf("%s%s.", formatted, v.Qualifier)
	}
	return fmt.Sprintf("%s %s%s", v.CapitalName(), formatted, v.Type)
}

func (v Value) CapitalName() string {
	return strings.Title(v.Name)
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.TYPE {
		// We only care about type declarations.
		if ok && decl.Tok == token.IMPORT {
			for _, s := range decl.Specs {
				spec := s.(*ast.ImportSpec)

				// Only add un-named imports for now
				if spec.Name == nil {
					parts := strings.Split(spec.Path.Value, "/")
					name := strings.Replace(parts[len(parts)-1], "\"", "", -1)
					imp := &Import{
						Name:     name,
						Path:     spec.Path.Value,
						Required: false,
					}
					f.imports.Add(imp)
				}
			}
		}
		return true
	}

	spec := decl.Specs[0].(*ast.TypeSpec)
	ident := spec.Name

	// Look for an interface type with methods, not named `_`
	specType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		// We only care about interfaces with methods
		return true
	}
	methods := specType.Methods.List
	if len(methods) == 0 {
		return true
	}
	name := ident.Name
	if name == "_" {
		return true
	}

	// Continue walking if the name doesn't match a name we're looking for
	namefound := false
	for _, i := range f.interfaceNames {
		namefound = i == name
		if namefound {
			break
		}
	}
	if !namefound {
		return true
	}

	interfacedec := &InterfaceDeclaration{
		Name: name,
	}

	// Add each method to our interfacedec
	for _, method := range methods {
		_, ok := method.Type.(*ast.FuncType)
		if ok {
			interfacedec.addMethod(method, f.imports)
		}
	}

	f.interfaces = append(f.interfaces, interfacedec)

	return false
}

// Maps an ast.Field reference to an array of Values
func extractValues(f *ast.Field, i int, prefix string, imports *ImportSet) []*Value {
	var values []*Value
	var names []string

	elliptical := false
	chandir := 0
	ispointer := false
	qualifier := ""
	fieldType := ""

	if len(f.Names) == 0 {
		n := fmt.Sprintf("%s%d", prefix, i)
		names = append(names, n)
	} else {
		for _, n := range f.Names {
			names = append(names, n.Name)
		}
	}

	// Check if we're dealing with an ellipse
	topType := f.Type
	ellipsis, ok := topType.(*ast.Ellipsis)
	if ok {
		elliptical = true
		topType = ellipsis.Elt
	}

	// Check if we're dealing with a channel
	chantype, ok := topType.(*ast.ChanType)
	if ok {
		chandir = int(chantype.Dir)
		topType = chantype.Value
	}

	// Check if we're dealing with a pointer
	starType, ok := topType.(*ast.StarExpr)
	if ok {
		ispointer = true
		topType = starType.X
	}

	typeParseFailure := "charlatan: failed to parse type: %s"

	// Check if the type is a qualified identifier (from a package)
	selectorType, isqualified := topType.(*ast.SelectorExpr)
	_, isinterface := topType.(*ast.InterfaceType)
	_, isstruct := topType.(*ast.StructType)
	if isinterface {
		fieldType = "interface{}"
	} else if isstruct {
		fieldType = "struct{}"
	} else if isqualified {
		selectedName, ok := selectorType.X.(*ast.Ident)
		if !ok {
			fmt.Println(fmt.Errorf(typeParseFailure, f.Type))
		}
		qualifier = selectedName.Name
		imports.RequireByName(selectedName.Name)

		fieldType = selectorType.Sel.Name
	} else {
		selectedName, ok := topType.(*ast.Ident)
		if !ok {
			fmt.Println(fmt.Errorf(typeParseFailure, f.Type))
		}
		fieldType = selectedName.Name
	}

	for _, name := range names {
		v := &Value{
			Name:       name,
			Type:       fieldType,
			Pointer:    ispointer,
			Elliptical: elliptical,
			Qualifier:  qualifier,
			ChanDir:    chandir,
		}
		values = append(values, v)
	}
	return values
}

func (i *InterfaceDeclaration) addMethod(m *ast.Field, imps *ImportSet) {
	functype, ok := m.Type.(*ast.FuncType)
	if !ok {
		return
	}

	method := &Method{
		InterfaceName: i.Name,
		Name:          m.Names[0].Name,
	}

	// `Params.List` can be 0-length, but `Results` can be nil
	for i, p := range functype.Params.List {
		values := extractValues(p, i, "arg", imps)
		method.Params = append(method.Params, values...)
	}

	if functype.Results != nil {
		for i, r := range functype.Results.List {
			values := extractValues(r, i, "ret", imps)
			method.Results = append(method.Results, values...)
		}
	}
	i.Methods = append(i.Methods, method)
}
