package main

import (
	"fmt"
	"go/ast"
	"strings"
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
	for i, parameter := range functionType.Params.List {
		values := extractValues(parameter, i, "arg", imports)
		method.Parameters = append(method.Parameters, values...)
	}

	if functionType.Results != nil {
		for i, result := range functionType.Results.List {
			values := extractValues(result, i, "ret", imports)
			method.Results = append(method.Results, values...)
		}
	}

	i.Methods = append(i.Methods, method)
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

// Method represents a method in an interface's method set
type Method struct {
	Interface      string
	Name           string
	Parameters     []*Value
	Results        []*Value
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
		for i, value := range m.Parameters {
			params[i] = value.functionDeclarationFormat()
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
		for i, value := range m.Parameters {
			params[i] = value.argumentFormat()
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
		for i, value := range m.Results {
			params[i] = value.functionDeclarationFormat()
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
		for i, value := range m.Results {
			params[i] = value.argumentFormat()
		}
		m.resultsCall = strings.Join(params, ", ")
	}

	return m.resultsCall
}

// Value represents a Parameter or Result of a Method
type Value struct {
	Name       string // name if a named parameter/result, else null string
	Type       string
	Qualifier  string
	Pointer    bool
	Elliptical bool
	ChanDir    int
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
