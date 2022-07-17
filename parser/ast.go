package parser

import "goc/token"

type Node interface {
	Num | BinOp
}

type AST struct {
}

type Num struct {
	AST
	Token *token.Token
}

type BinOp struct {
	AST
	Left  interface{}
	Token *token.Token
	Right interface{}
}

type UnaryOp struct {
	AST
	Token *token.Token
	Expr  interface{}
}
