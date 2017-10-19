package main

import (
	"fmt"
	"go/ast"
	"strings"
)

var (
	identSymGen = SymbolGenerator{Prefix: "ident"}
)

// Import represents a declared import
type Import struct {
	Name     string // the package's name
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

func (r *ImportSet) Remove(ri *Import) {
	for index, i := range r.imports {
		if i.Name == ri.Name && i.Path == ri.Path {
			r.imports = append(r.imports[:index], r.imports[index+1:]...)
			return
		}
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

func (r *ImportSet) GetAll() []*Import {
	return r.imports
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
		if imp.Name == s {
			r.imports[i].Required = true
		}
	}
}

// Interface represents a declared interface.
type Interface struct {
	Name    string
	Methods []*Method
}

func (i *Interface) addMethod(field *ast.Field, imports *ImportSet) {
	functionType, ok := field.Type.(*ast.FuncType)
	if !ok {
		return
	}

	method := &Method{
		Interface: i.Name,
		Name:      field.Names[0].Name,
	}

	// `Params.List` can be 0-length, but `Results` can be nil
	for _, parameter := range functionType.Params.List {
		identifiers := extractIdentifiers(parameter, imports)
		method.Parameters = append(method.Parameters, identifiers...)
	}

	if functionType.Results != nil {
		for _, result := range functionType.Results.List {
			identifiers := extractIdentifiers(result, imports)
			method.Results = append(method.Results, identifiers...)
		}
	}

	i.Methods = append(i.Methods, method)
}

func extractIdentifiers(field *ast.Field, imports *ImportSet) []*Identifier {
	identifierType := unwrap(field.Type, imports)

	if len(field.Names) == 0 {
		return []*Identifier{
			&Identifier{
				Name:      identSymGen.Next(),
				valueType: identifierType,
			},
		}
	}

	identifiers := make([]*Identifier, len(field.Names))
	for i, name := range field.Names {
		identifiers[i] = &Identifier{
			Name:      name.Name,
			valueType: identifierType,
		}
	}

	return identifiers
}

func unwrap(node ast.Expr, imports *ImportSet) Type {
	switch nodeType := node.(type) {
	case *ast.Ellipsis:
		return &Ellipsis{
			subType: unwrap(nodeType.Elt, imports),
		}
	case *ast.ChanType:
		switch nodeType.Dir {
		case ast.SEND:
			return &SendChannel{
				subType: unwrap(nodeType.Value, imports),
			}
		case ast.RECV:
			return &ReceiveChannel{
				subType: unwrap(nodeType.Value, imports),
			}
		case ast.SEND + ast.RECV:
			return &Channel{
				subType: unwrap(nodeType.Value, imports),
			}
		}
	case *ast.StarExpr:
		return &Pointer{
			subType: unwrap(nodeType.X, imports),
		}
	case *ast.InterfaceType:
		return &BasicType{
			Name: "interface{}",
		}
	case *ast.StructType:
		return &BasicType{
			Name: "struct{}",
		}
	case *ast.SelectorExpr:
		selector := nodeType.X.(*ast.Ident).Name
		imports.RequireByName(selector)
		return &BasicType{
			Qualifier: selector,
			Name:      nodeType.Sel.Name,
		}
	case *ast.Ident:
		return &BasicType{
			Name: nodeType.Name,
		}
	}

	return nil
}

// Method represents a method in an interface's method set
type Method struct {
	Interface      string
	Name           string
	Parameters     []*Identifier
	Results        []*Identifier
	parametersDecl string
	parametersCall string
	resultsDecl    string
	resultsCall    string
}

func (m *Method) ParametersDeclaration() string {
	if len(m.Parameters) == 0 {
		return ""
	}
	if m.parametersDecl == "" {
		params := make([]string, len(m.Parameters))
		for i, param := range m.Parameters {
			params[i] = param.ParameterFormat()
		}
		m.parametersDecl = strings.Join(params, ", ")
	}

	return m.parametersDecl
}

func (m *Method) ParametersReference() string {
	if len(m.Parameters) == 0 {
		return ""
	}
	if m.parametersCall == "" {
		params := make([]string, len(m.Parameters))
		for i, param := range m.Parameters {
			params[i] = param.ReferenceFormat()
		}
		m.parametersCall = strings.Join(params, ", ")
	}

	return m.parametersCall
}

func (m *Method) ResultsDeclaration() string {
	if len(m.Results) == 0 {
		return ""
	}
	if m.resultsDecl == "" {
		params := make([]string, len(m.Results))
		for i, param := range m.Results {
			params[i] = param.ParameterFormat()
		}
		m.resultsDecl = strings.Join(params, ", ")
	}

	return m.resultsDecl
}

func (m *Method) ResultsReference() string {
	if len(m.Results) == 0 {
		return ""
	}
	if m.resultsCall == "" {
		params := make([]string, len(m.Results))
		for i, param := range m.Results {
			params[i] = param.ReferenceFormat()
		}
		m.resultsCall = strings.Join(params, ", ")
	}

	return m.resultsCall
}

type Identifier struct {
	Name            string
	valueType       Type
	titleCase       string
	parameterFormat string
	referenceFormat string
	fieldFormat     string
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

type Type interface {
	ParameterFormat() string
	ReferenceFormat() string
	FieldFormat() string
}

type Ellipsis struct {
	subType Type
}

func (t *Ellipsis) ParameterFormat() string {
	return fmt.Sprintf("%s...", t.subType.ParameterFormat())
}

func (t *Ellipsis) ReferenceFormat() string {
	return "..."
}

func (t *Ellipsis) FieldFormat() string {
	return fmt.Sprintf("%s[]", t.subType.FieldFormat())
}

type Channel struct {
	subType Type
}

func (t *Channel) ParameterFormat() string {
	return fmt.Sprintf("chan %s", t.subType.ParameterFormat())
}

func (t *Channel) ReferenceFormat() string {
	return t.subType.ReferenceFormat()
}

func (t *Channel) FieldFormat() string {
	return fmt.Sprintf("chan %s", t.subType.FieldFormat())
}

type ReceiveChannel struct {
	subType Type
}

func (t *ReceiveChannel) ParameterFormat() string {
	return fmt.Sprintf("<-chan %s", t.subType.ParameterFormat())
}

func (t *ReceiveChannel) ReferenceFormat() string {
	return t.subType.ReferenceFormat()
}

func (t *ReceiveChannel) FieldFormat() string {
	return fmt.Sprintf("<-chan %s", t.subType.FieldFormat())
}

type SendChannel struct {
	subType Type
}

func (t *SendChannel) ParameterFormat() string {
	return fmt.Sprintf("chan<- %s", t.subType.ParameterFormat())
}

func (t *SendChannel) ReferenceFormat() string {
	return t.subType.ReferenceFormat()
}

func (t *SendChannel) FieldFormat() string {
	return fmt.Sprintf("chan<- %s", t.subType.FieldFormat())
}

type Pointer struct {
	subType Type
}

func (t *Pointer) ParameterFormat() string {
	return fmt.Sprintf("*%s", t.subType.ParameterFormat())
}

func (t *Pointer) ReferenceFormat() string {
	return t.subType.ReferenceFormat()
}

func (t *Pointer) FieldFormat() string {
	return fmt.Sprintf("*%s", t.subType.FieldFormat())
}

type BasicType struct {
	Name      string
	Qualifier string
}

func (t *BasicType) ParameterFormat() string {
	if t.Qualifier != "" {
		return fmt.Sprintf("%s.%s", t.Qualifier, t.Name)
	}

	return t.Name
}

func (t *BasicType) ReferenceFormat() string {
	return ""
}

func (t *BasicType) FieldFormat() string {
	if t.Qualifier != "" {
		return fmt.Sprintf("%s.%s", t.Qualifier, t.Name)
	}

	return t.Name
}
