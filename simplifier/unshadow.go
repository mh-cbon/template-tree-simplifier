package simplifier

import (
	"fmt"
	"strconv"
	// "github.com/mh-cbon/print-template-tree/printer"
	"text/template/parse"
)

// Unshadow process the tree until no more variable is shadowed.
func Unshadow(tree *parse.Tree) {
	t := &treeUnshadower{
		vars:           []string{},
		currentRenames: map[string]string{},
	}
	t.tree = tree
	t.browseToUnshadow(tree.Root)
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

// browseToUnshadow browses all nodes and unshadow variable
func (t *treeUnshadower) browseToUnshadow(l interface{}) {
	switch node := l.(type) {

	case *parse.ListNode:
		if node != nil {
			for _, child := range node.Nodes {
				t.browseToUnshadow(child)
			}
		}

	case *parse.ActionNode:
		t.browseToUnshadow(node.Pipe)

	case *parse.PipeNode:
		// note that it matters to handle cmds before declarations
		for _, child := range node.Cmds {
			t.browseToUnshadow(child)
		}
		for _, child := range node.Decl {
			t.handleDecl(child)
		}

	case *parse.CommandNode:
		for _, child := range node.Args {
			t.browseToUnshadow(child)
		}

	case *parse.RangeNode:
		t.browseToUnshadow(node.Pipe)
		t.browseToUnshadow(node.List)
		t.browseToUnshadow(node.ElseList)

	case *parse.IfNode:
		t.browseToUnshadow(node.Pipe)
		t.browseToUnshadow(node.List)
		t.browseToUnshadow(node.ElseList)

	case *parse.WithNode:
		t.browseToUnshadow(node.Pipe)
		t.browseToUnshadow(node.List)
		t.browseToUnshadow(node.ElseList)

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
	case *parse.FieldNode:
		// pass
	case *parse.VariableNode:
		// do the rename
		node.Ident[0] = t.getName(node.Ident[0])

	default:
		err := fmt.Errorf("treeUnshadower.browseToUnshadow: unhandled node type\n%v\n%#v", node, node)
		panic(err)
	}
}

// handleDecl should receive a variable node of a Pipe.Decl,
// it renames the variable declaration to unshadow it from previous declartions.
// if the variable does not need unshadowing, it is left intact.
func (t *treeUnshadower) handleDecl(node *parse.VariableNode) {
	// is it {{$v := ""}} or {{$v.field = "rr"}} ?
	if len(node.Ident) == 1 {
		// Is this var $v already defined ?
		if t.hasVar(node.Ident[0]) == false {
			// nop, then register it only.
			t.registerVar(node.Ident[0])
		} else {
			// yup, then rename it.
			node.Ident[0] = t.rename(node.Ident[0])
		}
	} else {
		// keep browsing.
		t.browseToUnshadow(node)
	}
}
