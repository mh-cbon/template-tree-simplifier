package simplifier_test

import (
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/mh-cbon/template-tree-simplifier/simplifier"
)

type type1 struct {
	Some []string
}

type type2 struct {
	Some string
}

type type3 struct {
	Some type2
}

type type4 struct {
	Some interface{}
}

type type5 struct {
	Some *type3
}

type type6 struct {
	Some         []type1
	SomeSlicePtr []*type1
}

type type7 struct {
	Some []*type7
}

func (t type4) Method() interface{}        { return nil }
func (t type4) MethodArgs(a string) string { return "" }

func TestTypeCheck(t *testing.T) {
	//-
	defFuncs := template.FuncMap{
		"split": strings.Split,
		"join":  func(sep string, a []string) string { return strings.Join(a, sep) },
		"up":    strings.ToUpper,
		"lower": strings.ToLower,
		"incr":  func(s int) int { return s + 1 },
		"mul":   func(s int, d int) int { return s * d },
		"intf":  func() interface{} { return nil },
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
			typecheck:    true,
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
			tplstr:       `{{range .}}{{end}}`,
			expectTplStr: `{{$var0 := .}}{{range $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         []string{},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf([]string{}),
					"$var0": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{range $x := .Some}}{{end}}`,
			expectTplStr: `{{$var0 := .Some}}{{range $tpl_x := $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf(""),
					"$tpl_x": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .}}{{range $x.Some}}{{end}}`,
			expectTplStr: `{{$tpl_x := .}}{{$var0 := $tpl_x.Some}}{{range $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type1{}),
					"$tpl_x": reflect.TypeOf(type1{}),
					"$var0":  reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{range $x, $y := .Some}}{{end}}`,
			expectTplStr: `{{$var0 := .Some}}{{range $tpl_x, $tpl_y := $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf(""),
					"$tpl_x": reflect.TypeOf(1),
					"$tpl_y": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{with $x := .Some}}{{end}}`,
			expectTplStr: `{{$var0 := .Some}}{{with $tpl_x := $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf([]string{}),
					"$tpl_x": reflect.TypeOf([]string{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .}}{{with $x.Some}}{{end}}`,
			expectTplStr: `{{$tpl_x := .}}{{$var0 := $tpl_x.Some}}{{with $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type1{}),
					"$tpl_x": reflect.TypeOf(type1{}),
					"$var0":  reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".": reflect.TypeOf([]string{}),
				},
			},
		},
		TestData{
			tplstr:       `{{with .}}{{$y := .}}{{end}}`,
			expectTplStr: `{{$var0 := .}}{{with $var0}}{{$tpl_y := .}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf(type1{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf(type1{}),
					"$tpl_y": reflect.TypeOf(type1{}),
				},
			},
		},
		TestData{
			tplstr:       `{{with .}}{{$y := $}}{{end}}`,
			expectTplStr: `{{$var0 := .}}{{with $var0}}{{$tpl_y := $}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf(type1{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf(type1{}),
					"$tpl_y": reflect.TypeOf(type1{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := up .Some}}`,
			expectTplStr: `{{$tpl_x := up .Some}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type2{Some: ""},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type2{}),
					"$tpl_x": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := up .Some}}{{$y := $x}}`,
			expectTplStr: `{{$tpl_x := up .Some}}{{$tpl_y := $tpl_x}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type2{Some: ""},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type2{}),
					"$tpl_x": reflect.TypeOf(""),
					"$tpl_y": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := ""}}`,
			expectTplStr: `{{$tpl_x := ""}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type2{Some: ""},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type2{}),
					"$tpl_x": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := 2}}`,
			expectTplStr: `{{$tpl_x := 2}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type2{Some: ""},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type2{}),
					"$tpl_x": reflect.TypeOf(2),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := true}}`,
			expectTplStr: `{{$tpl_x := true}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type2{Some: ""},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type2{}),
					"$tpl_x": reflect.TypeOf(true),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some}}{{$y := .Some.Some}}{{$z := .}}`,
			expectTplStr: `{{$tpl_x := .Some}}{{$tpl_y := .Some.Some}}{{$tpl_z := .}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type3{Some: type2{Some: ""}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type3{}),
					"$tpl_x": reflect.TypeOf(type2{}),
					"$tpl_y": reflect.TypeOf(""),
					"$tpl_z": reflect.TypeOf(type3{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some}}{{$y := .Some.Some}}{{$z := .}}`,
			expectTplStr: `{{$tpl_x := .Some}}{{$tpl_y := .Some.Some}}{{$tpl_z := .}}`,
			funcs:        defFuncs,
			typecheck:    true,
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
			expectTplStr: `{{$tpl_x := .Some}}{{$tpl_y := $tpl_x.Some}}`,
			funcs:        defFuncs,
			typecheck:    true,
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
			typecheck:    true,
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
			expectTplStr: `{{$tpl_x := .Some.Method}}`,
			funcs:        defFuncs,
			typecheck:    true,
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
			expectTplStr: `{{$tpl_x := .Some.MethodArgs "r"}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type4{Some: type4{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type4{}),
					"$tpl_x": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some}}`,
			expectTplStr: `{{$tpl_x := .Some}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type5{Some: &type3{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type5{}),
					"$tpl_x": reflect.TypeOf(&type3{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some.Some}}`,
			expectTplStr: `{{$tpl_x := .Some.Some}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type5{Some: &type3{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type5{}),
					"$tpl_x": reflect.TypeOf(type2{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := intf}}`,
			expectTplStr: `{{$tpl_x := intf}}`,
			funcs:        defFuncs,
			typecheck:    true,
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(nil),
					"$tpl_x": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{range $i,$v := .Some}}{{range $i, $v := $v.Some}}{{end}}{{end}}`,
			expectTplStr: `{{$var0 := .Some}}{{range $tpl_i, $tpl_v := $var0}}{{$var1 := $tpl_v.Some}}{{range $tpl_i, $tpl_v := $var1}}{{end}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type6{},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type6{}),
					"$var0": reflect.TypeOf([]type1{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf(type1{}),
					"$tpl_i": reflect.TypeOf(1),
					"$tpl_v": reflect.TypeOf(type1{}),
					"$var1":  reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf(""),
					"$tpl_i": reflect.TypeOf(1),
					"$tpl_v": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{range $i,$v := .SomeSlicePtr}}{{range $i, $v := $v.Some}}{{end}}{{end}}`,
			expectTplStr: `{{$var0 := .SomeSlicePtr}}{{range $tpl_i, $tpl_v := $var0}}{{$var1 := $tpl_v.Some}}{{range $tpl_i, $tpl_v := $var1}}{{end}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type6{},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type6{}),
					"$var0": reflect.TypeOf([]*type1{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf(&type1{}),
					"$tpl_i": reflect.TypeOf(1),
					"$tpl_v": reflect.TypeOf(&type1{}),
					"$var1":  reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf(""),
					"$tpl_i": reflect.TypeOf(1),
					"$tpl_v": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{range $i,$v := .Some}}{{range $i, $v := $v.Some}}{{end}}{{end}}`,
			expectTplStr: `{{$var0 := .Some}}{{range $tpl_i, $tpl_v := $var0}}{{$var1 := $tpl_v.Some}}{{range $tpl_i, $tpl_v := $var1}}{{end}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type7{},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type7{}),
					"$var0": reflect.TypeOf([]*type7{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf(&type7{}),
					"$tpl_i": reflect.TypeOf(1),
					"$tpl_v": reflect.TypeOf(&type7{}),
					"$var1":  reflect.TypeOf([]*type7{}),
				},
				map[string]reflect.Type{
					".":      reflect.TypeOf(&type7{}),
					"$tpl_i": reflect.TypeOf(1),
					"$tpl_v": reflect.TypeOf(&type7{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := "e"}}{{define "rr"}}{{$x := "x"}}{{end}}{{template "rr" $x}}`,
			expectTplStr: `{{$tpl_x := "e"}}{{template "rr" $tpl_x}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type7{},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":      reflect.TypeOf(type7{}),
					"$tpl_x": reflect.TypeOf(""),
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

func typechecktemplate(t *template.Template, testData TestData) (*template.Template, *simplifier.State) {
	ret, err := t.Clone()
	if err != nil {
		panic(err)
	}
	var typeCheck *simplifier.State
	for _, t := range ret.Templates() {
		if t.Tree != nil {
			simplifier.Simplify(t.Tree)
			typeCheck = simplifier.TypeCheck(t.Tree, testData.data, testData.funcs)
		}
	}
	return ret, typeCheck
}
