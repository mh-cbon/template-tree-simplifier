package simplifier

import (
	"fmt"
	"reflect"
	"text/template/parse"
)

// TypeCheck browses the tree to identify variable types.
func TypeCheck(tree *parse.Tree, data interface{}, funcs map[string]interface{}) *State {
	s := &State{
		currentScope: -1,
		vars:         []map[string]reflect.Type{},
	}
	t := &treeTypecheck{
		funcs: funcs,
		tree:  tree,
	}
	s.Add()
	s.Enter()
	s.AddVar(".", reflect.TypeOf(data))
	t.process(tree, s)
	s.Leave()
	return s
}

// treeTypecheck ...
type treeTypecheck struct {
	tree  *parse.Tree
	funcs map[string]interface{}
}

// State ...
type State struct {
	currentScope int
	vars         []map[string]reflect.Type
}

// Add a new scope level.
func (s *State) Add() {
	s.vars = append(s.vars, map[string]reflect.Type{})
}

// Enter into a scope level.
func (s *State) Enter() {
	s.currentScope++
}

// Len returns the number of scopes.
func (s *State) Len() int {
	return len(s.vars)
}

// Level returns the current scope level.
func (s *State) Level() int {
	return s.currentScope
}

// Leave a scope level.
func (s *State) Leave() {
	s.currentScope--
}

// Current scope vars.
func (s *State) Current() map[string]reflect.Type {
	return s.vars[s.currentScope]
}

// Root scope vars.
func (s *State) Root() map[string]reflect.Type {
	return s.vars[0]
}

// Dot is the current scope dot.
func (s *State) Dot() reflect.Type {
	return s.Current()["."]
}

// RootDot is the root scope dot.
func (s *State) RootDot() reflect.Type {
	return s.Root()["."]
}

// AddVar in the current scope level.
func (s *State) AddVar(name string, r reflect.Type) {
	s.Current()[name] = r
}

// HasVar tells if current level contains given variable name.
func (s *State) HasVar(name string) bool {
	_, ok := s.Current()[name]
	return ok
}

// GetVar get a variable in current scope.
func (s *State) GetVar(name string) reflect.Type {
	r, _ := s.Current()[name]
	return r
}

// FindVar starting from current scope level to the root.
func (s *State) FindVar(name string) reflect.Type {
	for i := s.currentScope; i >= 0; i-- {
		if v, ok := s.vars[i][name]; ok {
			return v
		}
	}
	if name == "$" {
		return s.RootDot()
	}
	return nil
}

// BrowsePathType ...
func (s *State) BrowsePathType(path []string, val reflect.Type) reflect.Type {
	for _, p := range path {
		if val.Kind() == reflect.Interface {
			return val
		}
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		field, found := val.FieldByName(p)
		if !found {
			meth, found := val.MethodByName(p)
			if !found {
				err := fmt.Errorf("State.BrowsePathType: path %v not found in type %v", path, val)
				panic(err)
			}
			val = meth.Type.Out(0)
		} else {
			val = field.Type
		}
	}
	return val
}

// IsMethodPath ...
func (s *State) IsMethodPath(path []string, val reflect.Type) bool {
	for _, p := range path {
		if val.Kind() == reflect.Interface {
			return false
		}
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() == reflect.Struct {
			field, found := val.FieldByName(p)
			if !found {
				_, found := val.MethodByName(p)
				return found
			}
			val = field.Type
		}
	}
	return false
}

// ReflectPath ...
func (s *State) ReflectPath(path []string, val reflect.Type) reflect.Type {
	for _, p := range path {
		if val.Kind() == reflect.Interface {
			return val
		}
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		field, found := val.FieldByName(p)
		if !found {
			meth, found := val.MethodByName(p)
			if found {
				return meth.Type
			}
			return nil
		}
		val = field.Type
	}
	return val
}

// process the tree until no more simplification can be done.
func (t *treeTypecheck) process(tree *parse.Tree, state *State) {
	t.browseNodes(tree.Root, state)
}

// browseNodes recursively.
func (t *treeTypecheck) browseNodes(l interface{}, state *State) {
	switch node := l.(type) {

	case *parse.ListNode:
		if node != nil {
			for _, child := range node.Nodes {
				t.browseNodes(child, state)
			}
		}

	case *parse.ActionNode:
		t.typeCheckActionNode(node, state)
		t.browseNodes(node.Pipe, state)

	case *parse.RangeNode:
		t.enterRangeNode(node, state)
		t.browseNodes(node.Pipe, state)
		t.browseNodes(node.List, state)
		t.browseNodes(node.ElseList, state)
		state.Leave()

	case *parse.IfNode:
		t.browseNodes(node.Pipe, state)
		t.browseNodes(node.List, state)
		t.browseNodes(node.ElseList, state)

	case *parse.WithNode:
		t.enterWithNode(node, state)
		t.browseNodes(node.Pipe, state)
		t.browseNodes(node.List, state)
		t.browseNodes(node.ElseList, state)
		state.Leave()

	case *parse.TemplateNode:
		if node.Pipe != nil {
			t.browseNodes(node.Pipe, state)
		}

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
		err := fmt.Errorf("treeTypecheck.browseNodes: unhandled node type\n%v\n%#v", node, node)
		panic(err)
	}
}

