package main

import "fmt"

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
