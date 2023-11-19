package main

import "fmt"

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
