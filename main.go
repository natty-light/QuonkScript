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
	scope := &runtime.Scope{Parent: nil, Variables: make(map[string]runtime.Variable)}
	fmt.Println("REPL v0.1")
	in := bufio.NewReader(os.Stdin)

	scope.DeclareVariable("true", runtime.MakeBoolean(true), true)
	scope.DeclareVariable("false", runtime.MakeBoolean(false), true)
	// https://www.youtube.com/watch?v=uwKnc4w15nk&list=PL_2VhOvlMk4UHGqYCLWc6GO8FaPl8fQTh&index=5 if you want to have null be an identifier which I do not right now
	// scope.DeclareVariable("null", runtime.MakeNull())

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
