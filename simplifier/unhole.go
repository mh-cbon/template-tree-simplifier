package simplifier

import (
	"fmt"
	"reflect"
	"strings"
	"text/template/parse"
)

// Unhole process the tree until no more type holes subsist in the tree.
// a type hole being an interface{} value found inside a property path.
// Exemple:
// there is a call such a.b.c.d
// if b is an interface{} value,
// there is no way to statically type check the subsequent c.d path,
// it can only be done at runtime via reflection.
// This process will solve that problem
// by identifying such situation and isolate
// type checkable and type uncheckable parts.
// so such node {{a.b.c.d}} will become {{browsePropertyPath a.b "c.d"}}.
// Where browsePropertyPath is a new identifier (func of funcmap).
// Note1: the case may occur with variable/identifier nodes too,
// lets imagine {{$z := a}}{{$z.b.c.d}}
// Note2: this will work only after simplify occured
func Unhole(tree *parse.Tree, state *State, funcs map[string]interface{}) {
	unhole := &treeUnhole{tree: tree, funcs: funcs}
	unhole.process(state)
}

// treeTypecheck ...
type treeUnhole struct {
	tree  *parse.Tree
	funcs map[string]interface{}
}

// process the tree until no more simplification can be done.
func (t *treeUnhole) process(state *State) {
	state.Enter()
	t.browseNodes(t.tree.Root, state)
	state.Leave()
}

// browseNodes recursively.
func (t *treeUnhole) browseNodes(l interface{}, state *State) {
	switch node := l.(type) {

	case *parse.ListNode:
		if node != nil {
			for _, child := range node.Nodes {
				t.browseNodes(child, state)
			}
		}

	case *parse.ActionNode:
		t.unholeActionNode(node, state)
		t.browseNodes(node.Pipe, state)

	case *parse.RangeNode:
		state.Enter()
		t.browseNodes(node.Pipe, state)
		t.browseNodes(node.List, state)
		t.browseNodes(node.ElseList, state)
		state.Leave()

	case *parse.IfNode:
		t.browseNodes(node.Pipe, state)
		t.browseNodes(node.List, state)
		t.browseNodes(node.ElseList, state)

	case *parse.WithNode:
		state.Enter()
		t.browseNodes(node.Pipe, state)
		t.browseNodes(node.List, state)
		t.browseNodes(node.ElseList, state)
		state.Leave()

	case *parse.PipeNode:
		for _, c := range node.Decl {
			t.browseNodes(c, state)
		}
		for _, c := range node.Cmds {
			t.browseNodes(c, state)
		}

	case *parse.CommandNode:
		for _, c := range node.Args {
			t.browseNodes(c, state)
		}

	case *parse.VariableNode:
		//pass
	case *parse.IdentifierNode:
		//pass
	case *parse.StringNode:
		//pass
	case *parse.NumberNode:
		//pass
	case *parse.BoolNode:
		//pass
	case *parse.DotNode:
		//pass
	case *parse.FieldNode:
		//pass
	case *parse.TextNode:
		//pass

	default:
		fmt.Printf("%#v\n", node)
		fmt.Printf("!!! Unhandled %T\n", node)
		panic("unhandled")
	}
}

func (t *treeUnhole) unholeActionNode(node *parse.ActionNode, state *State) {
	var x []interface{}
	reflectInterface := reflect.TypeOf(x).Elem()
	/*
					  look for
					  {{$some := .b.c}}
					  {{$some := $x.b.c}}
				    check its property path types for an interface{}
		        suppose .b is interface{}
		        it transforms into
		        {{$some := browse .b "c"}}
		        or
		        {{$some := browse $x.b "c"}}
	*/
	if len(node.Pipe.Decl) == 1 && len(node.Pipe.Cmds) == 1 {
		decl := node.Pipe.Decl[0]
		if variable, ok := node.Pipe.Cmds[0].Args[0].(*parse.VariableNode); ok && len(variable.Ident) > 1 {
			if state.GetVar(decl.Ident[0]) == reflectInterface {
				typedPath, unTypedPath := splitTypedPath(variable.Ident, state.GetVar(variable.Ident[0]))
				if len(unTypedPath) > 0 {
					args := []parse.Node{}
					i := &parse.IdentifierNode{
						NodeType: parse.NodeIdentifier,
						Ident:    "browsePropertyPath",
					}
					args = append(args, i)
					if len(typedPath) == 0 {
						v := &parse.VariableNode{
							NodeType: parse.NodeVariable,
							Ident:    variable.Ident[:1],
						}
						args = append(args, v)
						t := &parse.StringNode{
							NodeType: parse.NodeString,
							Text:     strings.Join(unTypedPath[1:], "."),
							Quoted:   "\"" + strings.Join(unTypedPath[1:], ".") + "\"",
						}
						args = append(args, t)
					} else {
						v := &parse.VariableNode{
							NodeType: parse.NodeVariable,
							Ident:    typedPath,
						}
						args = append(args, v)
						t := &parse.StringNode{
							NodeType: parse.NodeString,
							Text:     strings.Join(unTypedPath, "."),
							Quoted:   "\"" + strings.Join(unTypedPath, ".") + "\"",
						}
						args = append(args, t)
					}
					identArgs := node.Pipe.Cmds[0].Args[1:]
					args = append(args, identArgs...)
					node.Pipe.Cmds[0].Args = append(node.Pipe.Cmds[0].Args[:0], args...)
				}
			}
		} else if field, ok := node.Pipe.Cmds[0].Args[0].(*parse.FieldNode); ok && len(field.Ident) > 1 {
			if state.GetVar(decl.Ident[0]) == reflectInterface {
				typedPath, unTypedPath := splitTypedPath(field.Ident, state.Dot())
				if len(unTypedPath) > 0 {
					args := []parse.Node{}
					i := &parse.IdentifierNode{
						NodeType: parse.NodeIdentifier,
						Ident:    "browsePropertyPath",
					}
					args = append(args, i)
					if len(typedPath) > 0 {
						v := &parse.FieldNode{
							NodeType: parse.NodeField,
							Ident:    typedPath,
						}
						args = append(args, v)
					} else {
						v := &parse.DotNode{
							NodeType: parse.NodeDot,
						}
						args = append(args, v)
					}
					t := &parse.StringNode{
						NodeType: parse.NodeString,
						Text:     strings.Join(unTypedPath, "."),
						Quoted:   "\"" + strings.Join(unTypedPath, ".") + "\"",
					}
					args = append(args, t)
					identArgs := node.Pipe.Cmds[0].Args[1:]
					args = append(args, identArgs...)
					node.Pipe.Cmds[0].Args = append(node.Pipe.Cmds[0].Args[:0], args...)
				}
			}
		}
	} else {
		panic("unhandled")
	}
}

func splitTypedPath(path []string, val reflect.Type) ([]string, []string) {
	for i, p := range path {
		if val.Kind() == reflect.Interface {
			return path[:i], path[i:]
		}
		field, found := val.FieldByName(p)
		if found {
			if field.Type.Kind() == reflect.Struct {
				val = field.Type
			} else {
				return path[:i], path[i:]
			}
		} else {
			meth, found := val.MethodByName(p)
			if found {
				if meth.Type.NumOut() == 0 {
					panic("method is void..")
				}
				if meth.Type.Out(0).Kind() == reflect.Struct {
					val = field.Type
				} else {
					panic("unhandled method return value type")
				}
			} else {
				panic("field/method not found")
			}
		}
	}
	return path, []string{}
}
