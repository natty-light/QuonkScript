package parser

import (
	"QuonkScript/lexer"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	tokens []lexer.Token
}

func (P *Parser) at() lexer.Token {
	return P.tokens[0]
}

func (P *Parser) next() lexer.Token {
	// Pull out first token
	prev := P.at()
	// Remove prev
	P.tokens = P.tokens[1:]

	return prev
}

// Orders of precedence

// AssignmentExpr
// MemberExpr
// FunctionCall
// LogicalExpr
// ComparisonExpr
// AdditiveExpr
// MultiplicativeExpr
// UnaryExpr
// PrimaryExpr

func (P *Parser) ProduceAST(src string) Program {
	// Create token array
	P.tokens = lexer.Tokenize(src)
	program := Program{Kind: ProgramNode, Body: make([]Stmt, 0)}

	for P.NotEOF() {
		program.Body = append(program.Body, P.ParseStatement())
	}

	return program
}

func (P *Parser) NotEOF() bool {
	return P.tokens[0].Type != lexer.EOF
}

func (P *Parser) ParseStatement() Stmt {
	return P.ParseExpr()
}

func (P *Parser) ParseExpr() Expr {
	return P.ParseAdditiveExpr()
}

func (P *Parser) ParsePrimaryExpr() Expr {
	token := P.at().Type

	switch token {
	case lexer.Identifier:
		return Ident{Symbol: P.next().Value, ExprStmt: ExprStmt{Kind: IdentifierNode}}
	case lexer.Number:
		val, _ := strconv.ParseFloat(P.next().Value, 64)
		return NumericLiteral{Value: val, ExprStmt: ExprStmt{Kind: NumericLiteralNode}}
	default:
		// Do something better than panicking here
		panic("Unexpeceted token found during parsing")
	}
}

// left to right precedence
func (P *Parser) ParseAdditiveExpr() Expr {
	left := P.ParsePrimaryExpr()

	for P.at().Value == "+" || P.at().Value == "-" {
		operator := P.next().Value
		right := P.ParsePrimaryExpr()

		// This bubbles up the expr
		left = BinaryExpr{ExprStmt: ExprStmt{Kind: BinaryExprNode}, Left: left, Right: right, Operator: operator}
	}
	return left
}

func PrintAST(prog Program) {
	bytes, err := json.MarshalIndent(prog, "", "    ")
	if err != nil {
		return
	}
	str := string(bytes)
	str = strings.ReplaceAll(str, "\"Kind\": 1", "Program")
	str = strings.ReplaceAll(str, "\"Kind\": 2", "NumericLiteral")
	str = strings.ReplaceAll(str, "\"Kind\": 3", "Identifier")
	str = strings.ReplaceAll(str, "\"Kind\": 4", "BinaryExpr")

	fmt.Println(str)
}
