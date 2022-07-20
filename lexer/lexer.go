package lexer

import (
	"fmt"
	"goc/token"
	"io"
	"strings"
	"unicode"
)

const MAX_TEMP_RUNES_SIZE = 8

var RESERVED_KEYWORDS = map[string]*token.Token{
	"var": {Type: token.TOKEN_TYPE_VAR, Value: "var"},
}

func New(source string) *Lexer {
	return &Lexer{
		reader: &runeReader{
			runes: []rune(source),
			index: 0,
		},
	}
}

type Lexer struct {
	reader *runeReader
}

func (l *Lexer) nextChar() (ch rune, err error) {
	ch, _, err = l.reader.ReadRune()
	return
}

func (l *Lexer) peek(t int) (ch rune, err error) {
	return l.reader.PeekRune(t)
}

type checkFunc func(int, rune) bool

func (l *Lexer) read(check checkFunc) (string, error) {
	var builder strings.Builder
	index := 0
	for {
		ch, err := l.nextChar()
		if err == io.EOF {
			return builder.String(), nil
		} else if err != nil {
			return builder.String(), err
		}

		if check(index, ch) {
			builder.WriteRune(ch)
			index++
		} else {
			l.reader.UnreadRune()
			return builder.String(), nil
		}
	}
}

func (l *Lexer) readNumber() (string, error) {
	return l.read(
		func(i int, r rune) bool {
			return unicode.IsDigit(r) || r == '.' || r == 'x' || r == 'b' || r == 'o'
		},
	)
}

func (l *Lexer) readWord() (string, error) {
	return l.read(func(i int, r rune) bool {
		return unicode.IsLetter(r) || r == '_' || r > 127 || unicode.IsNumber(r)
	})
}

func (l *Lexer) readSpace() (string, error) {
	return l.read(func(i int, r rune) bool {
		return unicode.IsSpace(r)
	})
}

func (l *Lexer) readCommentLine() (string, error) {
	return l.read(func(i int, r rune) bool {
		return r != '\n' && r != '\r'
	})
}

func (l *Lexer) readCommentLines() (string, error) {
	var last rune
	var stop bool
	return l.read(func(i int, r rune) bool {
		if stop {
			return false
		}

		if last == '*' && r == '/' {
			stop = true
		}
		last = r
		return true
	})
}

func (l *Lexer) readOp(op string) string {
	rs := []rune(op)
	op, _ = l.read(func(i int, r rune) bool {
		if i < len(op) && rs[i] == r {
			i++
			return true
		}
		return false
	})

	return op
}

func (l *Lexer) NextToken() *token.Token {
	for {
		ch, err := l.peek(0)
		if err == io.EOF {
			return token.New(token.TOKEN_TYPE_EOF, "EOF")
		} else if err != nil {
			panic(fmt.Sprintf("read next token failed: %v\n", err))
		}

		if unicode.IsSpace(ch) {
			_, err := l.readSpace()
			if err != nil {
				panic(fmt.Sprintf("read space failed: %v\n", err))
			}
			continue
		}

		if ch == '/' {
			ch, err := l.peek(1)
			if err != nil && err != io.EOF {
				panic(fmt.Sprintf("peek failed: %v\n", err))
			} else if ch == '/' {
				_, err := l.readCommentLine()
				if err != nil && err != io.EOF {
					panic(fmt.Sprintf("read comment line failed: %v\n", err))
				}
				continue
			} else if ch == '*' {
				_, err := l.readCommentLines()
				if err != nil && err != io.EOF {
					panic(fmt.Sprintf("read comment lines failed: %v\n", err))
				}
				continue
			}
		}

		if unicode.IsNumber(ch) {
			numStr, err := l.readNumber()
			if err != nil {
				fmt.Printf("read number failed: %v\n", err)
			}
			if strings.Contains(numStr, ".") {
				return token.New(token.TOKEN_TYPE_FLOAT, numStr)
			}
			return token.New(token.TOKEN_TYPE_INTEGER, numStr)
		}

		if unicode.IsLetter(ch) || ch == '_' || ch > 127 {
			word, err := l.readWord()
			if err != nil {
				fmt.Printf("read word failed: %v\n", err)
			}
			if t, exists := RESERVED_KEYWORDS[word]; exists {
				return t
			}
			return token.New(token.TOKEN_TYPE_ID, word)
		}

		switch ch {
		case '+':
			return token.New(token.TOKEN_TYPE_OP_ADD, l.readOp("+"))
		case '-':
			return token.New(token.TOKEN_TYPE_OP_SUB, l.readOp("-"))
		case '*':
			return token.New(token.TOKEN_TYPE_OP_MUL, l.readOp("*"))
		case '/':
			return token.New(token.TOKEN_TYPE_OP_MUL, l.readOp("/"))
		case '%':
			return token.New(token.TOKEN_TYPE_OP_MOD, l.readOp("%"))
		case '(':
			return token.New(token.TOKEN_TYPE_LPAREN, l.readOp("("))
		case ')':
			return token.New(token.TOKEN_TYPE_RPAREN, l.readOp(")"))
		case '{':
			return token.New(token.TOKEN_TYPE_LBRACE, l.readOp("{"))
		case '}':
			return token.New(token.TOKEN_TYPE_RBRACE, l.readOp("}"))
		case ';':
			return token.New(token.TOKEN_TYPE_SEMI, l.readOp(";"))
		case '=':
			return token.New(token.TOKEN_TYPE_ASSIGN, l.readOp("="))
		}

		panic(fmt.Sprintf("unexpected character: %c", ch))
	}

}
