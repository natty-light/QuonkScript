package main

import (
	"regexp"
	"strings"
)

type TokenType int

const (
	Number TokenType = iota + 1
	Identifier
	Equals

	OpenParen
	CloseParen

	BinaryOperator

	Let
	EOL
)

const (
	leftParen  = "("
	rightParen = ")"
	addSym     = "+"
	multSym    = "*"
	divSym     = "/"
	subSym     = "-"
	eqSym      = "="
	semicolon  = ";"
)

type Token struct {
	Value string
	Type  TokenType
}

func isAlpha(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(s)
}

func isNumeric(s string) bool {
	return regexp.MustCompile(`^[0-9]+$`).MatchString(s)
}

func isSkipable(s string) bool {
	return s == " " || s == "\n" || s == "\t"
}

func token(Type TokenType, Value string) Token {
	return Token{Type: Type, Value: Value}
}

func getKeywordMap() map[string]TokenType {
	return map[string]TokenType{
		"let": Let,
	}
}

func Tokenize(source string) []Token {
	tokens := make([]Token, 0)
	keywords := getKeywordMap()

	src := strings.Split(source, "")

	// Build each token
	for len(src) > 0 {

		// src[0] will always be defined because len(src) > 0
		char := src[0]

		switch char {
		case leftParen:
			tokens = append(tokens, token(OpenParen, char))
			// Pop first char
			src = src[1:]
		case rightParen:
			tokens = append(tokens, token(CloseParen, char))
			src = src[1:]
		case addSym:
			fallthrough
		case subSym:
			fallthrough
		case divSym:
			fallthrough
		case multSym:
			tokens = append(tokens, token(BinaryOperator, char))
			src = src[1:]
		case eqSym:
			tokens = append(tokens, token(Equals, char))
			src = src[1:]
		case semicolon:
			tokens = append(tokens, token(EOL, char))
			src = src[1:]
		default:
			// Handle multichar token
			if isNumeric(char) {
				num := ""
				// while there are characters left to Parse and the characters are numeric
				// We don't use char here because we want to process entire multichar number within this switch case
				for len(src) > 0 && isNumeric(src[0]) {
					num += src[0]
					src = src[1:]
				}

				tokens = append(tokens, token(Number, num))

			} else if isAlpha(char) {
				ident := "" // ident could be a variable name, or it could be a keyword
				for len(src) > 0 && isAlpha(src[0]) {
					ident += src[0]
					src = src[1:]
				}

				// check for reserved keyword
				// a miss will be the types zero value, so ""
				reserved := keywords[ident]
				// TokenType is iota + 1 so TokenType will always be greater than 0
				if reserved == 0 {
					tokens = append(tokens, token(Identifier, ident))
				} else {
					tokens = append(tokens, token(reserved, ident))
				}

			} else if isSkipable(char) {
				src = src[1:]
			} else {
				panic("Unrecognized character")
			}
		}
	}
	return tokens
}
