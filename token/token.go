package token

import "fmt"

type TokenType int

const (
	TOKEN_TYPE_UNKNOWN TokenType = iota
	TOKEN_TYPE_INTEGER
	TOKEN_TYPE_FLOAT
	TOKEN_TYPE_OP_ADD
	TOKEN_TYPE_OP_SUB
	TOKEN_TYPE_OP_MUL
	TOKEN_TYPE_OP_DIV
	TOKEN_TYPE_OP_MOD
	TOKEN_TYPE_EOF
	TOKEN_TYPE_LPAREN
	TOKEN_TYPE_RPAREN
)

func (t TokenType) String() string {
	return fmt.Sprintf("%v", TOKEN_TYPE_MAP[t])
}

var TOKEN_TYPE_MAP = map[TokenType]string{
	TOKEN_TYPE_INTEGER: "INTEGER",
	TOKEN_TYPE_FLOAT:   "FLOAT",
	TOKEN_TYPE_OP_ADD:  "OP_ADD",
	TOKEN_TYPE_OP_SUB:  "OP_SUB",
	TOKEN_TYPE_OP_MUL:  "OP_MUL",
	TOKEN_TYPE_OP_DIV:  "OP_DIV",
	TOKEN_TYPE_OP_MOD:  "OP_MOD",
	TOKEN_TYPE_EOF:     "EOF",
	TOKEN_TYPE_LPAREN:  "LPAREN",
	TOKEN_TYPE_RPAREN:  "RPAREN",
}

var print_token = false

func SetPrintToken(isPrint bool) {
	print_token = isPrint
}

type Token struct {
	Type  TokenType
	Value string
}

func (t *Token) String() string {
	return fmt.Sprintf(`TOKEN(%v, "%v")`, t.Type, t.Value)
}

func New(tokenType TokenType, value string) *Token {
	t := &Token{
		Type:  tokenType,
		Value: value,
	}

	if print_token {
		fmt.Println(t)
	}

	return t
}
