package main

import (
	"fmt"
)

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
