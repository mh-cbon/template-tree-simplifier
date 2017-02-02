package simplifier_test

import (
	"testing"
	"text/template"

	"github.com/mh-cbon/template-tree-simplifier/simplifier"
)

func TestUseDot(t *testing.T) {
	//-
	defFuncs := template.FuncMap{
		"getSlice": func() []string { return []string{} },
	}
	// add more test with $ and numbers and bool
	testTable := []TestData{
		TestData{
			tplstr:       `{{$y := "r"}}{{$y}}{{$y := "r"}}{{$y}}`,
			funcs:        defFuncs,
			usedot:       true,
			expectDotUse: false,
		},
		TestData{
			tplstr:       `{{$y := .}}{{$y}}`,
			funcs:        defFuncs,
			usedot:       true,
			expectDotUse: true,
		},
		TestData{
			tplstr:       `{{.}}`,
			funcs:        defFuncs,
			usedot:       true,
			expectDotUse: true,
		},
		TestData{
			tplstr:       `{{$y := getSlice}}{{range $y}}{{.}}{{end}}`,
			funcs:        defFuncs,
			usedot:       true,
			expectDotUse: false,
		},
		TestData{
			tplstr:       `{{$y := "rr"}}{{with $y}}{{.}}{{end}}`,
			funcs:        defFuncs,
			usedot:       true,
			expectDotUse: false,
		},
		TestData{
			tplstr:       `{{range .}}{{.}}{{end}}`,
			funcs:        defFuncs,
			data:         []string{},
			usedot:       true,
			expectDotUse: true,
		},
		TestData{
			tplstr:       `{{with .}}{{.}}{{end}}`,
			funcs:        defFuncs,
			data:         []string{},
			usedot:       true,
			expectDotUse: true,
		},
		TestData{
			tplstr:       `{{if true}}{{.}}{{end}}`,
			funcs:        defFuncs,
			data:         []string{},
			usedot:       true,
			expectDotUse: true,
		},
		TestData{
			tplstr:       `{{if true}}{{else}}{{.}}{{end}}`,
			funcs:        defFuncs,
			data:         []string{},
			usedot:       true,
			expectDotUse: true,
		},
	}

	for i, testData := range testTable {
		if execTestData(testData, t, i) == false {
			break
		}
	}
}

func usedottemplate(t *template.Template) bool {
	ret, err := t.Clone()
	if err != nil {
		panic(err)
	}
	for _, t := range ret.Templates() {
		if t.Tree != nil {
			return simplifier.IsUsingDot(t.Tree)
		}
	}
	return false
}