func (t *treeTypecheck) typeCheckActionNode(node *parse.ActionNode, state *State) bool {
	/*
		looks for
		{{$some := .Field.Node}}
		{{$some := $y}}
		{{$some := $y.Field}}
		{{$some := .}}
		{{$some := up $w}}
		{{$some := "what"}}
		{{$some := 4}}
		{{$some := true}}
	*/
	if len(node.Pipe.Decl) > 0 && len(node.Pipe.Decl[0].Ident) == 1 {
		varName := node.Pipe.Decl[0].Ident[0]
		if len(node.Pipe.Cmds) == 1 && len(node.Pipe.Cmds[0].Args) > 0 {
			if field, ok := node.Pipe.Cmds[0].Args[0].(*parse.FieldNode); ok {
				r := state.BrowsePathType(field.Ident, state.Dot())
				state.AddVar(varName, r)

			} else if variable, ok := node.Pipe.Cmds[0].Args[0].(*parse.VariableNode); ok {
				rightVarType := state.FindVar(variable.Ident[0])
				if rightVarType == nil {
					panic(fmt.Errorf("%v\nVariable not found %v in %v", t.tree.Root.String(), variable.Ident[0], node))
				}
				if len(variable.Ident) > 1 {
					rightVarType = state.BrowsePathType(variable.Ident[1:], rightVarType)
				}
				state.AddVar(varName, rightVarType)

			} else if _, ok := node.Pipe.Cmds[0].Args[0].(*parse.DotNode); ok {
				rightVarType := state.Dot()
				state.AddVar(varName, rightVarType)

			} else if ident, ok := node.Pipe.Cmds[0].Args[0].(*parse.IdentifierNode); ok {
				funcRetType := t.getFuncValueType(ident.Ident)
				state.AddVar(varName, funcRetType)

			} else if _, ok := node.Pipe.Cmds[0].Args[0].(*parse.StringNode); ok {
				state.AddVar(varName, reflect.TypeOf(""))

			} else if _, ok := node.Pipe.Cmds[0].Args[0].(*parse.NumberNode); ok {
				state.AddVar(varName, reflect.TypeOf(0))

			} else if _, ok := node.Pipe.Cmds[0].Args[0].(*parse.BoolNode); ok {
				state.AddVar(varName, reflect.TypeOf(true))

			}
		} else {
			err := fmt.Errorf("treeTypecheck.typeCheckActionNode: unhandled length of node.Pipe.Decl or node.Pipe.Cmds\n%v\n%#v", node, node)
			panic(err)
		}
	}
	return false
}

func (t *treeTypecheck) enterRangeNode(node *parse.RangeNode, state *State) bool {
	var newDotType reflect.Type
	if len(node.Pipe.Cmds) == 1 {
		if variable, ok := node.Pipe.Cmds[0].Args[0].(*parse.VariableNode); ok {
			//-
			rightVarType := state.FindVar(variable.Ident[0])
			if len(variable.Ident) > 1 {
				rightVarType = state.BrowsePathType(variable.Ident[1:], rightVarType)
			}
			newDotType = rightVarType

		} else if _, ok := node.Pipe.Cmds[0].Args[0].(*parse.DotNode); ok {
			newDotType = state.Dot()

		} else {
			err := fmt.Errorf("treeTypecheck.enterRangeNode: unhandled type of Arg[0]\n%v\n%#v", node, node)
			panic(err)
		}
	} else {
		err := fmt.Errorf("treeTypecheck.enterRangeNode: unhandled length of node.Pipe.Cmds\n%v\n%#v", node, node)
		panic(err)
	}
	if newDotType == nil {
		err := fmt.Errorf("treeTypecheck.enterRangeNode: new dot type not found\n%v\n%#v", node, node)
		panic(err)
	}
	state.Add()
	state.Enter()
	state.AddVar(".", newDotType.Elem())
	if len(node.Pipe.Decl) > 0 {
		// add the new var to the new scope
		if len(node.Pipe.Decl) == 1 {
			state.AddVar(node.Pipe.Decl[0].Ident[0], state.Dot())

		} else {
			state.AddVar(node.Pipe.Decl[0].Ident[0], reflect.TypeOf(1))
			state.AddVar(node.Pipe.Decl[1].Ident[0], state.Dot())
		}
	}
	return false
}

func (t *treeTypecheck) enterWithNode(node *parse.WithNode, state *State) bool {
	var newDotType reflect.Type
	if len(node.Pipe.Cmds) == 1 {
		if variable, ok := node.Pipe.Cmds[0].Args[0].(*parse.VariableNode); ok {
			//-
			rightVarType := state.FindVar(variable.Ident[0])
			if len(variable.Ident) > 1 {
				rightVarType = state.BrowsePathType(variable.Ident[1:], rightVarType)
			}
			newDotType = rightVarType

		} else if _, ok := node.Pipe.Cmds[0].Args[0].(*parse.DotNode); ok {
			newDotType = state.Dot()

		} else {
			err := fmt.Errorf("treeTypecheck.enterWithNode: unhandled type of Arg[0]\n%v\n%#v", node, node)
			panic(err)
		}
	} else {
		err := fmt.Errorf("treeTypecheck.enterWithNode: unhandled length of node.Pipe.Cmds\n%v\n%#v", node, node)
		panic(err)
	}
	if newDotType == nil {
		err := fmt.Errorf("treeTypecheck.enterWithNode: new dot type not found\n%v\n%#v", node, node)
		panic(err)
	}
	state.Add()
	state.Enter()
	state.AddVar(".", newDotType)
	if len(node.Pipe.Decl) > 0 {
		// add the new var to the new scope
		if len(node.Pipe.Decl) == 1 {
			state.AddVar(node.Pipe.Decl[0].Ident[0], state.Dot())
		} else {
			err := fmt.Errorf("treeTypecheck.enterWithNode: unhandled length of node.Pipe.Decl\n%v\n%#v", node, node)
			panic(err)
		}
	}
	return false
}

func (t *treeTypecheck) getFuncValueType(name string) reflect.Type {
	if f, ok := t.funcs[name]; ok {
		fR := reflect.TypeOf(f)
		if fR.NumOut() > 0 {
			return fR.Out(0)
		}
	}
	return nil
}
