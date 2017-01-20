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
			expectTplStr: `{{$var0 := .Some}}{{range $tplX := $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf(""),
					"$tplX": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .}}{{range $x.Some}}{{end}}`,
			expectTplStr: `{{$tplX := .}}{{$var0 := $tplX.Some}}{{range $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$tplX": reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{range $x, $y := .Some}}{{end}}`,
			expectTplStr: `{{$var0 := .Some}}{{range $tplX, $tplY := $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf(""),
					"$tplX": reflect.TypeOf(1),
					"$tplY": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{with $x := .Some}}{{end}}`,
			expectTplStr: `{{$var0 := .Some}}{{with $tplX := $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf([]string{}),
					"$tplX": reflect.TypeOf([]string{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .}}{{with $x.Some}}{{end}}`,
			expectTplStr: `{{$tplX := .}}{{$var0 := $tplX.Some}}{{with $var0}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$tplX": reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".": reflect.TypeOf([]string{}),
				},
			},
		},
		TestData{
			tplstr:       `{{with .}}{{$y := .}}{{end}}`,
			expectTplStr: `{{$var0 := .}}{{with $var0}}{{$tplY := .}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf(type1{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$tplY": reflect.TypeOf(type1{}),
				},
			},
		},
		TestData{
			tplstr:       `{{with .}}{{$y := $}}{{end}}`,
			expectTplStr: `{{$var0 := .}}{{with $var0}}{{$tplY := $}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type1{Some: []string{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$var0": reflect.TypeOf(type1{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$tplY": reflect.TypeOf(type1{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := up .Some}}`,
			expectTplStr: `{{$tplX := up .Some}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type2{Some: ""},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type2{}),
					"$tplX": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := up .Some}}{{$y := $x}}`,
			expectTplStr: `{{$tplX := up .Some}}{{$tplY := $tplX}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type2{Some: ""},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type2{}),
					"$tplX": reflect.TypeOf(""),
					"$tplY": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := ""}}`,
			expectTplStr: `{{$tplX := ""}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type2{Some: ""},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type2{}),
					"$tplX": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := 2}}`,
			expectTplStr: `{{$tplX := 2}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type2{Some: ""},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type2{}),
					"$tplX": reflect.TypeOf(2),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := true}}`,
			expectTplStr: `{{$tplX := true}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type2{Some: ""},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type2{}),
					"$tplX": reflect.TypeOf(true),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some}}{{$y := .Some.Some}}{{$z := .}}`,
			expectTplStr: `{{$tplX := .Some}}{{$tplY := .Some.Some}}{{$tplZ := .}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type3{Some: type2{Some: ""}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type3{}),
					"$tplX": reflect.TypeOf(type2{}),
					"$tplY": reflect.TypeOf(""),
					"$tplZ": reflect.TypeOf(type3{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some}}{{$y := .Some.Some}}{{$z := .}}`,
			expectTplStr: `{{$tplX := .Some}}{{$tplY := .Some.Some}}{{$tplZ := .}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type4{Some: type2{Some: ""}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type4{}),
					"$tplX": reflectInterface,
					"$tplY": reflectInterface,
					"$tplZ": reflect.TypeOf(type4{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some}}{{$y := $x.Some}}`,
			expectTplStr: `{{$tplX := .Some}}{{$tplY := $tplX.Some}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type4{Some: type2{Some: ""}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type4{}),
					"$tplX": reflectInterface,
					"$tplY": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Method}}`,
			expectTplStr: `{{$tplX := .Method}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type4{Some: type2{Some: ""}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type4{}),
					"$tplX": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some.Method}}`,
			expectTplStr: `{{$tplX := .Some.Method}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type4{Some: type4{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type4{}),
					"$tplX": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some.MethodArgs "r"}}`,
			expectTplStr: `{{$tplX := .Some.MethodArgs "r"}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type4{Some: type4{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type4{}),
					"$tplX": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some}}`,
			expectTplStr: `{{$tplX := .Some}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type5{Some: &type3{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type5{}),
					"$tplX": reflect.TypeOf(&type3{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := .Some.Some}}`,
			expectTplStr: `{{$tplX := .Some.Some}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type5{Some: &type3{}},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type5{}),
					"$tplX": reflect.TypeOf(type2{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := intf}}`,
			expectTplStr: `{{$tplX := intf}}`,
			funcs:        defFuncs,
			typecheck:    true,
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(nil),
					"$tplX": reflectInterface,
				},
			},
		},
		TestData{
			tplstr:       `{{range $i,$v := .Some}}{{range $i, $v := $v.Some}}{{end}}{{end}}`,
			expectTplStr: `{{$var0 := .Some}}{{range $tplI, $tplV := $var0}}{{$var1 := $tplV.Some}}{{range $tplI, $tplV := $var1}}{{end}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type6{},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type6{}),
					"$var0": reflect.TypeOf([]type1{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf(type1{}),
					"$tplI": reflect.TypeOf(1),
					"$tplV": reflect.TypeOf(type1{}),
					"$var1": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf(""),
					"$tplI": reflect.TypeOf(1),
					"$tplV": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{range $i,$v := .SomeSlicePtr}}{{range $i, $v := $v.Some}}{{end}}{{end}}`,
			expectTplStr: `{{$var0 := .SomeSlicePtr}}{{range $tplI, $tplV := $var0}}{{$var1 := $tplV.Some}}{{range $tplI, $tplV := $var1}}{{end}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type6{},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type6{}),
					"$var0": reflect.TypeOf([]*type1{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf(&type1{}),
					"$tplI": reflect.TypeOf(1),
					"$tplV": reflect.TypeOf(&type1{}),
					"$var1": reflect.TypeOf([]string{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf(""),
					"$tplI": reflect.TypeOf(1),
					"$tplV": reflect.TypeOf(""),
				},
			},
		},
		TestData{
			tplstr:       `{{range $i,$v := .Some}}{{range $i, $v := $v.Some}}{{end}}{{end}}`,
			expectTplStr: `{{$var0 := .Some}}{{range $tplI, $tplV := $var0}}{{$var1 := $tplV.Some}}{{range $tplI, $tplV := $var1}}{{end}}{{end}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type7{},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type7{}),
					"$var0": reflect.TypeOf([]*type7{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf(&type7{}),
					"$tplI": reflect.TypeOf(1),
					"$tplV": reflect.TypeOf(&type7{}),
					"$var1": reflect.TypeOf([]*type7{}),
				},
				map[string]reflect.Type{
					".":     reflect.TypeOf(&type7{}),
					"$tplI": reflect.TypeOf(1),
					"$tplV": reflect.TypeOf(&type7{}),
				},
			},
		},
		TestData{
			tplstr:       `{{$x := "e"}}{{define "rr"}}{{$x := "x"}}{{end}}{{template "rr" $x}}`,
			expectTplStr: `{{$tplX := "e"}}{{template "rr" $tplX}}`,
			funcs:        defFuncs,
			typecheck:    true,
			data:         type7{},
			checkedTypes: []map[string]reflect.Type{
				map[string]reflect.Type{
					".":     reflect.TypeOf(type7{}),
					"$tplX": reflect.TypeOf(""),
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
