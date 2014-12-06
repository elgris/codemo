package models

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

const FuncNameMain = "main"

var NoMainErr = errors.New("could not find func main to start parse from")

type (
	Node interface {
		Type() string
	}

	Expr interface {
		Node
		exprNode()
	}

	baseNode struct {
		NodeType string `json:"type"`
	}

	baseExpr struct {
		baseNode
	}

	Block struct {
		baseNode
		Children []Node `json:"children"`
	}

	Assignment struct {
		baseNode
		Left  Expr `json:"left"`
		Right Expr `json:"right"`
	}

	Loop struct {
		baseNode
	}

	LoopRange struct {
		baseNode
	}

	Cond struct {
		baseNode
	}

	Arithm struct {
		baseExpr
		Left  Expr `json:"left"`
		Right Expr `json:"right"`
	}

	Call struct {
		baseExpr
		Func string `json:"func"`
		Args []Expr `json:"args"`
	}

	Const struct {
		baseExpr
		Value string `json:"value"`
	}

	Var struct {
		baseExpr
		Name string `json:"name"`
	}

	Array struct {
		baseExpr
		Items []Expr `json:"items"`
	}
)

func newBlock() *Block           { return &Block{baseNode: baseNode{NodeType: "block"}} }
func newAssignment() *Assignment { return &Assignment{baseNode: baseNode{NodeType: "assignment"}} }
func newLoop() *Loop             { return &Loop{baseNode: baseNode{NodeType: "loop"}} }
func newLoopRange() *LoopRange   { return &LoopRange{baseNode: baseNode{NodeType: "loopRange"}} }
func newCond() *Cond             { return &Cond{baseNode: baseNode{NodeType: "cond"}} }

func newArithm() *Arithm { return &Arithm{baseExpr: baseExpr{baseNode: baseNode{NodeType: "arithm"}}} }
func newCall() *Call     { return &Call{baseExpr: baseExpr{baseNode: baseNode{NodeType: "call"}}} }
func newConst() *Const   { return &Const{baseExpr: baseExpr{baseNode: baseNode{NodeType: "const"}}} }
func newVar() *Var       { return &Var{baseExpr: baseExpr{baseNode: baseNode{NodeType: "var"}}} }
func newArray() *Array   { return &Array{baseExpr: baseExpr{baseNode: baseNode{NodeType: "array"}}} }

func (n *baseNode) Type() string { return n.NodeType }

func (n *baseExpr) exprNode() {}

func ParseSrc(src string) (Node, error) {
	fset := token.NewFileSet()

	fullSrc := fmt.Sprintf("package codemotst\n func main() {\n %s \n }", src)

	f, err := parser.ParseFile(fset, "", fullSrc, 0)
	if err != nil {
		return nil, err
	}

	// var bf bytes.Buffer
	// ast.Fprint(&bf, fset, f, ast.NotNilFilter)

	// fmt.Println(bf.String())

	return parseAst(f)
}

func parseAst(f *ast.File) (node Node, err error) {
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == FuncNameMain {
			return parseBlock(funcDecl.Body)
			// TODO: parse contents of Main
			// return parsed result
		}
	}

	return nil, NoMainErr
}

func parseStmt(stmt ast.Stmt) (Node, error) {
	switch stmt.(type) {
	case *ast.BlockStmt:
		return parseBlock(stmt.(*ast.BlockStmt))
	case *ast.AssignStmt:
		return parseAssign(stmt.(*ast.AssignStmt))
	case *ast.ExprStmt:
		return parseExpr(stmt.(*ast.ExprStmt).X)
	case *ast.RangeStmt:
		// TODO parse RangeLoop
		return nil, nil
	case *ast.ForStmt:
		// TODO parse ForLoop
		return nil, nil
	case *ast.IfStmt:
		// TODO parse RangeLoop
		return nil, nil
	case *ast.IncDecStmt:
		return nil, nil
	case *ast.DeclStmt:
		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf("still unsupported AST type: %T", stmt))
	}
}

func parseBlock(block *ast.BlockStmt) (node *Block, err error) {
	node = newBlock()
	var parsedStmt Node

	for _, stmt := range block.List {
		parsedStmt, err = parseStmt(stmt)
		if err != nil {
			return
		}
		node.Children = append(node.Children, parsedStmt)
	}

	return
}

func parseAssign(stmt *ast.AssignStmt) (node *Assignment, err error) {
	if len(stmt.Lhs) != 1 || len(stmt.Rhs) != 1 {
		return nil, errors.New("right now we can process assignments with only 1 receiver")
	}
	node = newAssignment()

	node.Left, err = parseExpr(stmt.Lhs[0])
	if err != nil {
		return nil, err
	}
	node.Right, err = parseExpr(stmt.Rhs[0])
	if err != nil {
		return nil, err
	}
	return node, nil
}

func parseExpr(expr ast.Expr) (Expr, error) {
	switch expr.(type) {
	case *ast.BasicLit:
		i := newConst()
		i.Value = expr.(*ast.BasicLit).Value

		return i, nil
	case *ast.Ident:
		i := newVar()
		i.Name = expr.(*ast.Ident).Name

		return i, nil
	case *ast.CallExpr:
		return parseCall(expr.(*ast.CallExpr))
	case *ast.CompositeLit:
		return parseComposite(expr.(*ast.CompositeLit))
	case *ast.BinaryExpr:
		return nil, nil
		//TODO
	default:
		return nil, errors.New(fmt.Sprintf("still unsupported AST expression: %T", expr))
	}

	return nil, nil
}

func parseExprs(exprs []ast.Expr) (is []Expr, err error) {
	is = make([]Expr, len(exprs))

	for ind, expr := range exprs {
		is[ind], err = parseExpr(expr)
		if err != nil {
			return nil, err
		}
	}

	return
}

func parseCall(expr *ast.CallExpr) (c *Call, err error) {
	fIdent, ok := expr.Fun.(*ast.Ident)

	if !ok {
		return nil, errors.New(fmt.Sprintf("only simple function calls are accepted. Given: %#v", expr.Fun))
	}

	c = newCall()
	c.Func = fIdent.Name
	c.Args, err = parseExprs(expr.Args)

	// switch fIndent.Name {
	// case "len":
	// 	c = Call{}
	// 	c.Func = fIndent.Name
	// 	c.Args = parseExprs(expr.Args)
	// default:
	// 	return nil, errors.New(fmt.Sprintf("still unsupported func call: %s", fIndent.Name))
	// }

	return
}

func parseComposite(expr *ast.CompositeLit) (c *Array, err error) {
	if _, ok := expr.Type.(*ast.ArrayType); !ok {
		return nil, errors.New(fmt.Sprintf("from composite types we process only arrays. Given: %#v", expr.Type))
	}

	c = newArray()
	c.Items = make([]Expr, len(expr.Elts))

	for i, e := range expr.Elts {
		c.Items[i], err = parseExpr(e)
		if err != nil {
			return nil, err
		}
	}

	return
}
