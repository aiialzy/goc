package main

import (
	"fmt"
	"goc/interpreter"
	"goc/lexer"
	"goc/parser"
	"goc/token"
	"os"
)

const (
	IS_PRINT_TOKEN = false
)

func main() {

	token.SetPrintToken(IS_PRINT_TOKEN)

	source, err := os.ReadFile("code.goc")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	aLexer := lexer.New(string(source))
	// printAllToken(aLexer)

	interpret(aLexer)
}

func interpret(aLexer *lexer.Lexer) {
	aParser := parser.New(aLexer)
	ir := interpreter.New(aParser)
	ir.Interpret()
	ir.PrintValues()
}

func printAllToken(aLexer *lexer.Lexer) {
	for {
		t := aLexer.NextToken()
		if t.Type == token.TOKEN_TYPE_EOF {
			break
		}

		fmt.Println(t)
	}
}
