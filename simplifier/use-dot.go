package simplifier

import (
	"fmt"
	"text/template/parse"
)

//IsUsingDot tells if a template tree consumes a dot value.
func IsUsingDot(tree *parse.Tree) bool {
	return browseNodesToCheckIfDotIsUsed(tree.Root)
}

// browseNodesToCheckIfDotIsUsed browses all tree nodes,
// returns true if a node using a dot {{.}} exists.
func browseNodesToCheckIfDotIsUsed(l interface{}) bool {
	switch node := l.(type) {

	case *parse.ListNode:
		if node != nil {
			for _, child := range node.Nodes {
				if browseNodesToCheckIfDotIsUsed(child) {
					return true
				}
			}
		}

	case *parse.ActionNode:
		if browseNodesToCheckIfDotIsUsed(node.Pipe) {
			return true
		}

	case *parse.PipeNode:
		for _, child := range node.Decl {
			if browseNodesToCheckIfDotIsUsed(child) {
				return true
			}
		}
		for _, child := range node.Cmds {
			if browseNodesToCheckIfDotIsUsed(child) {
				return true
			}
		}

	case *parse.CommandNode:
		for _, child := range node.Args {
			if browseNodesToCheckIfDotIsUsed(child) {
				return true
			}
		}

	case *parse.RangeNode:
		if browseNodesToCheckIfDotIsUsed(node.Pipe) {
			return true
		}
		// it is not needed to enter into range node,
		// the dot of a range body does not equal the root dot,
		// unless root dot is consumed into the range Pipe
		// if browseNodesToCheckIfDotIsUsed(node.List) {
		// 	return true
		// }
		// if browseNodesToCheckIfDotIsUsed(node.ElseList) {
		// 	return true
		// }

	case *parse.IfNode:
		if browseNodesToCheckIfDotIsUsed(node.Pipe) {
			return true
		}
		if browseNodesToCheckIfDotIsUsed(node.List) {
			return true
		}
		if browseNodesToCheckIfDotIsUsed(node.ElseList) {
			return true
		}

	case *parse.WithNode:
		if browseNodesToCheckIfDotIsUsed(node.BranchNode.Pipe) {
			return true
		}

		// it is not needed to enter into with node,
		// the dot of a range body does not equal the root dot,
		// unless root dot is consumed into the with Pipe
		// if browseNodesToCheckIfDotIsUsed(node.BranchNode.List) {
		// 	return true
		// }
		// if browseNodesToCheckIfDotIsUsed(node.BranchNode.ElseList) {
		// 	return true
		// }

	case *parse.TemplateNode:
		if node.Pipe != nil {
			if browseNodesToCheckIfDotIsUsed(node.Pipe) {
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
		return true // easy one
	case *parse.FieldNode:
		return true // easy one
	case *parse.TextNode:
		// pass

	default:
		err := fmt.Errorf("browseNodesToCheckIfDotIsUsed: unhandled node type\n%v\n%#v", node, node)
		panic(err)
	}
	return false
}
