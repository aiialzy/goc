package parser

import (
	"fmt"
	"goc/lexer"
	"goc/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken *token.Token
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		lexer:        l,
		currentToken: l.NextToken(),
	}
}

func (i *Parser) readNextToken(tokenType token.TokenType) {
	if i.currentToken.Type == tokenType {
		i.currentToken = i.lexer.NextToken()
	} else {
		panic(fmt.Sprintf("syntax error, expected %v, but got %v", tokenType, i.currentToken.Type))
	}
}

func (i *Parser) factor() interface{} {
	if i.currentToken.Type == token.TOKEN_TYPE_OP_SUB ||
		i.currentToken.Type == token.TOKEN_TYPE_OP_ADD {
		t := i.currentToken
		i.readNextToken(t.Type)
		return UnaryOp{
			Token: t,
			Expr:  i.factor(),
		}
	} else if i.currentToken.Type == token.TOKEN_TYPE_INTEGER {
		t := i.currentToken
		i.readNextToken(token.TOKEN_TYPE_INTEGER)
		return Num{Token: t}
	} else if i.currentToken.Type == token.TOKEN_TYPE_FLOAT {
		t := i.currentToken
		i.readNextToken(token.TOKEN_TYPE_FLOAT)
		return Num{Token: t}
	} else if i.currentToken.Type == token.TOKEN_TYPE_LPAREN {
		i.readNextToken(token.TOKEN_TYPE_LPAREN)
		result := i.Expr()
		i.readNextToken(token.TOKEN_TYPE_RPAREN)
		return result
	} else if i.currentToken.Type == token.TOKEN_TYPE_EOF {
		return nil
	}

	panic(fmt.Sprintf("unexpected token %v", i.currentToken))
}

func (i *Parser) term() interface{} {

	l := i.factor()
	if l == nil {
		return l
	}

	ops := map[token.TokenType]struct{}{
		token.TOKEN_TYPE_OP_MUL: {},
		token.TOKEN_TYPE_OP_DIV: {},
		token.TOKEN_TYPE_OP_MOD: {},
	}

	for {
		if _, exists := ops[i.currentToken.Type]; !exists {
			break
		}

		t := i.currentToken
		if t.Type == token.TOKEN_TYPE_OP_MUL {
			i.readNextToken(token.TOKEN_TYPE_OP_MUL)
		} else if t.Type == token.TOKEN_TYPE_OP_DIV {
			i.readNextToken(token.TOKEN_TYPE_OP_DIV)
		} else if t.Type == token.TOKEN_TYPE_OP_MOD {
			i.readNextToken(token.TOKEN_TYPE_OP_MOD)
		}
		r := i.factor()
		l = BinOp{
			Left:  l,
			Right: r,
			Token: t,
		}
	}

	return l
}

func (i *Parser) Expr() interface{} {

	l := i.term()

	ops := map[token.TokenType]struct{}{
		token.TOKEN_TYPE_OP_ADD: {},
		token.TOKEN_TYPE_OP_SUB: {},
	}

	for {
		if _, exists := ops[i.currentToken.Type]; !exists {
			break
		}

		t := i.currentToken
		if t.Type == token.TOKEN_TYPE_OP_ADD {
			i.readNextToken(token.TOKEN_TYPE_OP_ADD)
		} else if t.Type == token.TOKEN_TYPE_OP_SUB {
			i.readNextToken(token.TOKEN_TYPE_OP_SUB)
		}
		r := i.term()
		l = BinOp{
			Left:  l,
			Right: r,
			Token: t,
		}
	}

	return l
}
