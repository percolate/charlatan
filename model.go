package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"strconv"
	"strings"
)

var (
	identSymGen = symbolGenerator{Prefix: "ident"}
)

// Import represents a declared import
type Import struct {
	Name     string // the package's name
	Alias    string // the local alias for the package name
	Path     string // import path for the package
	Required bool   // is the import required in the charlatan output?
}

// ImportSet contains all the import declarations encountered
type ImportSet struct {
	imports []*Import
}

// Add inserts the given value into the set if it doesn't already exist
func (r *ImportSet) Add(value *Import) {
	if r.imports == nil {
		r.imports = []*Import{value}
	} else if !r.Contains(value) {
		r.imports = append(r.imports, value)
	}
}

// Contains returns true if the given value is in the set
func (r *ImportSet) Contains(value *Import) bool {
	for _, i := range r.imports {
		if i.Name == value.Name && i.Path == value.Path {
			return true
		}
	}

	return false
}

// GetRequired returns all imports referenced by the target interface
func (r *ImportSet) GetRequired() []*Import {
	result := make([]*Import, 0, len(r.imports))
	for _, imp := range r.imports {
		if imp.Required {
			result = append(result, imp)
		}
	}
	return result
}

// RequireByName marks an import symbol as required
func (r *ImportSet) RequireByName(s string) {
	for i, imp := range r.imports {
		if imp.Name == s || imp.Alias == s {
			r.imports[i].Required = true
		}
	}
}

// Interface represents a declared interface.
type Interface struct {
	Name    string
	Methods []*Method
	embeds  []string
}

func (i *Interface) addMethodFromField(field *ast.Field, imports *ImportSet) error {
	functionType, ok := field.Type.(*ast.FuncType)
	if !ok {
		return fmt.Errorf("internal error: expected *ast.FuncType, have: %#v", field)
	}

	method := &Method{
		Interface: i.Name,
		Name:      field.Names[0].Name,
	}

	identSymGen.reset()
	// `Params.List` can be 0-length, but `Results` can be nil
	for _, parameter := range functionType.Params.List {
		identifiers, err := extractIdentifiersFromField(parameter, imports)
		if err != nil {
			return err
		}
		method.Parameters = append(method.Parameters, identifiers...)
	}

	if functionType.Results != nil {
		for _, result := range functionType.Results.List {
			identifiers, err := extractIdentifiersFromField(result, imports)
			if err != nil {
				return err
			}
			method.Results = append(method.Results, identifiers...)
		}
	}

	i.Methods = append(i.Methods, method)
	return nil
}

func extractIdentifiersFromField(field *ast.Field, imports *ImportSet) ([]*Identifier, error) {
	identifierType, err := unwrapExpr(field.Type, imports)
	if err != nil {
		return nil, err
	}

	if len(field.Names) == 0 {
		return []*Identifier{
			{
				Name:      identSymGen.next(),
				ValueType: identifierType,
			},
		}, nil
	}

	identifiers := make([]*Identifier, len(field.Names))
	for i, name := range field.Names {
		identifiers[i] = &Identifier{
			Name:      name.Name,
			ValueType: identifierType,
		}
	}

	return identifiers, nil
}

func unwrapExpr(node ast.Expr, imports *ImportSet) (t Type, err error) {
	switch nodeType := node.(type) {
	case *ast.Ellipsis:
		var subType Type
		subType, err = unwrapExpr(nodeType.Elt, imports)
		if err != nil {
			return
		}
		t = &Ellipsis{
			subType: subType,
		}
	case *ast.ArrayType:
		var subType Type
		subType, err = unwrapExpr(nodeType.Elt, imports)
		if err != nil {
			return
		}
		a := &Array{
			subType: subType,
		}
		if nodeType.Len != nil {
			if lit, ok := nodeType.Len.(*ast.BasicLit); ok {
				a.scale = lit.Value
			} else {
				err = fmt.Errorf("internal error: unsupported array len type node: %#v", nodeType.Len)
				return
			}
		}
		t = a
	case *ast.MapType:
		var keyType Type
		keyType, err = unwrapExpr(nodeType.Key, imports)
		if err != nil {
			return
		}
		var subType Type
		subType, err = unwrapExpr(nodeType.Value, imports)
		if err != nil {
			return
		}
		t = &Map{keyType: keyType, subType: subType}
	case *ast.ChanType:
		var subType Type
		subType, err = unwrapExpr(nodeType.Value, imports)
		if err != nil {
			return
		}
		switch nodeType.Dir {
		case ast.SEND:
			t = &SendChannel{
				subType: subType,
			}
		case ast.RECV:
			t = &ReceiveChannel{
				subType: subType,
			}
		case ast.SEND + ast.RECV:
			t = &Channel{
				subType: subType,
			}
		}
	case *ast.StarExpr:
		var subType Type
		subType, err = unwrapExpr(nodeType.X, imports)
		if err != nil {
			return
		}
		t = &Pointer{
			subType: subType,
		}
	case *ast.InterfaceType, *ast.StructType, *ast.FuncType:
		var buf bytes.Buffer
		if err = format.Node(&buf, token.NewFileSet(), nodeType); err != nil {
			return
		}
		t = &BasicType{
			Name: buf.String(),
		}
	case *ast.SelectorExpr:
		selector := nodeType.X.(*ast.Ident).Name
		imports.RequireByName(selector)
		t = &BasicType{
			Qualifier: selector,
			Name:      nodeType.Sel.Name,
		}
	case *ast.Ident:
		t = &BasicType{
			Name: nodeType.Name,
		}
	default:
		err = fmt.Errorf("internal error: unsupported parameter type: %#v", nodeType)
	}

	return
}

