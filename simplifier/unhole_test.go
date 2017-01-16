package simplifier_test

import (
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/mh-cbon/template-tree-simplifier/funcmap"
	"github.com/mh-cbon/template-tree-simplifier/simplifier"
)

func TestUnhole(t *testing.T) {
	//-
	defFuncs := template.FuncMap{
		"split":              strings.Split,
		"join":               func(sep string, a []string) string { return strings.Join(a, sep) },
		"up":                 strings.ToUpper,
		"lower":              strings.ToLower,
		"incr":               func(s int) int { return s + 1 },
		"mul":                func(s int, d int) int { return s * d },
		"intf":               func() interface{} { return nil },
		"browsePropertyPath": funcmap.BrowsePropertyPath,
	}
	// create a reflect of interface{}
	var x []interface{}
	reflectInterface := reflect.TypeOf(x).Elem()
	// add more test with $ and numbers and bool
	testTable := []TestData{
		TestData{
			tplstr:       `{{range .Some}}{{end}}`,
			expectTplStr: `{{$var0 := .Some}}{{range $var0}}{{end}}`,
			funcs:        defFuncs,
			unhole:       true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some}}{{$y := .Some.Some}}{{$z := .}}`,
			expectTplStr: `{{$tpl_x := .Some}}{{$tpl_y := browsePropertyPath . "Some.Some"}}{{$tpl_z := .}}`,
			funcs:        defFuncs,
			unhole:       true,
			data:         type4{Some: type2{Some: ""}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type4{}),
					"$tpl_x": reflectInterface,
					"$tpl_y": reflectInterface,
					"$tpl_z": reflect.TypeOf(type4{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some}}{{$y := .Some.Some}}{{$z := .}}`,
			expectTplStr: `{{$tpl_x := .Some}}{{$tpl_y := browsePropertyPath . "Some.Some"}}{{$tpl_z := .}}`,
			funcs:        defFuncs,
			unhole:       true,
			data:         type4{Some: type2{Some: ""}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type4{}),
					"$tpl_x": reflectInterface,
					"$tpl_y": reflectInterface,
					"$tpl_z": reflect.TypeOf(type4{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some}}{{$y := $x.Some}}`,
			expectTplStr: `{{$tpl_x := .Some}}{{$tpl_y := browsePropertyPath $tpl_x "Some"}}`,
			funcs:        defFuncs,
			unhole:       true,
			data:         type4{Some: type2{Some: ""}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type4{}),
					"$tpl_x": reflectInterface,
					"$tpl_y": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Method}}`,
			expectTplStr: `{{$tpl_x := .Method}}`,
			funcs:        defFuncs,
			unhole:       true,
			data:         type4{Some: type2{Some: ""}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type4{}),
					"$tpl_x": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some.Method}}`,
			expectTplStr: `{{$tpl_x := browsePropertyPath . "Some.Method"}}`,
			funcs:        defFuncs,
			unhole:       true,
			data:         type4{Some: type4{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type4{}),
					"$tpl_x": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some.MethodArgs "r"}}`,
			expectTplStr: `{{$tpl_x := browsePropertyPath . "Some.MethodArgs" "r"}}`,
			funcs:        defFuncs,
			unhole:       true,
			data:         type4{Some: type4{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type4{}),
					"$tpl_x": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{$x := intf}}`,
			expectTplStr: `{{$tpl_x := intf}}`,
			funcs:        defFuncs,
			unhole:       true,
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(nil),
					"$tpl_x": reflectInterface,
				},
			},
		},
	}

	for _, testData := range testTable {
		if execTestData(testData, t) == false {
			break
		}
	}
}

func unholetemplate(t *template.Template, testData TestData) (*template.Template, *simplifier.State) {
	ret, err := t.Clone()
	if err != nil {
		panic(err)
	}
	var typeCheck *simplifier.State
	for _, t := range ret.Templates() {
		if t.Tree != nil {
			simplifier.Simplify(t.Tree)
			typeCheck = simplifier.TypeCheck(t.Tree, testData.data, testData.funcs)
			simplifier.Unhole(t.Tree, typeCheck, testData.data, testData.funcs)
		}
	}
	return ret, typeCheck
}
