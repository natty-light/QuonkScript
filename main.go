package main

import (
	"QuonkScript/parser"
	"QuonkScript/runtime"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	repl()
}

func repl() {
	p := parser.Parser{}
	scope := &runtime.Scope{Parent: nil, Variables: make(map[string]runtime.RuntimeValue)}
	fmt.Println("REPL v0.1")
	in := bufio.NewReader(os.Stdin)

	scope.DeclareVariable("x", runtime.NumberValue{TypedValue: runtime.TypedValue{Type: runtime.NumberValueType}, Value: 100})

	for {
		fmt.Print("> ")
		input, err := in.ReadString('\n')

		if err != nil {
			fmt.Println("Invalid Input")
			return
		}
		if strings.Contains(input, "exit") {
			os.Exit(0)
		}
		prog := p.ProduceAST(input)

		result := runtime.Evaluate(prog, scope)
		fmt.Println(result)
	}
}
