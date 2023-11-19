package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
)

func unwrapExpr(node ast.Expr, imports *ImportSet) (t Type, err error) {
	switch nodeType := node.(type) {
	case *ast.Ellipsis:
		var subType Type
		subType, err = unwrapExpr(nodeType.Elt, imports)
		if err != nil {
			return
		}
		t = &Ellipsis{
			subType: subType,
		}
	case *ast.ArrayType:
		var subType Type
		subType, err = unwrapExpr(nodeType.Elt, imports)
		if err != nil {
			return
		}
		a := &Array{
			subType: subType,
		}
		if nodeType.Len != nil {
			if lit, ok := nodeType.Len.(*ast.BasicLit); ok {
				a.scale = lit.Value
			} else {
				err = fmt.Errorf("internal error: unsupported array len type node: %#v", nodeType.Len)
				return
			}
		}
		t = a
	case *ast.MapType:
		var keyType Type
		keyType, err = unwrapExpr(nodeType.Key, imports)
		if err != nil {
			return
		}
		var subType Type
		subType, err = unwrapExpr(nodeType.Value, imports)
		if err != nil {
			return
		}
		t = &Map{keyType: keyType, subType: subType}
	case *ast.ChanType:
		var subType Type
		subType, err = unwrapExpr(nodeType.Value, imports)
		if err != nil {
			return
		}
		switch nodeType.Dir {
		case ast.SEND:
			t = &SendChannel{
				subType: subType,
			}
		case ast.RECV:
			t = &ReceiveChannel{
				subType: subType,
			}
		case ast.SEND + ast.RECV:
			t = &Channel{
				subType: subType,
			}
		}
	case *ast.StarExpr:
		var subType Type
		subType, err = unwrapExpr(nodeType.X, imports)
		if err != nil {
			return
		}
		t = &Pointer{
			subType: subType,
		}
	case *ast.InterfaceType, *ast.StructType, *ast.FuncType:
		var buf bytes.Buffer
		if err = format.Node(&buf, token.NewFileSet(), nodeType); err != nil {
			return
		}
		t = &BasicType{
			Name: buf.String(),
		}
	case *ast.SelectorExpr:
		selector := nodeType.X.(*ast.Ident).Name
		imports.RequireByName(selector)
		t = &BasicType{
			Qualifier: selector,
			Name:      nodeType.Sel.Name,
		}
	case *ast.Ident:
		t = &BasicType{
			Name: nodeType.Name,
		}
	case *ast.IndexExpr:
		// ignore this expression kind
	default:
		err = fmt.Errorf("internal error: unsupported parameter type for expr: %#v", nodeType)
	}

	return
}
