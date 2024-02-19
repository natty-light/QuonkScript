package main

import (
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

		prog := p.ProduceAST(input)
		parser.PrintAST(prog)
	}
}
