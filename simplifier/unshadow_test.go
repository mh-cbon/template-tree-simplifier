package simplifier_test

import (
	"strings"
	"testing"
	"text/template"

	"github.com/mh-cbon/template-tree-simplifier/simplifier"
)

func TestUnshadow(t *testing.T) {
	//-
	defFuncs := template.FuncMap{
		"split": strings.Split,
		"join":  func(sep string, a []string) string { return strings.Join(a, sep) },
		"up":    strings.ToUpper,
		"lower": strings.ToLower,
		"incr":  func(s int) int { return s + 1 },
		"mul":   func(s int, d int) int { return s * d },
	}
	// ad more test with $ and numbers and bool
	testTable := []TestData{
		TestData{
			tplstr:       `{{$y := "r"}}{{$y}}{{$y := "r"}}{{$y}}`,
			expectTplStr: `{{$y := "r"}}{{$y}}{{$y_shadow := "r"}}{{$y_shadow}}`,
			funcs:        defFuncs,
			unshadow:     true,
		},
		TestData{
			tplstr:       `{{$y := "r"}}{{$y}}{{$y := .}}{{$y}}{{$y.X}}`,
			expectTplStr: `{{$y := "r"}}{{$y}}{{$y_shadow := .}}{{$y_shadow}}{{$y_shadow.X}}`,
			funcs:        defFuncs,
			data:         struct{ X string }{X: ""},
			unshadow:     true,
		},
		TestData{
			tplstr:       `{{$y := "r"}}{{$y := "r"}}{{up $y}}`,
			expectTplStr: `{{$y := "r"}}{{$y_shadow := "r"}}{{up $y_shadow}}`,
			funcs:        defFuncs,
			unshadow:     true,
		},
	}

	for _, testData := range testTable {
		if execTestData(testData, t) == false {
			break
		}
	}
}

func unshadowtemplate(t *template.Template) *template.Template {
	ret, err := t.Clone()
	if err != nil {
		panic(err)
	}
	for _, t := range ret.Templates() {
		if t.Tree != nil {
			simplifier.Unshadow(t.Tree)
		}
	}
	return ret
}
