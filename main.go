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
	args := os.Args

	if len(args) == 1 {
		// If no filename was passed as a command line argument, run the repl
		repl()
	} else {
		// Script name should be second arg
		run(args[1])
	}

}

func repl() {
	p := parser.Parser{}
	scope := &runtime.Scope{Parent: nil, Variables: make(map[string]runtime.Variable)}
	fmt.Println("REPL v0.1")
	in := bufio.NewReader(os.Stdin)

	runtime.SetupScope(scope)
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

func run(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Honk! Cannot read file %s\n", filename)
		return
	}

	src := string(bytes)
	p := parser.Parser{}
	scope := &runtime.Scope{Parent: nil, Variables: make(map[string]runtime.Variable)}
	runtime.SetupScope(scope)

	prog := p.ProduceAST(src)
	// parser.PrintAST(prog)
	result := runtime.Evaluate(prog, scope)

	fmt.Println(result)
}
