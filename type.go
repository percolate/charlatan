package main

import (
	"fmt"
	"go/types"
	"strconv"
)

// Type is the interface for all built-in types
type Type interface {
	// ParameterFormat returns syntax for a parameter declaration
	ParameterFormat() string
	// ReferenceFormat returns the syntax for a reference
	ReferenceFormat() string
	// FieldFormat returns the syntax for a field declaration
	FieldFormat() string
}

// BasicType represents all built-in simple types
type BasicType struct {
	Name            string
	Qualifier       string
	parameterFormat string
	fieldFormat     string
}

// ParameterFormat returns syntax for a parameter declaration
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

// ReferenceFormat returns the syntax for a reference
func (t *BasicType) ReferenceFormat() string {
	return ""
}

// FieldFormat returns the syntax for a field declaration
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

func unwrapType(t types.Type, imports *ImportSet) (r Type, err error) {
	switch actual := t.(type) {
	case *types.Array:
		var subType Type
		subType, err = unwrapType(actual.Elem(), imports)
		if err != nil {
			return
		}
		a := &Array{subType: subType}
		if actual.Len() != 0 {
			a.scale = strconv.FormatInt(actual.Len(), 10)
		}
		r = a
	case *types.Slice:
		var subType Type
		subType, err = unwrapType(actual.Elem(), imports)
		if err != nil {
			return
		}
		r = &Array{subType: subType}
	case *types.Map:
		var keyType Type
		keyType, err = unwrapType(actual.Key(), imports)
		if err != nil {
			return
		}
		var subType Type
		subType, err = unwrapType(actual.Elem(), imports)
		if err != nil {
			return
		}
		r = &Map{keyType: keyType, subType: subType}
	case *types.Chan:
		var subType Type
		subType, err = unwrapType(actual.Elem(), imports)
		if err != nil {
			return
		}
		switch actual.Dir() {
		case types.SendOnly:
			r = &SendChannel{subType: subType}
		case types.RecvOnly:
			r = &ReceiveChannel{subType: subType}
		case types.SendRecv:
			r = &Channel{subType: subType}
		}
	case *types.Pointer:
		var subType Type
		subType, err = unwrapType(actual.Elem(), imports)
		if err != nil {
			return nil, err
		}
		r = &Pointer{subType: subType}
	case *types.Interface, *types.Struct, *types.Signature:
		r = &BasicType{Name: actual.String()}
	case *types.Named:
		b := &BasicType{Name: actual.Obj().Name()}
		if actual.Obj().Pkg() != nil {
			b.Qualifier = actual.Obj().Pkg().Name()
			imports.RequireByName(b.Qualifier)
		}
		r = b
	case *types.Basic:
		r = &BasicType{Name: actual.Name()}
	default:
		err = fmt.Errorf("internal error: unsupported parameter type for type: %#v", actual)
	}

	return
}
