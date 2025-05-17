package main

import (
	"fmt"
	"go/ast"
)

type FileVisitEntry struct {
	fv *fileVisitor
}

func (fve *FileVisitEntry) Get() File {
	if fve.fv != nil {
		return fve.fv.Get()
	}
	return File{}
}

func (fve *FileVisitEntry) Visit(node ast.Node) ast.Visitor {
	if n, ok := node.(*ast.File); ok {
		fve.fv = &fileVisitor{
			pkg: n.Name.String(),
		}
		return fve.fv
	}
	return fve
}

type File struct {
	Pkg     string
	Types   []Type
	Imports []string
}

type fileVisitor struct {
	pkg     string
	types   []*typeVisitor
	imports []string
}

func (fv *fileVisitor) Get() File {
	types := make([]Type, 0, len(fv.types))
	for _, typ := range fv.types {
		types = append(types, typ.Get())
	}
	return File{
		Pkg:     fv.pkg,
		Types:   types,
		Imports: fv.imports,
	}
}

func (fv *fileVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		tv := &typeVisitor{
			name:   n.Name.String(),
			fields: make([]Field, 0),
		}
		fv.types = append(fv.types, tv)
		return tv
	case *ast.ImportSpec:
		fv.imports = append(fv.imports, fmt.Sprintf("%s %s", n.Name.String(), n.Path.Value))
	}

	return fv
}

type Type struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}

type typeVisitor struct {
	name   string
	fields []Field
}

func (tv *typeVisitor) Get() Type {
	return Type{
		Name:   tv.name,
		Fields: tv.fields,
	}
}

func (tv *typeVisitor) Visit(node ast.Node) ast.Visitor {
	if fn, ok := node.(*ast.Field); ok {
		var typeName string
		switch n := fn.Type.(type) {
		case *ast.Ident:
			typeName = n.String()
		case *ast.StarExpr:
			switch xType := n.X.(type) {
			case *ast.Ident:
				typeName = fmt.Sprintf("*%s", xType.String())
			case *ast.SelectorExpr:
				typeName = fmt.Sprintf("*%s.%s", xType.X.(*ast.Ident).String(), xType.Sel.String())
			}
		default:
			fmt.Println(n)
		}

		tv.fields = append(tv.fields, Field{
			Name: fn.Names[0].String(),
			Type: typeName,
		})
		return nil
	}

	return tv
}