func (i *Interface) addMethodFromType(f *types.Func, imports *ImportSet) error {
	method := &Method{
		Interface: i.Name,
		Name:      f.Name(),
	}

	sig := f.Type().(*types.Signature)
	parameters, err := extractIdentifiersFromTuple(sig.Params(), imports)
	if err != nil {
		return err
	}
	if sig.Variadic() {
		last := parameters[len(parameters)-1]
		if avt, ok := last.ValueType.(*Array); ok {
			last.ValueType = &Ellipsis{subType: avt.subType}
		} else {
			last.ValueType = &Ellipsis{subType: last.ValueType}
		}
	}
	method.Parameters = append(method.Parameters, parameters...)

	results, err := extractIdentifiersFromTuple(sig.Results(), imports)
	if err != nil {
		return err
	}
	method.Results = append(method.Results, results...)

	i.Methods = append(i.Methods, method)

	return nil
}

func extractIdentifiersFromTuple(tuple *types.Tuple, imports *ImportSet) ([]*Identifier, error) {
	if 0 == tuple.Len() {
		return nil, nil
	}

	idents := make([]*Identifier, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		p := tuple.At(i)
		identifierType, err := unwrapType(p.Type(), imports)
		if err != nil {
			return nil, err
		}
		ident := &Identifier{
			Name:      p.Name(),
			ValueType: identifierType,
		}
		if "" == ident.Name {
			ident.Name = identSymGen.next()
		}
		idents[i] = ident
	}

	return idents, nil
}

func unwrapType(t types.Type, imports *ImportSet) (r Type, err error) {
	switch actual := t.(type) {
	case *types.Array:
		var subType Type
		subType, err = unwrapType(actual.Elem(), imports)
		if err != nil {
			return
		}
		a := &Array{subType: subType}
		if actual.Len() != 0 {
			a.scale = strconv.FormatInt(actual.Len(), 10)
		}
		r = a
	case *types.Slice:
		var subType Type
		subType, err = unwrapType(actual.Elem(), imports)
		if err != nil {
			return
		}
		r = &Array{subType: subType}
	case *types.Map:
		var keyType Type
		keyType, err = unwrapType(actual.Key(), imports)
		if err != nil {
			return
		}
		var subType Type
		subType, err = unwrapType(actual.Elem(), imports)
		if err != nil {
			return
		}
		r = &Map{keyType: keyType, subType: subType}
	case *types.Chan:
		var subType Type
		subType, err = unwrapType(actual.Elem(), imports)
		if err != nil {
			return
		}
		switch actual.Dir() {
		case types.SendOnly:
			r = &SendChannel{subType: subType}
		case types.RecvOnly:
			r = &ReceiveChannel{subType: subType}
		case types.SendRecv:
			r = &Channel{subType: subType}
		}
	case *types.Pointer:
		var subType Type
		subType, err = unwrapType(actual.Elem(), imports)
		if err != nil {
			return nil, err
		}
		r = &Pointer{subType: subType}
	case *types.Interface, *types.Struct, *types.Signature:
		r = &BasicType{Name: actual.String()}
	case *types.Named:
		b := &BasicType{Name: actual.Obj().Name()}
		if actual.Obj().Pkg() != nil {
			b.Qualifier = actual.Obj().Pkg().Name()
			imports.RequireByName(b.Qualifier)
		}
		r = b
	case *types.Basic:
		r = &BasicType{Name: actual.Name()}
	default:
		err = fmt.Errorf("internal error: unsupported parameter type: %#v", actual)
	}

	return
}

// Method represents a method in an interface's method set
type Method struct {
	Interface             string
	Name                  string
	Parameters            []*Identifier
	Results               []*Identifier
	parametersDeclaration string
	resultsDeclaration    string
	parametersCall        string
	resultsCall           string
	parametersSignature   string
	resultsSignature      string
}

