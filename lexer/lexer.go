package lexer

import (
	"QuonkScript/utils"
	"fmt"
	"regexp"
	"strings"
)

type TokenType int

const (
	// Literals
	Null TokenType = iota + 1
	Number
	Identifier

	// Keywords
	Mut
	Const
	True
	False
	If
	Else
	Elseif
	Func
	Return

	// Grouping and operations
	Equals
	Semicolon
	OpenParen
	CloseParen
	BinaryOperator
	OpenCurlyBracket
	CloseCurlyBracket
	Comma
	Colon
	OpenSquareBracket
	CloseSquareBracket
	Dot
	Equality
	GreaterThan
	LessThan
	GreaterEqualTo
	LessEqualTo
	NotEqual
	And
	Or

	EOF // End of File
)

const (
	leftParen          = "("
	rightParen         = ")"
	addSym             = "+"
	multSym            = "*"
	divSym             = "/"
	subSym             = "-"
	eqSym              = "="
	modSym             = "%"
	semi               = ";"
	leftCurlyBracket   = "{"
	rightCurlyBracket  = "}"
	comma              = ","
	colon              = ":"
	leftSquareBracket  = "["
	rightSquareBracket = "]"
	dot                = "."
	greaterThan        = ">"
	lessThan           = "<"
	equality           = "=="
	leq                = "<="
	geq                = ">="
	bang               = "!"
	ampersand          = "&"
	pipe               = "|"
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
	return s == " " || s == "\n" || s == "\t" || s == "\r"
}

func token(Type TokenType, Value string) Token {
	return Token{Type: Type, Value: Value}
}

func getKeywordMap() map[string]TokenType {
	return map[string]TokenType{
		"mut":    Mut,
		"const":  Const,
		"null":   Null,
		"true":   True,
		"false":  False,
		"if":     If,
		"else":   Else,
		"elseif": Elseif,
		"func":   Func,
		"return": Return,
	}
}

func Tokenize(source string) []Token {
	tokens := make([]Token, 0)
	keywords := getKeywordMap()

	src := strings.Split(source, "")
	remaining := len(src)

	// Build each token
	for len(src) > 0 {

		// src[0] will always be defined because len(src) > 0
		char := src[0]

		switch char {
		case leftParen:
			tokens = append(tokens, token(OpenParen, char))
			// Pop first char
			src = utils.Pop(src)
			remaining--

		case rightParen:
			tokens = append(tokens, token(CloseParen, char))
			src = utils.Pop(src)
			remaining--

		case leftCurlyBracket:
			tokens = append(tokens, token(OpenCurlyBracket, char))
			src = utils.Pop(src)
			remaining--

		case rightCurlyBracket:
			tokens = append(tokens, token(CloseCurlyBracket, char))
			src = utils.Pop(src)
			remaining--

		case leftSquareBracket:
			tokens = append(tokens, token(OpenSquareBracket, char))
			src = utils.Pop(src)
			remaining--
		case rightSquareBracket:
			tokens = append(tokens, token(CloseSquareBracket, char))
			src = utils.Pop(src)
			remaining--
		case addSym:
			fallthrough
		case subSym:
			fallthrough
		case divSym:
			fallthrough
		case modSym:
			fallthrough
		case multSym:
			tokens = append(tokens, token(BinaryOperator, char))
			src = utils.Pop(src)
			remaining--
		case eqSym:
			// check for equality symbol here
			src = utils.Pop(src)
			remaining--
			if remaining > 0 { // If there are symbols left to parse
				nextChar := src[0] // This is safe because remaining > 0
				// another equals sign
				if nextChar == eqSym {
					tokens = append(tokens, token(Equality, "=="))
					src = utils.Pop(src)
					remaining--
				} else {
					tokens = append(tokens, token(Equals, char))
				}
			} else {
				tokens = append(tokens, token(Equals, char))
			}
		case greaterThan: // >= or >
			// need to check next char
			src = utils.Pop(src)
			remaining--
			if remaining > 0 { // If there are symbols left to parse
				nextChar := src[0] // This is safe because remaining > 0
				// looking for equals sign
				if nextChar == eqSym {
					tokens = append(tokens, token(GreaterEqualTo, ">="))
					src = utils.Pop(src)
					remaining--
				} else {
					tokens = append(tokens, token(GreaterThan, char))
				}
			} else {
				tokens = append(tokens, token(GreaterThan, char))
			}
		case lessThan: // <= or <
			// need to check next char
			src = utils.Pop(src)
			remaining--
			if remaining > 0 { // If there are symbols left to parse
				nextChar := src[0] // This is safe because remaining > 0
				// looking for equals sign
				if nextChar == eqSym {
					tokens = append(tokens, token(LessEqualTo, "<="))
					src = utils.Pop(src)
					remaining--
				} else {
					tokens = append(tokens, token(LessThan, char))
				}
			} else {
				tokens = append(tokens, token(LessThan, char))
			}
		case bang:
			// need to check next char
			src = utils.Pop(src)
			remaining--
			if remaining > 0 { // If there are symbols left to parse
				nextChar := src[0] // This is safe because remaining > 0
				// looking for equals sign
				if nextChar == eqSym {
					tokens = append(tokens, token(NotEqual, "!="))
					src = utils.Pop(src)
					remaining--
				}
			}
		case ampersand:
			src = utils.Pop(src)
			remaining--
			if remaining > 0 {
				nextChar := src[0]
				if nextChar == ampersand {
					tokens = append(tokens, token(And, "&&"))
					src = utils.Pop(src)
					remaining--
				} else {
					panic(fmt.Sprintf("Honk! Unrecognized character %s", char))
				}
			}
		case pipe:
			src = utils.Pop(src)
			remaining--
			if remaining > 0 {
				nextChar := src[0]
				if nextChar == pipe {
					tokens = append(tokens, token(And, "||"))
					src = utils.Pop(src)
					remaining--
				} else {
					panic(fmt.Sprintf("Honk! Unrecognized character %s", char))
				}
			}
		case semi:
			tokens = append(tokens, token(Semicolon, char))
			src = utils.Pop(src)
			remaining--
		case colon:
			tokens = append(tokens, token(Colon, char))
			src = utils.Pop(src)
			remaining--
		case comma:
			tokens = append(tokens, token(Comma, char))
			src = utils.Pop(src)
			remaining--
		case dot:
			tokens = append(tokens, token(Dot, char))
			src = utils.Pop(src)
			remaining--
		default:
			// Handle multichar token
			if isNumeric(char) {
				num := ""
				// while there are characters left to Parse and the characters are numeric
				// We don't use char here because we want to process entire multichar number within this switch case
				for len(src) > 0 && isNumeric(src[0]) {
					num += src[0]
					src = utils.Pop(src)
					remaining--
				}

				tokens = append(tokens, token(Number, num))

			} else if isAlpha(char) {
				ident := "" // ident could be a variable name, or it could be a keyword
				for len(src) > 0 && isAlpha(src[0]) {
					ident += src[0]
					src = utils.Pop(src)
					remaining--
				}

				// check for reserved keyword
				// a miss will be the types zero value, so 0
				reserved := keywords[ident]

				// TokenType is iota + 1 so TokenType will always be greater than 0
				if reserved == 0 {
					tokens = append(tokens, token(Identifier, ident))
				} else {
					tokens = append(tokens, token(reserved, ident))
				}

			} else if isSkipable(char) {
				src = utils.Pop(src)
				remaining--
			} else {
				panic(fmt.Sprintf("Unrecognized character %s", char))
			}
		}
	}
	return append(tokens, token(EOF, "EOF"))
}
