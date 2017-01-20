package simplifier

import (
	"fmt"
	"text/template/parse"
)

//PrintsAnything tells if a template tree prints anything.
func PrintsAnything(tree *parse.Tree) bool {
	return browseNodesToCheckIfItPrintsAnything(tree.Root)
}

// browseNodesToCheckIfItPrintsAnything browses all tree nodes,
// returns true if a node using a dot {{.}} exists.
func browseNodesToCheckIfItPrintsAnything(l interface{}) bool {
	switch node := l.(type) {

	case *parse.ListNode:
		if node != nil {
			for _, child := range node.Nodes {
				if browseNodesToCheckIfItPrintsAnything(child) {
					return true
				}
			}
		}

	case *parse.ActionNode:
		if len(node.Pipe.Decl) == 0 && len(node.Pipe.Cmds) > 0 {
			return true
		}
		if browseNodesToCheckIfItPrintsAnything(node.Pipe) {
			return true
		}

	case *parse.PipeNode:
		for _, child := range node.Decl {
			if browseNodesToCheckIfItPrintsAnything(child) {
				return true
			}
		}
		for _, child := range node.Cmds {
			if browseNodesToCheckIfItPrintsAnything(child) {
				return true
			}
		}

	case *parse.CommandNode:
		for _, child := range node.Args {
			if browseNodesToCheckIfItPrintsAnything(child) {
				return true
			}
		}

	case *parse.RangeNode:
		if browseNodesToCheckIfItPrintsAnything(node.Pipe) {
			return true
		}
		if browseNodesToCheckIfItPrintsAnything(node.List) {
			return true
		}
		if browseNodesToCheckIfItPrintsAnything(node.ElseList) {
			return true
		}

	case *parse.IfNode:
		if browseNodesToCheckIfItPrintsAnything(node.Pipe) {
			return true
		}
		if browseNodesToCheckIfItPrintsAnything(node.List) {
			return true
		}
		if browseNodesToCheckIfItPrintsAnything(node.ElseList) {
			return true
		}

	case *parse.WithNode:
		if browseNodesToCheckIfItPrintsAnything(node.BranchNode.Pipe) {
			return true
		}
		if browseNodesToCheckIfItPrintsAnything(node.BranchNode.List) {
			return true
		}
		if browseNodesToCheckIfItPrintsAnything(node.BranchNode.ElseList) {
			return true
		}

	case *parse.TemplateNode:
		if node.Pipe != nil {
			if browseNodesToCheckIfItPrintsAnything(node.Pipe) {
				return true
			}
		}

	case *parse.VariableNode:
		// -
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
		return true

	default:
		err := fmt.Errorf("browseNodesToCheckIfItPrintsAnything: unhandled node type\n%v\n%#v", node, node)
		panic(err)
	}
	return false
}
