package simplifier_test

import (
	"testing"
	"text/template"

	"github.com/mh-cbon/template-tree-simplifier/simplifier"
)

func TestPrintsAnything(t *testing.T) {
	//-
	defFuncs := template.FuncMap{
		"getSlice": func() []string { return []string{} },
	}
	// add more test with $ and numbers and bool
	testTable := []TestData{
		TestData{
			tplstr:         `{{.}}`,
			funcs:          defFuncs,
			printsanything: true,
			expectPrints:   true,
		},
		TestData{
			tplstr:         `{{$y := "r"}}{{$y}}{{$y := "r"}}{{$y}}`,
			funcs:          defFuncs,
			printsanything: true,
			expectPrints:   true,
		},
		TestData{
			tplstr:         `{{range .}}{{.}}{{end}}`,
			funcs:          defFuncs,
			data:           []string{},
			printsanything: true,
			expectPrints:   true,
		},
		TestData{
			tplstr:         `{{$y := "hello"}}`,
			funcs:          defFuncs,
			printsanything: true,
			expectPrints:   false,
		},
		TestData{
			tplstr:         `{{range .}}{{end}}`,
			funcs:          defFuncs,
			data:           []string{},
			printsanything: true,
			expectPrints:   false,
		},
		TestData{
			tplstr:         `TestNode`,
			funcs:          defFuncs,
			data:           []string{},
			printsanything: true,
			expectPrints:   true,
		},
	}

	for i, testData := range testTable {
		if execTestData(testData, t, i) == false {
			break
		}
	}
}

func printsanythingtemplate(t *template.Template) bool {
	ret, err := t.Clone()
	if err != nil {
		panic(err)
	}
	for _, t := range ret.Templates() {
		if t.Tree != nil {
			return simplifier.PrintsAnything(t.Tree)
		}
	}
	return false
}
