package main

import (
	_ "embed"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"text/template"
)

//go:embed tpl.gohtml
var ormTemplate string

func gen(w io.Writer, srcFile string) error {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, srcFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	sv := &SingleFileEntryVisitor{}
	ast.Walk(sv, f)

	tpl := template.New("orm-gen")
	parse, err := tpl.Parse(ormTemplate)
	if err != nil {
		return err
	}

	return parse.Execute(w, sv.Get())
}
