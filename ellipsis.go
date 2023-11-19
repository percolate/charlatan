package main

import "fmt"

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
