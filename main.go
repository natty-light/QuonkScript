package main

import (
	"QuonkScript/lexer"
	parser "QuonkScript/parser"
	"bufio"
	"fmt"
	"os"
)

func main() {
	repl()
}

func repl() {
	p := parser.Parser{}
	fmt.Println("REPL v0.1")
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := in.ReadString('\n')

		if err != nil {
			fmt.Println("Invalid Input")
			return
		}

		if input == "exit" {
			os.Exit(0)
		}
		tokens := lexer.Tokenize(input)
		fmt.Println(tokens)
		prog := p.ProduceAST(input)
		parser.PrintAST(prog)
	}
}
