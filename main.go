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

	fd, err := os.Open("code.goc")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	aLexer := lexer.New(fd)
	aParser := parser.New(aLexer)
	ir := interpreter.New(aParser)
	result, valueType := ir.Interpret()
	fmt.Println(result, valueType)
}
