package main

import "go/ast"

// SingleFileEntryVisitor 单文件入口 visitor
type SingleFileEntryVisitor struct {
	fv *FileVisitor
}

func (s *SingleFileEntryVisitor) Get() File {
	types := make([]Type, 0, len(s.fv.tvs))

	for _, tv := range s.fv.tvs {
		types = append(types, Type{
			Name:   tv.name,
			Fields: tv.fields,
		})
	}
	return File{
		Package: s.fv.pkg,
		Types:   types,
	}
}

func (s *SingleFileEntryVisitor) Visit(node ast.Node) ast.Visitor {
	fn, ok := node.(*ast.File)
	if !ok {
		return s
	}

	s.fv = &FileVisitor{
		pkg: fn.Name.String(),
	}

	return s.fv
}

type File struct {
	Package string
	Types   []Type
}

type FileVisitor struct {
	pkg string
	tvs []*TypeVisitor
}

func (f *FileVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		tv := &TypeVisitor{name: n.Name.String()}
		f.tvs = append(f.tvs, tv)
		return tv
	}

	return f
}

type Type struct {
	Name   string
	Fields []Field
}

type TypeVisitor struct {
	name   string
	fields []Field
}

func (t *TypeVisitor) Visit(node ast.Node) ast.Visitor {
	fn, ok := node.(*ast.Field)
	if !ok {
		return t
	}

	for _, name := range fn.Names {
		t.fields = append(t.fields, Field{Name: name.String()})
	}

	return t
}

type Field struct {
	Name string
	Type string
}
