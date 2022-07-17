package lexer

import (
	"bufio"
	"fmt"
	"goc/token"
	"io"
	"os"
	"strings"
	"unicode"
)

func New(ioreader io.Reader) *Lexer {
	reader := bufio.NewReader(ioreader)
	return &Lexer{
		reader: reader,
	}
}

type Lexer struct {
	reader *bufio.Reader
}

func (l *Lexer) nextChar() (ch rune, err error) {
	ch, _, err = l.reader.ReadRune()
	return
}

func (l *Lexer) peek() (ch rune, err error) {
	ch, _, err = l.reader.ReadRune()
	if err != nil {
		return
	}

	err = l.reader.UnreadRune()
	if err != nil {
		ch = 0
	}
	return
}

type checkFunc func(rune) bool

func (l *Lexer) read(check checkFunc) (string, error) {
	var builder strings.Builder
	for {
		ch, err := l.nextChar()
		if err == io.EOF {
			return builder.String(), nil
		} else if err != nil {
			return builder.String(), err
		}

		if check(ch) {
			builder.WriteRune(ch)
		} else {
			l.reader.UnreadRune()
			return builder.String(), nil
		}
	}
}

func (l *Lexer) readNumber() (string, error) {
	return l.read(
		func(r rune) bool {
			return unicode.IsDigit(r) || r == '.' || r == 'x' || r == 'b' || r == 'o'
		},
	)
}

func (l *Lexer) readSpace() (string, error) {
	return l.read(unicode.IsSpace)
}

func (l *Lexer) readOp(op string) string {
	i := 0
	rs := []rune(op)
	op, _ = l.read(func(r rune) bool {
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
		ch, err := l.peek()
		if err == io.EOF {
			return token.New(token.TOKEN_TYPE_EOF, "EOF")
		} else if err != nil {
			fmt.Printf("read next token failed: %v\n", err)
			os.Exit(1)
		}

		if unicode.IsSpace(ch) {
			_, err := l.readSpace()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			continue
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

		switch ch {
		case '+':
			return token.New(token.TOKEN_TYPE_OP_ADD, l.readOp("+"))
		case '-':
			return token.New(token.TOKEN_TYPE_OP_SUB, l.readOp("-"))
		case '*':
			return token.New(token.TOKEN_TYPE_OP_MUL, l.readOp("*"))
		case '/':
			return token.New(token.TOKEN_TYPE_OP_DIV, l.readOp("/"))
		case '%':
			return token.New(token.TOKEN_TYPE_OP_MOD, l.readOp("%"))
		case '(':
			return token.New(token.TOKEN_TYPE_LPAREN, l.readOp("("))
		case ')':
			return token.New(token.TOKEN_TYPE_RPAREN, l.readOp(")"))
		}

		panic(fmt.Sprintf("unrecognize character: %c", ch))
	}

}
