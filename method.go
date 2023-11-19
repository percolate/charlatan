package main

import "strings"

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