// ParametersDeclaration returns the formal declaration syntax for the method's parameters
func (m *Method) ParametersDeclaration() string {
	if len(m.Parameters) == 0 {
		return ""
	}
	if m.parametersDeclaration == "" {
		idents := make([]string, len(m.Parameters))
		for i, ident := range m.Parameters {
			idents[i] = ident.ParameterFormat()
		}
		m.parametersDeclaration = strings.Join(idents, ", ")
	}

	return m.parametersDeclaration
}

// ResultsDeclaration returns the formal declaration syntax for the method's results
func (m *Method) ResultsDeclaration() string {
	if len(m.Results) == 0 {
		return ""
	}
	if m.resultsDeclaration == "" {
		idents := make([]string, len(m.Results))
		for i, ident := range m.Results {
			idents[i] = ident.ParameterFormat()
		}
		m.resultsDeclaration = strings.Join(idents, ", ")
	}

	return m.resultsDeclaration
}

// ParametersReference returns the sytax to reference the method's parameters
func (m *Method) ParametersReference() string {
	if len(m.Parameters) == 0 {
		return ""
	}
	if m.parametersCall == "" {
		idents := make([]string, len(m.Parameters))
		for i, ident := range m.Parameters {
			idents[i] = ident.ReferenceFormat()
		}
		m.parametersCall = strings.Join(idents, ", ")
	}

	return m.parametersCall
}

// ResultsReference returns the syntax to reference the method's results
func (m *Method) ResultsReference() string {
	if len(m.Results) == 0 {
		return ""
	}
	if m.resultsCall == "" {
		idents := make([]string, len(m.Results))
		for i, ident := range m.Results {
			idents[i] = ident.ReferenceFormat()
		}
		m.resultsCall = strings.Join(idents, ", ")
	}

	return m.resultsCall
}

// ParametersSignature returns the type declaration syntax for the methods parameters
func (m *Method) ParametersSignature() string {
	if len(m.Parameters) == 0 {
		return ""
	}
	if m.parametersSignature == "" {
		idents := make([]string, len(m.Parameters))
		for i, ident := range m.Parameters {
			idents[i] = ident.Signature()
		}
		m.parametersSignature = strings.Join(idents, ", ")
	}

	return m.parametersSignature
}

// ResultsSignature returns the type declaration syntax for the methods results
func (m *Method) ResultsSignature() string {
	if len(m.Results) == 0 {
		return ""
	}
	if m.resultsSignature == "" {
		idents := make([]string, len(m.Results))
		for i, ident := range m.Results {
			idents[i] = ident.Signature()
		}
		m.resultsSignature = strings.Join(idents, ", ")
	}

	return m.resultsSignature
}

// Identifier is a declared identifier
type Identifier struct {
	Name            string
	ValueType       Type
	titleCase       string
	parameterFormat string
	referenceFormat string
	fieldFormat     string
	signature       string
}

// TitleCase returns the identifier's name in title case
func (i *Identifier) TitleCase() string {
	if i.titleCase == "" {
		i.titleCase = strings.Title(i.Name)
	}
	return i.titleCase
}

// ParameterFormat returns the syntax to use the identifier as a parameter
func (i *Identifier) ParameterFormat() string {
	if i.parameterFormat == "" {
		i.parameterFormat = fmt.Sprintf("%s %s", i.Name, i.ValueType.ParameterFormat())
	}

	return i.parameterFormat
}

// ReferenceFormat returns the syntax to refer to the identifier
func (i *Identifier) ReferenceFormat() string {
	if i.referenceFormat == "" {
		i.referenceFormat = fmt.Sprintf("%s%s", i.Name, i.ValueType.ReferenceFormat())
	}

	return i.referenceFormat
}

// FieldFormat returns the syntax to use the identifier as a field
func (i *Identifier) FieldFormat() string {
	if i.fieldFormat == "" {
		i.fieldFormat = fmt.Sprintf("%s %s", i.TitleCase(), i.ValueType.FieldFormat())
	}

	return i.fieldFormat
}

// Signature returns the syntax type signature syntax for the identifier
func (i *Identifier) Signature() string {
	if i.signature == "" {
		i.signature = i.ValueType.ParameterFormat()
	}

	return i.signature
}

// Type is the interface for all built-in types
type Type interface {
	// ParameterFormat returns syntax for a parameter declaration
	ParameterFormat() string
	// ReferenceFormat returns the syntax for a reference
	ReferenceFormat() string
	// FieldFormat returns the syntax for a field declaration
	FieldFormat() string
}

// Array is the built-in array type
type Array struct {
	subType         Type
	scale           string
	parameterFormat string
	fieldFormat     string
}

// ParameterFormat returns syntax for a parameter declaration
func (t *Array) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("[%s]%s", t.scale, t.subType.ParameterFormat())

	return t.parameterFormat
}

