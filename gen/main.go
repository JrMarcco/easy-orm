package main

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed tpl.gohtml
var tpl string

type OrmFile struct {
	File
	Ops []string
}

func main() {
	src := os.Args[1]
	dstDir := filepath.Dir(src)
	filename := filepath.Base(src)

	index := strings.LastIndex(filename, ".")
	dst := filepath.Join(dstDir, filename[:index]+".gen.go")

	f, err := os.Create(dst)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = f.Close() }()

	if err = gen(f, src); err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("gofmt", "-w", dst)
	if err = cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success to generate: ", dst)
}

func gen(w io.Writer, srcFile string) error {
	var err error

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, srcFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	fve := &FileVisitEntry{}
	ast.Walk(fve, f)
	file := fve.Get()

	t := template.New("easyorm_gen")
	if t, err = t.Parse(tpl); err != nil {
		return err
	}

	return t.Execute(w, OrmFile{
		File: file,
		Ops:  []string{"Eq", "Ne", "Gt", "Ge", "Lt", "Le"},
	})
}
