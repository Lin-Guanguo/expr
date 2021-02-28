package main

import (
	"bytes"
	"expr/lexer"
	"expr/parser"
	"fmt"
)

func simpleTestParse(input string) *parser.BlockExpr {
	buf := bytes.NewBufferString(input)
	l := lexer.New(buf)
	p := parser.New(l)

	block := p.ParseProgram()
	if len(p.Errors) != 0 {
		fmt.Printf("%s\n", block.String(0))
		for i, e := range p.Errors {
			fmt.Printf("[%d] %s\n", i, e.Error())
		}
	}
	return block
}
