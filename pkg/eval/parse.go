package eval

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/fugiman/tyrantbot/pkg/message"
)

const SCAFFOLD = `
package main

func main() {%s}
`

func Parse(s string) (*Code, error) {
	fileset := token.NewFileSet()
	r, err := parser.ParseFile(fileset, "!addquote", fmt.Sprintf(SCAFFOLD, s), 0)
	if err != nil {
		return nil, err
	}

	for _, id := range r.Unresolved {
		if _, ok := builtins[id.Name]; !ok {
			return nil, fmt.Errorf("Invalid identifier %q. Please ensure you use := to declare variables", id.Name)
		}
	}

	code := &Code{}
	main := r.Scope.Objects["main"].Decl.(*ast.FuncDecl).Body.List
	for _, statement := range main {
		line, err := buildExpr(statement)
		if err != nil {
			return nil, err
		}
		*code = append(*code, line)
	}
	return code, nil
}

func buildExpr(s ast.Node) (Expr, error) {
	switch v := s.(type) {
	case *ast.ExprStmt:
		return buildExpr(v.X)
	case *ast.BasicLit:
		return buildLiteral(v)
	case *ast.Ident:
		return buildVariable(v)
	case *ast.SelectorExpr:
		return buildSelector(v)
	case *ast.AssignStmt:
		return buildAssign(v)
	case *ast.CallExpr:
		return buildCall(v)
	}

	return nil, fmt.Errorf("Unknown ast.Node type: %v", s)
}

func buildLiteral(s *ast.BasicLit) (Expr, error) {
	expr := func(msg *message.Message, vars map[string]interface{}) interface{} {
		return s.Value
	}

	return expr, nil
}

func buildVariable(s *ast.Ident) (Expr, error) {
	expr := func(msg *message.Message, vars map[string]interface{}) interface{} {
		return vars[s.Name]
	}

	return expr, nil
}

func buildSelector(s *ast.SelectorExpr) (Expr, error) {
	var lookup evalFunc
	name := fmt.Sprintf("%s_%s", s.X.(*ast.Ident).Name, s.Sel.Name)

	if f, ok := builtins[name]; !ok {
		return nil, fmt.Errorf("Invalid selector: %s", name)
	} else {
		lookup = f
	}

	expr := func(msg *message.Message, vars map[string]interface{}) interface{} {
		return lookup(msg)
	}

	return expr, nil
}

func buildAssign(s *ast.AssignStmt) (Expr, error) {
	var name string
	if len(s.Lhs) != 1 {
		return nil, fmt.Errorf("Malformed variable assignment: Invalid left hand side length")
	}
	if id, ok := s.Lhs[0].(*ast.Ident); !ok {
		return nil, fmt.Errorf("Malformed variable assignment: Invalid left hand side identifier")
	} else {
		name = id.Name
	}
	if len(s.Lhs) != 1 {
		return nil, fmt.Errorf("Malformed variable assignment: Invalid right hand side length")
	}

	f, err := buildExpr(s.Rhs[0])
	if err != nil {
		return nil, fmt.Errorf("Malformed variable assignment: %v", err)
	}

	expr := func(msg *message.Message, vars map[string]interface{}) interface{} {
		vars[name] = f(msg, vars)
		return nil
	}

	return expr, nil
}

func buildCall(s *ast.CallExpr) (Expr, error) {
	name := s.Fun.(*ast.Ident).Name

	var f evalFunc
	if v, ok := builtins[name]; ok {
		f = v
	}

	args := []Expr{}
	for _, e := range s.Args {
		v, err := buildExpr(e)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}

	expr := func(msg *message.Message, vars map[string]interface{}) interface{} {
		if f == nil {
			f = vars[name].(evalFunc)
		}
		a := []interface{}{}
		for _, v := range args {
			a = append(a, v(msg, vars))
		}
		return f(a...)
	}

	return expr, nil
}
