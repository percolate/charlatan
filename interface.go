package main

import (
	"fmt"
	"go/ast"
	"go/types"
)

var (
	identSymGen = symbolGenerator{Prefix: "ident"}
)

// Interface represents a declared interface.
type Interface struct {
	Name    string
	Methods []*Method
	embeds  []string
}

func (i *Interface) addMethodFromField(field *ast.Field, imports *ImportSet) error {
	functionType, ok := field.Type.(*ast.FuncType)
	if !ok {
		return fmt.Errorf("internal error: expected *ast.FuncType, have: %#v", field)
	}

	method := &Method{
		Interface: i.Name,
		Name:      field.Names[0].Name,
	}

	identSymGen.reset()
	// `Params.List` can be 0-length, but `Results` can be nil
	for _, parameter := range functionType.Params.List {
		identifiers, err := extractIdentifiersFromField(parameter, imports)
		if err != nil {
			return err
		}
		method.Parameters = append(method.Parameters, identifiers...)
	}

	if functionType.Results != nil {
		for _, result := range functionType.Results.List {
			identifiers, err := extractIdentifiersFromField(result, imports)
			if err != nil {
				return err
			}
			method.Results = append(method.Results, identifiers...)
		}
	}

	i.Methods = append(i.Methods, method)
	return nil
}

func (i *Interface) addMethodFromType(f *types.Func, imports *ImportSet) error {
	method := &Method{
		Interface: i.Name,
		Name:      f.Name(),
	}

	sig := f.Type().(*types.Signature)
	parameters, err := extractIdentifiersFromTuple(sig.Params(), imports)
	if err != nil {
		return err
	}
	if sig.Variadic() {
		last := parameters[len(parameters)-1]
		if avt, ok := last.ValueType.(*Array); ok {
			last.ValueType = &Ellipsis{subType: avt.subType}
		} else {
			last.ValueType = &Ellipsis{subType: last.ValueType}
		}
	}
	method.Parameters = append(method.Parameters, parameters...)

	results, err := extractIdentifiersFromTuple(sig.Results(), imports)
	if err != nil {
		return err
	}
	method.Results = append(method.Results, results...)

	i.Methods = append(i.Methods, method)

	return nil
}

func extractIdentifiersFromField(field *ast.Field, imports *ImportSet) ([]*Identifier, error) {
	identifierType, err := unwrapExpr(field.Type, imports)
	if err != nil {
		return nil, err
	}

	if len(field.Names) == 0 {
		return []*Identifier{
			{
				Name:      identSymGen.next(),
				ValueType: identifierType,
			},
		}, nil
	}

	identifiers := make([]*Identifier, len(field.Names))
	for i, name := range field.Names {
		identifiers[i] = &Identifier{
			Name:      name.Name,
			ValueType: identifierType,
		}
	}

	return identifiers, nil
}

func extractIdentifiersFromTuple(tuple *types.Tuple, imports *ImportSet) ([]*Identifier, error) {
	if 0 == tuple.Len() {
		return nil, nil
	}

	idents := make([]*Identifier, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		p := tuple.At(i)
		identifierType, err := unwrapType(p.Type(), imports)
		if err != nil {
			return nil, err
		}
		ident := &Identifier{
			Name:      p.Name(),
			ValueType: identifierType,
		}
		if "" == ident.Name {
			ident.Name = identSymGen.next()
		}
		idents[i] = ident
	}

	return idents, nil
}