// ReferenceFormat returns the syntax for a reference
func (t *Array) ReferenceFormat() string {
	return ""
}

// FieldFormat returns the syntax for a field declaration
func (t *Array) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("[%s]%s", t.scale, t.subType.FieldFormat())

	return t.fieldFormat
}

// Map is the built-in map type
type Map struct {
	keyType         Type
	subType         Type
	parameterFormat string
	fieldFormat     string
}

// ParameterFormat returns syntax for a parameter declaration
func (t *Map) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("map[%s]%s", t.keyType.ParameterFormat(), t.subType.ParameterFormat())

	return t.parameterFormat
}

// ReferenceFormat returns the syntax for a reference
func (t *Map) ReferenceFormat() string {
	return ""
}

// FieldFormat returns the syntax for a field declaration
func (t *Map) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("map[%s]%s", t.keyType.FieldFormat(), t.subType.FieldFormat())

	return t.fieldFormat
}

// Ellipsis is the built-in vararg type
type Ellipsis struct {
	subType         Type
	parameterFormat string
	fieldFormat     string
}

// ParameterFormat returns syntax for a parameter declaration
func (t *Ellipsis) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("...%s", t.subType.ParameterFormat())

	return t.parameterFormat
}

// ReferenceFormat returns the syntax for a reference
func (t *Ellipsis) ReferenceFormat() string {
	return "..."
}

// FieldFormat returns the syntax for a field declaration
func (t *Ellipsis) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("[]%s", t.subType.FieldFormat())

	return t.fieldFormat
}

// Channel is the built-in bidirectional channel
type Channel struct {
	subType         Type
	parameterFormat string
	fieldFormat     string
}

// ParameterFormat returns syntax for a parameter declaration
func (t *Channel) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("chan %s", t.subType.ParameterFormat())

	return t.parameterFormat
}

// ReferenceFormat returns the syntax for a reference
func (t *Channel) ReferenceFormat() string {
	return ""
}

// FieldFormat returns the syntax for a field declaration
func (t *Channel) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("chan %s", t.subType.FieldFormat())

	return t.fieldFormat
}

// ReceiveChannel is the built-in receive-only channel
type ReceiveChannel struct {
	subType         Type
	parameterFormat string
	fieldFormat     string
}

// ParameterFormat returns syntax for a parameter declaration
func (t *ReceiveChannel) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("<-chan %s", t.subType.ParameterFormat())

	return t.parameterFormat
}

// ReferenceFormat returns the syntax for a reference
func (t *ReceiveChannel) ReferenceFormat() string {
	return ""
}

// FieldFormat returns the syntax for a field declaration
func (t *ReceiveChannel) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("<-chan %s", t.subType.FieldFormat())

	return t.fieldFormat
}

// SendChannel is the built-in send-only channel
type SendChannel struct {
	subType         Type
	parameterFormat string
	fieldFormat     string
}

// ParameterFormat returns syntax for a parameter declaration
func (t *SendChannel) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("chan<- %s", t.subType.ParameterFormat())

	return t.parameterFormat
}

// ReferenceFormat returns the syntax for a reference
func (t *SendChannel) ReferenceFormat() string {
	return ""
}

// FieldFormat returns the syntax for a field declaration
func (t *SendChannel) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("chan<- %s", t.subType.FieldFormat())

	return t.fieldFormat
}

// Pointer is the built-in pointer type
type Pointer struct {
	subType         Type
	parameterFormat string
	fieldFormat     string
}

// ParameterFormat returns syntax for a parameter declaration
func (t *Pointer) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("*%s", t.subType.ParameterFormat())

	return t.parameterFormat
}

// ReferenceFormat returns the syntax for a reference
func (t *Pointer) ReferenceFormat() string {
	return t.subType.ReferenceFormat()
}

// FieldFormat returns the syntax for a field declaration
func (t *Pointer) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("*%s", t.subType.FieldFormat())

	return t.fieldFormat
}

// BasicType represents all built-in simple types
type BasicType struct {
	Name            string
	Qualifier       string
	parameterFormat string
	fieldFormat     string
}

// ParameterFormat returns syntax for a parameter declaration
func (t *BasicType) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	if t.Qualifier != "" {
		t.parameterFormat = fmt.Sprintf("%s.%s", t.Qualifier, t.Name)
	} else {
		t.parameterFormat = t.Name
	}

	return t.parameterFormat
}

// ReferenceFormat returns the syntax for a reference
func (t *BasicType) ReferenceFormat() string {
	return ""
}

// FieldFormat returns the syntax for a field declaration
func (t *BasicType) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	if t.Qualifier != "" {
		t.fieldFormat = fmt.Sprintf("%s.%s", t.Qualifier, t.Name)
	} else {
		t.fieldFormat = t.Name
	}

	return t.fieldFormat
}
