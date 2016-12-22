package simplifier_test

import (
	"bytes"
	"fmt"
	"github.com/mh-cbon/template-tree-simplifier/simplifier"
	"strings"
	"testing"
	"text/template"
)

type TestData struct {
	tplstr       string
	data         interface{}
	expectTplStr string
	funcs        template.FuncMap
}

func TestAll(t *testing.T) {
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
			tplstr:       `{{"son" | split "wat"}}`,
			expectTplStr: `{{split "wat" "son"}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{. | up}}`,
			expectTplStr: `{{up .}}`,
			funcs:        defFuncs,
			data:         "hello",
		},
		TestData{
			tplstr:       `{{.S | up}}`,
			expectTplStr: `{{up .S}}`,
			funcs:        defFuncs,
			data:         struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr:       `{{1 | incr}}`,
			expectTplStr: `{{incr 1}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{1 | incr | incr | mul 2}}`,
			expectTplStr: `{{$var0 := incr 1}}{{$var1 := incr $var0}}{{mul 2 $var1}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{mul (1 | incr) (1 | incr)}}`,
			expectTplStr: `{{$var0 := incr 1}}{{$var1 := incr 1}}{{mul $var0 $var1}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{"some" | split ("what" | up)}}`,
			expectTplStr: `{{$var0 := up "what"}}{{split $var0 "some"}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{"some" | split (.S | up)}}`,
			expectTplStr: `{{$var0 := up .S}}{{split $var0 "some"}}`,
			funcs:        defFuncs,
			data:         struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr:       `{{"some" | split (("what" | lower) | up)}}`,
			expectTplStr: `{{$var1 := lower "what"}}{{$var0 := up $var1}}{{split $var0 "some"}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{"some" | split ((.S | lower) | up)}}`,
			expectTplStr: `{{$var1 := lower .S}}{{$var0 := up $var1}}{{split $var0 "some"}}`,
			funcs:        defFuncs,
			data:         struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr:       `{{("what" | lower) | split "" | join "" | up}}`,
			expectTplStr: `{{$var0 := lower "what"}}{{$var1 := split "" $var0}}{{$var2 := join "" $var1}}{{up $var2}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{up "what" | lower}}`,
			expectTplStr: `{{$var0 := up "what"}}{{lower $var0}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{up .S | lower}}`,
			expectTplStr: `{{$var0 := up .S}}{{lower $var0}}`,
			funcs:        defFuncs,
			data:         struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr:       `{{$t := ("what" | up)}}`,
			expectTplStr: `{{$tpl_t := (up "what")}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{$t := (.S | up)}}`,
			expectTplStr: `{{$tpl_t := (up .S)}}`,
			funcs:        defFuncs,
			data:         struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr: `{{$t := ("what" | up)}}
{{$k := up $t}}`,
			expectTplStr: `{{$tpl_t := (up "what")}}
{{$tpl_k := up $tpl_t}}`,
			funcs: defFuncs,
		},
		TestData{
			tplstr: `{{$t := (.S | up)}}
{{$k := up $t}}`,
			expectTplStr: `{{$tpl_t := (up .S)}}
{{$tpl_k := up $tpl_t}}`,
			funcs: defFuncs,
			data:  struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr:       `{{if true}}{{end}}`,
			expectTplStr: `{{if true}}{{end}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{if eq 1 2}}{{end}}`,
			expectTplStr: `{{$var0 := eq 1 2}}{{if $var0}}{{end}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr: `{{if true}}
  {{"some" | split ("what" | up)}}
{{end}}`,
			expectTplStr: `{{if true}}
  {{$var0 := up "what"}}{{split $var0 "some"}}
{{end}}`,
			funcs: defFuncs,
		},
		TestData{
			tplstr: `{{if true}}
  {{"some" | split (.S | up)}}
{{end}}`,
			expectTplStr: `{{if true}}
  {{$var0 := up .S}}{{split $var0 "some"}}
{{end}}`,
			funcs: defFuncs,
			data:  struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr:       `{{if eq ("what" | up | lower) "what"}}{{end}}`,
			expectTplStr: `{{$var2 := up "what"}}{{$var1 := lower $var2}}{{$var0 := eq $var1 "what"}}{{if $var0}}{{end}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{if eq (.S | up | lower) .S}}{{end}}`,
			expectTplStr: `{{$var2 := up .S}}{{$var1 := lower $var2}}{{$var0 := eq $var1 .S}}{{if $var0}}{{end}}`,
			funcs:        defFuncs,
			data:         struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr: `{{if eq ("what" | up | lower) "what"}}
  {{"some" | split ("what" | up)}}
{{end}}`,
			expectTplStr: `{{$var2 := up "what"}}{{$var1 := lower $var2}}{{$var0 := eq $var1 "what"}}{{if $var0}}
  {{$var3 := up "what"}}{{split $var3 "some"}}
{{end}}`,
			funcs: defFuncs,
		},
		TestData{
			tplstr: `{{if eq (.S | up | lower) "what"}}
  {{.S | split ("what" | up)}}
{{end}}`,
			expectTplStr: `{{$var2 := up .S}}{{$var1 := lower $var2}}{{$var0 := eq $var1 "what"}}{{if $var0}}
  {{$var3 := up "what"}}{{split $var3 .S}}
{{end}}`,
			funcs: defFuncs,
			data:  struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr: `{{if not (eq ("what" | up | lower) "what")}}
  {{"some" | split ("what" | up)}}
{{end}}`,
			expectTplStr: `{{$var3 := up "what"}}{{$var2 := lower $var3}}{{$var1 := eq $var2 "what"}}{{$var0 := not $var1}}{{if $var0}}
  {{$var4 := up "what"}}{{split $var4 "some"}}
{{end}}`,
			funcs: defFuncs,
		},
		TestData{
			tplstr: `{{if not (eq (.S | up | lower) "what")}}
  {{"some" | split (.S | up)}}
{{end}}`,
			expectTplStr: `{{$var3 := up .S}}{{$var2 := lower $var3}}{{$var1 := eq $var2 "what"}}{{$var0 := not $var1}}{{if $var0}}
  {{$var4 := up .S}}{{split $var4 "some"}}
{{end}}`,
			funcs: defFuncs,
			data:  struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr:       `{{$var0 := up "what" | lower}}{{$var0}}`,
			expectTplStr: `{{$var0 := up "what"}}{{$tpl_var0 := lower $var0}}{{$tpl_var0}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{$var0 := up .S | lower}}{{$var0}}`,
			expectTplStr: `{{$var0 := up .S}}{{$tpl_var0 := lower $var0}}{{$tpl_var0}}`,
			funcs:        defFuncs,
			data:         struct{ S string }{S: "hello"},
		},
		TestData{
			tplstr:       `{{if not (eq "a" "b")}}not eq{{end}}`,
			expectTplStr: `{{$var1 := eq "a" "b"}}{{$var0 := not $var1}}{{if $var0}}not eq{{end}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{if eq ("what" |lower) ("what" |lower)}}{{end}}`,
			expectTplStr: `{{$var1 := lower "what"}}{{$var2 := lower "what"}}{{$var0 := eq $var1 $var2}}{{if $var0}}{{end}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{$var0 := eq (lower "up") "what"}}`,
			expectTplStr: `{{$var0 := lower "up"}}{{$tpl_var0 := eq $var0 "what"}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr: `{{range .List}}
{{.}}
{{end}}`,
			expectTplStr: `{{range .List}}
{{.}}
{{end}}`,
			funcs: defFuncs,
			data:  struct{ List []string }{List: []string{"what"}},
		},
		TestData{
			tplstr: `{{range split "what" ""}}
{{.}}
{{end}}`,
			expectTplStr: `{{$var0 := split "what" ""}}{{range $var0}}
{{.}}
{{end}}`,
			funcs: defFuncs,
			data:  struct{ List []string }{List: []string{"what"}},
		},
		TestData{
			tplstr: `{{range $i,$v := split "what" ""}}
{{$i}} -> {{$v}}
{{end}}`,
			expectTplStr: `{{$var0 := split "what" ""}}{{range $tpl_i, $tpl_v := $var0}}
{{$tpl_i}} -> {{$tpl_v}}
{{end}}`,
			funcs: defFuncs,
			data:  struct{ List []string }{List: []string{"what"}},
		},
		TestData{
			tplstr: `{{range $i,$v := "some" | split (("what" | lower) | up)}}
{{$i}} -> {{$v}}
{{end}}`,
			expectTplStr: `{{$var2 := lower "what"}}{{$var1 := up $var2}}{{$var0 := split $var1 "some"}}{{range $tpl_i, $tpl_v := $var0}}
{{$tpl_i}} -> {{$tpl_v}}
{{end}}`,
			funcs: defFuncs,
			data:  struct{ List []string }{List: []string{"what"}},
		},
		TestData{
			tplstr: `{{if false}}
{{"some" | split (("what" | lower) | up)}}
{{else}}
{{"some" | split (("what" | lower) | up)}}
{{end}}`,
			expectTplStr: `{{if false}}
{{$var1 := lower "what"}}{{$var0 := up $var1}}{{split $var0 "some"}}
{{else}}
{{$var3 := lower "what"}}{{$var2 := up $var3}}{{split $var2 "some"}}
{{end}}`,
			funcs: defFuncs,
			data:  struct{ List []string }{List: []string{"what"}},
		},
		TestData{
			tplstr: `{{if false}}
  {{"some" | split (("what" | lower) | up)}}
{{else}}
  {{"some" | split (("what" | lower) | up)}}
  {{if false}}
    {{"some" | split (("what" | lower) | up)}}
  {{else}}
    {{"some" | split (("what" | lower) | up)}}
  {{end}}
{{end}}`,
			expectTplStr: `{{if false}}
  {{$var1 := lower "what"}}{{$var0 := up $var1}}{{split $var0 "some"}}
{{else}}
  {{$var3 := lower "what"}}{{$var2 := up $var3}}{{split $var2 "some"}}
  {{if false}}
    {{$var5 := lower "what"}}{{$var4 := up $var5}}{{split $var4 "some"}}
  {{else}}
    {{$var7 := lower "what"}}{{$var6 := up $var7}}{{split $var6 "some"}}
  {{end}}
{{end}}`,
			funcs: defFuncs,
			data:  struct{ List []string }{List: []string{"what"}},
		},
		TestData{
			tplstr: `{{range .List}}
  This is the range branch
  {{"some" | split (("what" | lower) | up)}}
{{else}}
  This is the else branch
  {{"some" | split (("what" | lower) | up)}}
{{end}}`,
			expectTplStr: `{{range .List}}
  This is the range branch
  {{$var1 := lower "what"}}{{$var0 := up $var1}}{{split $var0 "some"}}
{{else}}
  This is the else branch
  {{$var3 := lower "what"}}{{$var2 := up $var3}}{{split $var2 "some"}}
{{end}}`,
			funcs: defFuncs,
			data:  struct{ List []string }{List: []string{"what"}},
		},
		TestData{
			tplstr:       `{{with $x := "output"}}{{. | up}}{{end}}`,
			expectTplStr: `{{with $x := "output"}}{{up .}}{{end}}`,
			funcs:        defFuncs,
		},
		TestData{
			tplstr:       `{{with $x := "output"}}{{$ | up}}{{end}}`,
			expectTplStr: `{{with $x := "output"}}{{up $}}{{end}}`,
			funcs:        defFuncs,
			data:         "hello",
		},
		TestData{
			tplstr:       `{{with $x := "output" | lower}}{{$ | up}}{{. | up}}{{end}}`,
			expectTplStr: `{{$var0 := lower "output"}}{{with $x := $var0}}{{up $}}{{up .}}{{end}}`,
			funcs:        defFuncs,
			data:         "hello",
		},
		TestData{
			tplstr:       `{{with $x := "output" | lower}}{{$ | up}}{{. | up | lower}}{{end}}`,
			expectTplStr: `{{$var0 := lower "output"}}{{with $x := $var0}}{{up $}}{{$var1 := up .}}{{lower $var1}}{{end}}`,
			funcs:        defFuncs,
			data:         "hello",
		},
		TestData{
			tplstr:       `{{with $x := "" | lower}}{{else}}{{$ | up}}{{. | up | lower}}{{end}}`,
			expectTplStr: `{{$var0 := lower ""}}{{with $x := $var0}}{{else}}{{up $}}{{$var1 := up .}}{{lower $var1}}{{end}}`,
			funcs:        defFuncs,
			data:         "hello",
		},
	}

	for _, testData := range testTable {
		// ensure the template is valid and working
		tpl, err := template.New("").Funcs(testData.funcs).Parse(testData.tplstr)
		if err != nil {
			t.Logf("ORIGINAL TEMPLATE:\n%v\n", testData.tplstr)
			t.Errorf("error while compiling original template: %v", err)
			break
		}
		// execute template, check everything is still fine
		originalOut, err := exectemplate(tpl, testData.data)
		if err != nil {
			t.Logf("ORIGINAL TEMPLATE:\n%v\n", testData.tplstr)
			t.Errorf("error while executing original template: %v", err)
			break
		}
		// do the simplification
		simplified := simplifytemplate(tpl)
		// execute simplified template, check everything is still fine
		simplifiedOut, err := exectemplate(simplified, testData.data)
		if err != nil {
			t.Logf("ORIGINAL TEMPLATE:\n%v\n", testData.tplstr)
			t.Logf("SIMPLIFIED TEMPLATE\n%v\n", simplified.Tree.Root.String())
			t.Errorf("error while executing simplified template: %v", err)
			break
		}
		// ensure both output are eq
		if originalOut != simplifiedOut {
			t.Errorf("Output are different\nORIGINAL\n%v\nSIMPLIFIED\n%v\n",
				originalOut, simplifiedOut)
		}
		// ensure the new template matches expected simplified template
		s := fmt.Sprintf("%v", simplified.Tree.Root.String())
		if s != testData.expectTplStr {
			t.Errorf("Simplified template is not as expected\nEXPECTED\n%v\nSIMPLIFIED\n%v\n",
				testData.expectTplStr, s)
		}
	}
}

func simplifytemplate(t *template.Template) *template.Template {
	ret, err := t.Clone()
	if err != nil {
		panic(err)
	}
	for _, t := range ret.Templates() {
		if t.Tree != nil {
			simplifier.Simplify(t.Tree)
		}
	}
	return ret
}

func exectemplate(t *template.Template, data interface{}) (string, error) {
	var b bytes.Buffer
	return b.String(), t.Execute(&b, data)
}
