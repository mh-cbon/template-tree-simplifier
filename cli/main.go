package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/mh-cbon/print-template-tree/printer"
	"github.com/mh-cbon/template-tree-simplifier/simplifier"
)

type tplData struct {
	Some string
}

func (t tplData) Method(s string) string {
	return s
}

func main() {
	file := "cli/tpl/test.tpl"
	fmt.Println(file)

	funcs := template.FuncMap{
		"split":     strings.Split,
		"join":      func(sep string, a []string) string { return strings.Join(a, sep) },
		"up":        strings.ToUpper,
		"lower":     strings.ToLower,
		"fnWithErr": func(a string) (string, error) { return a, nil },
		"incr":      func(s int) int { return s + 1 },
		"mul":       func(s int, d int) int { return s * d },
	}

	// data := "what"
	data := tplData{}

	tpl, err := template.New("").Funcs(funcs).ParseFiles(file)
	if err != nil {
		panic(err)
	}

	for _, t := range tpl.Templates() {
		if t.Tree != nil {
			fmt.Printf("------------------\n")
			fmt.Printf("Tree.Name=%v\n\n", t.Name())
			fmt.Println("BEFORE TRANSFORM")
			printTemplateInfo(t, data)
			simplifier.Unshadow(t.Tree)
			simplifier.Simplify(t.Tree)
			fmt.Println("AFTER TRANSFORM")
			printTemplateInfo(t, data)
		}
	}
}

func printTemplateInfo(t *template.Template, data interface{}) {
	fmt.Println("TEMPLATE CONTENT")
	printer.PrintContent(t.Tree)
	fmt.Println("TEMPLATE EXECUTION")
	printTemplateExec(t, data)
	fmt.Println("TEMPLATE TREE")
	printer.Print(t.Tree)
}

func printTemplateExec(t *template.Template, data interface{}) {
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%q\n", b.String())
}
