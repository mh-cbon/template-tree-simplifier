package simplifier

import (
	"fmt"
	html "html/template"
	text "text/template"
	"text/template/parse"
)

// Transform fully simplify a template.
// it accepts *text.Template or *html.Template,
// it panics if the value type is unexpected.
func Transform(some interface{}, data interface{}, funcs map[string]interface{}) {
	if t, ok := some.(*text.Template); ok {
		for _, tpl := range t.Templates() {
			if tpl.Tree != nil {
				TransformTree(tpl.Tree, data, funcs)
			}
		}

	} else if h, ok := some.(*html.Template); ok {
		for _, tpl := range h.Templates() {
			if tpl.Tree != nil {
				TransformTree(tpl.Tree, data, funcs)
			}
		}
	} else {
		err := fmt.Errorf("Transform: unhandled template type %T", some)
		panic(err)
	}
}

// TransformTree fully simplify a template Tree.
func TransformTree(tree *parse.Tree, data interface{}, funcs map[string]interface{}) *State {
	Unshadow(tree)
	Simplify(tree)
	typeCheck := TypeCheck(tree, data, funcs)
	Unhole(tree, typeCheck, funcs)
	return typeCheck
}
