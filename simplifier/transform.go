package simplifier

import (
	"fmt"
	html "html/template"
	text "text/template"
)

// Transform fully simplify a template Tree.
// it accepts *text.Template or *html.Template,
// it panics if the value type is unexpected.
func Transform(some interface{}, data interface{}, funcs map[string]interface{}) {
	if t, ok := some.(*text.Template); ok {
		for _, tpl := range t.Templates() {
			if tpl.Tree != nil {
				Unshadow(tpl.Tree)
				Simplify(tpl.Tree)
				typeCheck := TypeCheck(tpl.Tree, data, funcs)
				Unhole(tpl.Tree, typeCheck, funcs)
			}
		}

	} else if h, ok := some.(*html.Template); ok {
		for _, tpl := range h.Templates() {
			if tpl.Tree != nil {
				Unshadow(tpl.Tree)
				Simplify(tpl.Tree)
				typeCheck := TypeCheck(tpl.Tree, data, funcs)
				Unhole(tpl.Tree, typeCheck, funcs)
			}
		}
	} else {
		panic(fmt.Errorf("Wrong type %T", some))
	}
}
