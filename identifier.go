package main

import (
	"fmt"
	"strings"
)

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
