package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
)

var (
	identSymGen = SymbolGenerator{Prefix: "ident"}
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

func (r *ImportSet) Add(value *Import) {
	if r.imports == nil {
		r.imports = []*Import{value}
	} else if !r.Contains(value) {
		r.imports = append(r.imports, value)
	}
}

func (r *ImportSet) Contains(ri *Import) bool {
	for _, i := range r.imports {
		if i.Name == ri.Name && i.Path == ri.Path {
			return true
		}
	}

	return false
}

func (r *ImportSet) GetRequired() []*Import {
	result := make([]*Import, 0, len(r.imports))
	for _, imp := range r.imports {
		if imp.Required {
			result = append(result, imp)
		}
	}
	return result
}

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
}

func (i *Interface) addMethod(field *ast.Field, imports *ImportSet) error {
	functionType, ok := field.Type.(*ast.FuncType)
	if !ok {
		return fmt.Errorf("internal error: expected *ast.FuncType, have: %#v", field)
	}

	method := &Method{
		Interface: i.Name,
		Name:      field.Names[0].Name,
	}

	// `Params.List` can be 0-length, but `Results` can be nil
	for _, parameter := range functionType.Params.List {
		identifiers, err := extractIdentifiers(parameter, imports)
		if err != nil {
			return err
		}
		method.Parameters = append(method.Parameters, identifiers...)
	}

	if functionType.Results != nil {
		for _, result := range functionType.Results.List {
			identifiers, err := extractIdentifiers(result, imports)
			if err != nil {
				return err
			}
			method.Results = append(method.Results, identifiers...)
		}
	}

	i.Methods = append(i.Methods, method)
	return nil
}

func extractIdentifiers(field *ast.Field, imports *ImportSet) ([]*Identifier, error) {
	identifierType, err := unwrap(field.Type, imports)
	if err != nil {
		return nil, err
	}

	if len(field.Names) == 0 {
		return []*Identifier{
			&Identifier{
				Name:      identSymGen.Next(),
				valueType: identifierType,
			},
		}, nil
	}

	identifiers := make([]*Identifier, len(field.Names))
	for i, name := range field.Names {
		identifiers[i] = &Identifier{
			Name:      name.Name,
			valueType: identifierType,
		}
	}

	return identifiers, nil
}

func unwrap(node ast.Expr, imports *ImportSet) (t Type, err error) {
	switch nodeType := node.(type) {
	case *ast.Ellipsis:
		var subType Type
		subType, err = unwrap(nodeType.Elt, imports)
		if err != nil {
			return
		}
		t = &Ellipsis{
			subType: subType,
		}
	case *ast.ArrayType:
		var subType Type
		subType, err = unwrap(nodeType.Elt, imports)
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
	case *ast.ChanType:
		switch nodeType.Dir {
		case ast.SEND:
			var subType Type
			subType, err = unwrap(nodeType.Value, imports)
			if err != nil {
				return
			}
			t = &SendChannel{
				subType: subType,
			}
		case ast.RECV:
			var subType Type
			subType, err = unwrap(nodeType.Value, imports)
			if err != nil {
				return
			}
			t = &ReceiveChannel{
				subType: subType,
			}
		case ast.SEND + ast.RECV:
			var subType Type
			subType, err = unwrap(nodeType.Value, imports)
			if err != nil {
				return
			}
			t = &Channel{
				subType: subType,
			}
		}
	case *ast.StarExpr:
		var subType Type
		subType, err = unwrap(nodeType.X, imports)
		if err != nil {
			return
		}
		t = &Pointer{
			subType: subType,
		}
	case *ast.InterfaceType, *ast.StructType:
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
		err = fmt.Errorf("internal error: unsupported field type node: %#v", nodeType)
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

type Identifier struct {
	Name            string
	valueType       Type
	titleCase       string
	parameterFormat string
	referenceFormat string
	fieldFormat     string
	signature       string
}

func (i *Identifier) TitleCase() string {
	if i.titleCase == "" {
		i.titleCase = strings.Title(i.Name)
	}
	return i.titleCase
}

func (i *Identifier) ParameterFormat() string {
	if i.parameterFormat == "" {
		i.parameterFormat = fmt.Sprintf("%s %s", i.Name, i.valueType.ParameterFormat())
	}

	return i.parameterFormat
}

func (i *Identifier) ReferenceFormat() string {
	if i.referenceFormat == "" {
		i.referenceFormat = fmt.Sprintf("%s%s", i.Name, i.valueType.ReferenceFormat())
	}

	return i.referenceFormat
}

func (i *Identifier) FieldFormat() string {
	if i.fieldFormat == "" {
		i.fieldFormat = fmt.Sprintf("%s %s", i.TitleCase(), i.valueType.FieldFormat())
	}

	return i.fieldFormat
}

func (i *Identifier) Signature() string {
	if i.signature == "" {
		i.signature = i.valueType.ParameterFormat()
	}

	return i.signature
}

type Type interface {
	ParameterFormat() string
	ReferenceFormat() string
	FieldFormat() string
}

type Array struct {
	subType         Type
	scale           string
	parameterFormat string
	fieldFormat     string
}

func (t *Array) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("[%s]%s", t.scale, t.subType.ParameterFormat())

	return t.parameterFormat
}

func (t *Array) ReferenceFormat() string {
	return t.subType.ReferenceFormat()
}

func (t *Array) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("[%s]%s", t.scale, t.subType.FieldFormat())

	return t.fieldFormat
}

