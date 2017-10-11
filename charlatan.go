package main // import "github.com/percolate/charlatan"

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
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

// isDirectory returns true if the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
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
