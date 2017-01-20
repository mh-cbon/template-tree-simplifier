package simplifier

import (
	"fmt"
	"text/template/parse"

	"github.com/serenize/snaker"
)

// renameVariables browses all tree nodes,
// renames all encountered VariableNode Ident,
// adding a prefix "tpl_".
func renameVariables(l interface{}) {
	switch node := l.(type) {

	case *parse.ListNode:
		if node != nil {
			for _, child := range node.Nodes {
				renameVariables(child)
			}
		}

	case *parse.VariableNode:
		if node.Ident[0] != "$" {
			node.Ident[0] = "$tpl" + snaker.SnakeToCamel(node.Ident[0][1:]) // get ride of $ sign
		}

	case *parse.ActionNode:
		renameVariables(node.Pipe)

	case *parse.PipeNode:
		for _, child := range node.Decl {
			renameVariables(child)
		}
		for _, child := range node.Cmds {
			renameVariables(child)
		}

	case *parse.CommandNode:
		for _, child := range node.Args {
			renameVariables(child)
		}

	case *parse.RangeNode:
		renameVariables(node.Pipe)
		renameVariables(node.List)
		renameVariables(node.ElseList)

	case *parse.IfNode:
		renameVariables(node.Pipe)
		renameVariables(node.List)
		renameVariables(node.ElseList)

	case *parse.WithNode:
		renameVariables(node.BranchNode.Pipe)
		renameVariables(node.BranchNode.List)
		renameVariables(node.BranchNode.ElseList)

	case *parse.TemplateNode:
		if node.Pipe != nil {
			renameVariables(node.Pipe)
		}

	case *parse.StringNode:
		// pass
	case *parse.NumberNode:
		// pass
	case *parse.BoolNode:
		// pass
	case *parse.IdentifierNode:
		// pass
	case *parse.DotNode:
		// pass
	case *parse.FieldNode:
		// pass
	case *parse.TextNode:
		// pass

	default:
		err := fmt.Errorf("renameVariables: unhandled node type\n%v\n%#v", node, node)
		panic(err)
	}
}
