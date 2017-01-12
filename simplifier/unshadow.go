package simplifier

import (
	"fmt"
	"strconv"
	// "github.com/mh-cbon/print-template-tree/printer"
	"text/template/parse"
)

// Unshadow process the tree until no more variable is shadowed.
func Unshadow(tree *parse.Tree) {
	s := &treeUnshadower{
		vars:           []string{},
		currentRenames: map[string]string{},
	}
	s.process(tree)
}

// treeUnshadower holds
// a list of declared variable within the template
// a map of variable rename (original=>new name)
type treeUnshadower struct {
	tree           *parse.Tree
	vars           []string
	currentRenames map[string]string
}

// hasVar tells if a variable is already defined (shadowed)
func (t *treeUnshadower) hasVar(n string) bool {
	for _, s := range t.vars {
		if s == n {
			return true
		}
	}
	return false
}

// registerVar registers a new declared variable of the template.
func (t *treeUnshadower) registerVar(n string) {
	t.vars = append(t.vars, n)
}

// rename checks if given varname is already defined,
// create a new unique variable and returns it.
// if the variable is already unique, it is returned as is.
func (t *treeUnshadower) rename(varname string) string {
	if t.hasVar(varname) {
		n := varname + "_shadow"
		z := n
		i := -1
		for t.hasVar(z) {
			i++
			z = n + strconv.Itoa(i)
		}
		if i > -1 {
			n = z
		}
		t.currentRenames[varname] = n
		varname = n
		t.vars = append(t.vars, varname)
	}
	return varname
}

// getName returns the variable rename, or the variable name.
func (t *treeUnshadower) getName(n string) string {
	if t, ok := t.currentRenames[n]; ok {
		return t
	}
	return n
}

// process the tree until no more unshadowing can be done.
func (t *treeUnshadower) process(tree *parse.Tree) {
	t.tree = tree
	for t.browseToUnshadow(tree.Root) {
	}
}

func (t *treeUnshadower) browseToUnshadow(l interface{}) bool {
	switch node := l.(type) {

	case *parse.ListNode:
		if node != nil {
			for _, child := range node.Nodes {
				if t.browseToUnshadow(child) {
					return true
				}
			}
		}

	case *parse.ActionNode:
		if t.browseToUnshadow(node.Pipe) {
			return true
		}

	case *parse.PipeNode:
		for _, child := range node.Decl {
			if len(child.Ident) == 1 {
				if t.hasVar(child.Ident[0]) == false {
					t.registerVar(child.Ident[0])
				} else {
					child.Ident[0] = t.rename(child.Ident[0])
				}
			} else {
				if t.browseToUnshadow(child) {
					return true
				}
			}
		}
		for _, child := range node.Cmds {
			if t.browseToUnshadow(child) {
				return true
			}
		}

	case *parse.CommandNode:
		for _, child := range node.Args {
			if t.browseToUnshadow(child) {
				return true
			}
		}

	case *parse.RangeNode:
		if t.browseToUnshadow(node.Pipe) {
			return true
		}
		if t.browseToUnshadow(node.List) {
			return true
		}
		if t.browseToUnshadow(node.ElseList) {
			return true
		}

	case *parse.IfNode:
		if t.browseToUnshadow(node.Pipe) {
			return true
		}
		if t.browseToUnshadow(node.List) {
			return true
		}
		if t.browseToUnshadow(node.ElseList) {
			return true
		}

	case *parse.WithNode:
		if t.browseToUnshadow(node.Pipe) {
			return true
		}
		if t.browseToUnshadow(node.List) {
			return true
		}
		if t.browseToUnshadow(node.ElseList) {
			return true
		}
	case *parse.StringNode:
		// pass
	case *parse.IdentifierNode:
		// pass
	case *parse.NumberNode:
		// pass
	case *parse.BoolNode:
		// pass
	case *parse.TextNode:
		// pass
	case *parse.DotNode:
		// pass
	case *parse.VariableNode:
		// do the rename
		node.Ident[0] = t.getName(node.Ident[0])

	default:
		fmt.Printf("%#v\n", node)
		fmt.Printf("!!! Unhandled %T\n", node)
	}
	return false
}
