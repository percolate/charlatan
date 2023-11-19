package main

import "fmt"

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
