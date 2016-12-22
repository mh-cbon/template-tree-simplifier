## cli

cli is an useful cli tools for development.

It needs to `go get github.com/mh-cbon/print-template-tree`

Adjust the template at `cli/tpl/test.tpl`, adjust `funcs` and `data` within the `main.go`, run it.

It will pretty print information about the template before / after its processing, such as

```sh
cli/tpl/test.tpl
------------------
Tree.Name=test.tpl

BEFORE TRANSFORM
TEMPLATE CONTENT
{{if eq ("what" | lower) ("what" | lower)}}{{end}}

TEMPLATE EXECUTION
"\n"
TEMPLATE TREE
 parse.Tree
   parse.IfNode
    parse.PipeNode
     parse.CommandNode
      parse.IdentifierNode Ident="eq"
      parse.PipeNode
       parse.CommandNode
        parse.StringNode Quoted="\"what\"" Text="what"
       parse.CommandNode
        parse.IdentifierNode Ident="lower"
      parse.PipeNode
       parse.CommandNode
        parse.StringNode Quoted="\"what\"" Text="what"
       parse.CommandNode
        parse.IdentifierNode Ident="lower"
   parse.TextNode Text="\n"

AFTER TRANSFORM
TEMPLATE CONTENT
{{$var1 := lower "what"}}{{$var2 := lower "what"}}{{$var0 := eq $var1 $var2}}{{if $var0}}{{end}}

TEMPLATE EXECUTION
"\n"
TEMPLATE TREE
 parse.Tree
   parse.ActionNode
    parse.PipeNode
     parse.VariableNode Ident=["$var1"]
     parse.CommandNode
      parse.IdentifierNode Ident="lower"
      parse.StringNode Quoted="\"what\"" Text="what"
   parse.ActionNode
    parse.PipeNode
     parse.VariableNode Ident=["$var2"]
     parse.CommandNode
      parse.IdentifierNode Ident="lower"
      parse.StringNode Quoted="\"what\"" Text="what"
   parse.ActionNode
    parse.PipeNode
     parse.VariableNode Ident=["$var0"]
     parse.CommandNode
      parse.IdentifierNode Ident="eq"
      parse.VariableNode Ident=["$var1"]
      parse.VariableNode Ident=["$var2"]
   parse.IfNode
    parse.PipeNode
     parse.CommandNode
      parse.VariableNode Ident=["$var0"]
   parse.TextNode Text="\n"
```
