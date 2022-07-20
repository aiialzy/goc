package interpreter

import (
	"fmt"
	"goc/parser"
	"goc/token"
	"strconv"
	"strings"
)

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

func getValueType(v string) parser.ValueType {
	if strings.Contains(v, ".") {
		return parser.VALUE_TYPE_FLOAT
	}
	return parser.VALUE_TYPE_INTEGER
}

func cal(l, r *Value, op *token.Token) *Value {
	var fValue = &Value{
		Type: VALUE_TYPE_INTEGER,
	}
	if l.Type == VALUE_TYPE_FLOAT || r.Type == VALUE_TYPE_FLOAT {
		fValue.Type = VALUE_TYPE_FLOAT
	}

	var result interface{}
	if fValue.Type == VALUE_TYPE_INTEGER {
		switch op.Type {
		case token.TOKEN_TYPE_OP_ADD:
			result = l.Value.(int64) + r.Value.(int64)

		case token.TOKEN_TYPE_OP_SUB:
			result = l.Value.(int64) - r.Value.(int64)

		case token.TOKEN_TYPE_OP_MUL:
			result = l.Value.(int64) * r.Value.(int64)

		case token.TOKEN_TYPE_OP_DIV:
			result = l.Value.(int64) / r.Value.(int64)

		case token.TOKEN_TYPE_OP_MOD:
			result = l.Value.(int64) % r.Value.(int64)
		}
	} else {
		switch op.Type {
		case token.TOKEN_TYPE_OP_ADD:
			result = l.Value.(float64) + r.Value.(float64)

		case token.TOKEN_TYPE_OP_SUB:
			result = l.Value.(float64) - r.Value.(float64)

		case token.TOKEN_TYPE_OP_MUL:
			result = l.Value.(float64) * r.Value.(float64)

		case token.TOKEN_TYPE_OP_DIV:
			result = l.Value.(float64) / r.Value.(float64)
		}
	}

	fValue.Value = result

	return fValue
}
