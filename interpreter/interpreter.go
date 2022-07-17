package interpreter

import (
	"fmt"
	"goc/parser"
	"goc/token"
	"strconv"
	"strings"
)

type ValueType int

const (
	VALUE_TYPE_UNKNOWN ValueType = iota
	VALUE_TYPE_INTEGER
	VALUE_TYPE_FLOAT
)

var VALUE_TYPE_MAP = map[ValueType]string{
	VALUE_TYPE_UNKNOWN: "VALUE_TYPE_UNKNOWN",
	VALUE_TYPE_INTEGER: "VALUE_TYPE_INTEGER",
	VALUE_TYPE_FLOAT:   "VALUE_TYPE_FLOAT",
}

func (vt ValueType) String() string {
	return VALUE_TYPE_MAP[vt]
}

type Interpreter struct {
	parser *parser.Parser
}

func New(p *parser.Parser) *Interpreter {
	return &Interpreter{
		parser: p,
	}
}

func (i *Interpreter) Interpret() (string, ValueType) {
	tree := i.parser.Expr()
	return i.visit(tree)
}

func (i *Interpreter) visit(node interface{}) (string, ValueType) {
	switch node := node.(type) {
	case parser.BinOp:
		return i.visitBinOp(node)

	case parser.UnaryOp:
		return i.visitUnaryOp(node)

	case parser.Num:
		return i.visitNum(node)
	}

	panic(fmt.Sprintf("unexpected node %v", node))
}

func (i *Interpreter) visitBinOp(node parser.BinOp) (string, ValueType) {
	l, lValueType := i.visit(node.Left)
	r, rValueType := i.visit(node.Right)
	return cal(l, r, lValueType, rValueType, node.Token)
}

func (i *Interpreter) visitUnaryOp(node parser.UnaryOp) (string, ValueType) {
	expr, valueType := i.visit(node.Expr)
	if node.Token.Type == token.TOKEN_TYPE_OP_SUB {
		if strings.HasPrefix(expr, "-") {
			expr = expr[1:]
		} else {
			expr = "-" + expr
		}
	}

	return expr, valueType
}

func (i *Interpreter) visitNum(node parser.Num) (string, ValueType) {
	return node.Token.Value, getValueType(node.Token.Value)
}

func parseInt64(numStr string) int64 {
	var ret int64
	var err error
	if strings.HasPrefix(numStr, "0x") {
		ret, err = strconv.ParseInt(numStr[2:], 16, 64)
	} else if strings.HasPrefix(numStr, "0o") {
		ret, err = strconv.ParseInt(numStr[2:], 8, 64)
	} else if strings.HasPrefix(numStr, "0b") {
		ret, err = strconv.ParseInt(numStr[2:], 2, 64)
	} else if len(numStr) > 1 && strings.HasPrefix(numStr, "0") {
		ret, err = strconv.ParseInt(numStr[1:], 8, 64)
	} else {
		ret, err = strconv.ParseInt(numStr, 10, 64)
	}

	if err != nil {
		panic(fmt.Sprintf("invalid integer %v, err: %v", numStr, err))
	}

	return ret
}

func getValueType(v string) ValueType {
	if strings.Contains(v, ".") {
		return VALUE_TYPE_FLOAT
	}
	return VALUE_TYPE_INTEGER
}

func cal(l, r string, lValueType, rValueType ValueType, op *token.Token) (string, ValueType) {
	var iresult int64
	var fresult float64

	fValueType := VALUE_TYPE_INTEGER
	if lValueType == VALUE_TYPE_FLOAT || rValueType == VALUE_TYPE_FLOAT {
		fValueType = VALUE_TYPE_FLOAT
	}

	if fValueType == VALUE_TYPE_INTEGER {
		li := parseInt64(l)
		ri := parseInt64(r)
		switch op.Type {
		case token.TOKEN_TYPE_OP_ADD:
			iresult = li + ri

		case token.TOKEN_TYPE_OP_SUB:
			iresult = li - ri

		case token.TOKEN_TYPE_OP_MUL:
			iresult = li * ri

		case token.TOKEN_TYPE_OP_DIV:
			iresult = li / ri

		case token.TOKEN_TYPE_OP_MOD:
			iresult = li % ri
		}

		return strconv.Itoa(int(iresult)), VALUE_TYPE_INTEGER
	} else {
		li, _ := strconv.ParseFloat(l, 64)
		ri, _ := strconv.ParseFloat(r, 64)
		switch op.Type {
		case token.TOKEN_TYPE_OP_ADD:
			fresult = li + ri

		case token.TOKEN_TYPE_OP_SUB:
			fresult = li - ri

		case token.TOKEN_TYPE_OP_MUL:
			fresult = li * ri

		case token.TOKEN_TYPE_OP_DIV:
			fresult = li / ri
		}
		return strconv.FormatFloat(fresult, 'f', -1, 64), VALUE_TYPE_FLOAT
	}

}
