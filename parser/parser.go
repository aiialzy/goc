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

func (p *Parser) Parse() interface{} {
	return p.program()
}

func (p *Parser) readNextToken(tokenType token.TokenType) {
	if p.currentToken.Type == tokenType {
		p.currentToken = p.lexer.NextToken()
	} else {
		panic(fmt.Sprintf("syntax error, expected %v, but got %v", tokenType, p.currentToken.Type))
	}
}

func (p *Parser) factor() interface{} {
	if p.currentToken.Type == token.TOKEN_TYPE_OP_SUB ||
		p.currentToken.Type == token.TOKEN_TYPE_OP_ADD {
		t := p.currentToken
		p.readNextToken(t.Type)
		return UnaryOp{
			Token: t,
			Expr:  p.factor(),
		}
	} else if p.currentToken.Type == token.TOKEN_TYPE_INTEGER {
		t := p.currentToken
		p.readNextToken(token.TOKEN_TYPE_INTEGER)
		return Num{Token: t}
	} else if p.currentToken.Type == token.TOKEN_TYPE_FLOAT {
		t := p.currentToken
		p.readNextToken(token.TOKEN_TYPE_FLOAT)
		return Num{Token: t}
	} else if p.currentToken.Type == token.TOKEN_TYPE_LPAREN {
		p.readNextToken(token.TOKEN_TYPE_LPAREN)
		result := p.expr()
		p.readNextToken(token.TOKEN_TYPE_RPAREN)
		return result
	} else if p.currentToken.Type == token.TOKEN_TYPE_ID {
		t := p.currentToken
		p.readNextToken(token.TOKEN_TYPE_ID)
		return Var{
			Token: t,
			Value: t.Value,
		}
	} else if p.currentToken.Type == token.TOKEN_TYPE_EOF {
		return nil
	}

	panic(fmt.Sprintf("unexpected token %v", p.currentToken))
}

func (p *Parser) term() interface{} {

	l := p.factor()
	if l == nil {
		return l
	}

	ops := map[token.TokenType]struct{}{
		token.TOKEN_TYPE_OP_MUL: {},
		token.TOKEN_TYPE_OP_DIV: {},
		token.TOKEN_TYPE_OP_MOD: {},
	}

	for {
		if _, exists := ops[p.currentToken.Type]; !exists {
			break
		}

		t := p.currentToken
		if t.Type == token.TOKEN_TYPE_OP_MUL {
			p.readNextToken(token.TOKEN_TYPE_OP_MUL)
		} else if t.Type == token.TOKEN_TYPE_OP_DIV {
			p.readNextToken(token.TOKEN_TYPE_OP_DIV)
		} else if t.Type == token.TOKEN_TYPE_OP_MOD {
			p.readNextToken(token.TOKEN_TYPE_OP_MOD)
		}
		r := p.factor()
		l = BinOp{
			Left:  l,
			Right: r,
			Token: t,
		}
	}

	return l
}

func (p *Parser) expr() interface{} {

	l := p.term()

	ops := map[token.TokenType]struct{}{
		token.TOKEN_TYPE_OP_ADD: {},
		token.TOKEN_TYPE_OP_SUB: {},
	}

	for {
		if _, exists := ops[p.currentToken.Type]; !exists {
			break
		}

		t := p.currentToken
		if t.Type == token.TOKEN_TYPE_OP_ADD {
			p.readNextToken(token.TOKEN_TYPE_OP_ADD)
		} else if t.Type == token.TOKEN_TYPE_OP_SUB {
			p.readNextToken(token.TOKEN_TYPE_OP_SUB)
		}
		r := p.term()
		l = BinOp{
			Left:  l,
			Right: r,
			Token: t,
		}
	}

	return l
}

func (p *Parser) assignmentPart() interface{} {
	id := p.currentToken
	p.readNextToken(token.TOKEN_TYPE_ID)
	assign := p.currentToken
	p.readNextToken(token.TOKEN_TYPE_ASSIGN)
	expr := p.expr()
	p.readNextToken(token.TOKEN_TYPE_SEMI)

	return Assign{
		Left: Var{
			Token: id,
			Value: id.Value,
		},
		Right: expr,
		Token: assign,
	}
}

func (p *Parser) assignmentStatement() interface{} {
	p.readNextToken(token.TOKEN_TYPE_VAR)
	if p.currentToken.Type == token.TOKEN_TYPE_LPAREN {
		p.readNextToken(token.TOKEN_TYPE_LPAREN)
		t := p.assignCompoundStatement()
		p.readNextToken(token.TOKEN_TYPE_RPAREN)
		return t
	} else {
		return p.assignmentPart()
	}
}

func (p *Parser) statement() interface{} {
	switch p.currentToken.Type {
	case token.TOKEN_TYPE_VAR:
		return p.assignmentStatement()
	case token.TOKEN_TYPE_LBRACE:
		return p.compoundStatement()
	case token.TOKEN_TYPE_SEMI:
		return p.empty()
	}
	panic(fmt.Sprintf("syntax error, unexpected token %v", p.currentToken))
}

func (p *Parser) empty() interface{} {
	p.readNextToken(token.TOKEN_TYPE_SEMI)
	return NoOp{}
}

func (p *Parser) statementList() []interface{} {
	var statementList []interface{}
	for p.currentToken.Type != token.TOKEN_TYPE_RBRACE {
		statementList = append(statementList, p.statement())
	}

	return statementList
}

func (p *Parser) assignCompoundStatement() interface{} {
	var t Compound
	for p.currentToken.Type != token.TOKEN_TYPE_RPAREN {
		t.Children = append(t.Children, p.assignmentPart())
	}

	return t
}

func (p *Parser) compoundStatement() interface{} {
	var t Compound
	p.readNextToken(token.TOKEN_TYPE_LBRACE)
	t.Children = append(t.Children, p.statementList()...)
	p.readNextToken(token.TOKEN_TYPE_RBRACE)
	return t
}

func (p *Parser) program() interface{} {
	return p.compoundStatement()
}