type Ellipsis struct {
	subType         Type
	parameterFormat string
	fieldFormat     string
}

func (t *Ellipsis) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("...%s", t.subType.ParameterFormat())

	return t.parameterFormat
}

func (t *Ellipsis) ReferenceFormat() string {
	return "..."
}

func (t *Ellipsis) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("[]%s", t.subType.FieldFormat())

	return t.fieldFormat
}

type Channel struct {
	subType         Type
	parameterFormat string
	fieldFormat     string
}

func (t *Channel) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("chan %s", t.subType.ParameterFormat())

	return t.parameterFormat
}

func (t *Channel) ReferenceFormat() string {
	return t.subType.ReferenceFormat()
}

func (t *Channel) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("chan %s", t.subType.FieldFormat())

	return t.fieldFormat
}

type ReceiveChannel struct {
	subType         Type
	parameterFormat string
	fieldFormat     string
}

func (t *ReceiveChannel) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("<-chan %s", t.subType.ParameterFormat())

	return t.parameterFormat
}

func (t *ReceiveChannel) ReferenceFormat() string {
	return t.subType.ReferenceFormat()
}

func (t *ReceiveChannel) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("<-chan %s", t.subType.FieldFormat())

	return t.fieldFormat
}

type SendChannel struct {
	subType         Type
	parameterFormat string
	fieldFormat     string
}

func (t *SendChannel) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("chan<- %s", t.subType.ParameterFormat())

	return t.parameterFormat
}

func (t *SendChannel) ReferenceFormat() string {
	return t.subType.ReferenceFormat()
}

func (t *SendChannel) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("chan<- %s", t.subType.FieldFormat())

	return t.fieldFormat
}

type Pointer struct {
	subType         Type
	parameterFormat string
	fieldFormat     string
}

func (t *Pointer) ParameterFormat() string {
	if t.parameterFormat != "" {
		return t.parameterFormat
	}

	t.parameterFormat = fmt.Sprintf("*%s", t.subType.ParameterFormat())

	return t.parameterFormat
}

func (t *Pointer) ReferenceFormat() string {
	return t.subType.ReferenceFormat()
}

func (t *Pointer) FieldFormat() string {
	if t.fieldFormat != "" {
		return t.fieldFormat
	}

	t.fieldFormat = fmt.Sprintf("*%s", t.subType.FieldFormat())

	return t.fieldFormat
}

type BasicType struct {
	Name            string
	Qualifier       string
	parameterFormat string
	fieldFormat     string
}

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

func (t *BasicType) ReferenceFormat() string {
	return ""
}

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
