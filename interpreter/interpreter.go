package interpreter

import (
	"fmt"
	"goc/parser"
	"goc/token"
	"strconv"
)

var values = make(map[string]*Value)

type Interpreter struct {
	parser *parser.Parser
}

func New(p *parser.Parser) *Interpreter {
	return &Interpreter{
		parser: p,
	}
}

func (i *Interpreter) Interpret() *Value {
	tree := i.parser.Parse()
	return i.visit(tree)
}

func (i *Interpreter) PrintValues() {
	for k, v := range values {
		fmt.Printf("%v: type: %v, value: %+v\n", k, v.Type, v.Value)
	}
}

func (i *Interpreter) visit(node interface{}) *Value {
	switch node := node.(type) {
	case parser.BinOp:
		return i.visitBinOp(node)

	case parser.UnaryOp:
		return i.visitUnaryOp(node)

	case parser.Num:
		return i.visitNum(node)

	case parser.Assign:
		return i.visitAssign(node)

	case parser.Compound:
		return i.visitCompound(node)

	case parser.Var:
		return i.visitVar(node)

	case parser.NoOp:
		return i.visitNoOp(node)
	}

	panic(fmt.Sprintf("unexpected node %v", node))
}

func (i *Interpreter) visitBinOp(node parser.BinOp) *Value {
	l := i.visit(node.Left)
	r := i.visit(node.Right)
	return cal(l, r, node.Token)
}

func (i *Interpreter) visitUnaryOp(node parser.UnaryOp) *Value {
	expr := i.visit(node.Expr)
	if node.Token.Type == token.TOKEN_TYPE_OP_SUB {
		switch expr.Type {
		case VALUE_TYPE_INTEGER:
			expr.Value = -expr.Value.(int64)

		case VALUE_TYPE_FLOAT:
			expr.Value = -expr.Value.(float64)
		}
	}

	return expr
}

func (i *Interpreter) visitNum(node parser.Num) *Value {
	if node.Token.Type == token.TOKEN_TYPE_INTEGER {
		return &Value{
			Type:  VALUE_TYPE_INTEGER,
			Value: parseInt64(node.Token.Value),
		}
	} else if node.Token.Type == token.TOKEN_TYPE_FLOAT {
		f, _ := strconv.ParseFloat(node.Token.Value, 64)

		return &Value{
			Type:  VALUE_TYPE_FLOAT,
			Value: f,
		}
	}

	panic("unrecognize Num node")
}

func (i *Interpreter) visitCompound(node parser.Compound) *Value {
	for _, child := range node.Children {
		i.visit(child)
	}

	return nil
}

func (i *Interpreter) visitVar(node parser.Var) *Value {
	v, exists := values[node.Value.(string)]
	if !exists {
		panic(fmt.Sprintf("undefined variable %v", v))
	}

	return v
}

func (i *Interpreter) visitLeftVariable(node parser.Var) *Value {
	return &Value{
		Type:  VALUE_TYPE_UNKNOWN,
		Value: node.Value,
	}
}

func (i *Interpreter) visitAssign(node parser.Assign) *Value {
	leftVar, ok := node.Left.(parser.Var)
	if !ok {
		panic(fmt.Sprintf("expected node Var, but got %#v", node.Left))
	}
	left := i.visitLeftVariable(leftVar)
	right := i.visit(node.Right)
	values[left.Value.(string)] = right
	return nil
}

func (i *Interpreter) visitNoOp(node parser.NoOp) *Value {
	return nil
}
