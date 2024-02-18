package main

import "fmt"

func main() {
	src := "let x = 5;"

	tokens := Tokenize(src)

	fmt.Println(tokens)
}
