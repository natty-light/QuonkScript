package main

import (
	lexer "QuonkScript/Lexer"
	"fmt"
)

func main() {
	src := "let x = 5;"

	tokens := lexer.Tokenize(src)

	fmt.Println(tokens)
}
