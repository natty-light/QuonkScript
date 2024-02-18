package main

import (
	"fmt"
	"os"
)

func main() {

	args := os.Args[1:]
	filename := args[0]

	bytes, err := os.ReadFile(filename)

	if err != nil {
		panic(fmt.Sprintf("Unable to read file %s", filename))
	}

	src := string(bytes)
	tokens := Tokenize(src)

	fmt.Println(tokens)
}
