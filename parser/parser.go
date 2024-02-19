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

// returns first token in tokens array
func (P *Parser) at() lexer.Token {
	return P.tokens[0]
}

// removes first from tokens array and returns it
func (P *Parser) eat() lexer.Token {
	// Pull out first token
	prev := P.at()
	// Remove prev
	P.tokens = P.tokens[1:]

	return prev
}

func (P *Parser) eatExpected(expected lexer.TokenType, err string) lexer.Token {
	prev := P.at()
	P.tokens = P.tokens[1:]

	if prev.Type != lexer.TokenType(expected) {
		panic(err)
	}
	return prev
}

// lexes, tokenizes, and produces a Program AST
func (P *Parser) ProduceAST(src string) Program {
	// Create token array
	P.tokens = lexer.Tokenize(src)
	program := Program{Kind: ProgramNode, Body: make([]Stmt, 0)}

	for P.NotEOF() {
		// Push expressions onto body
		program.Body = append(program.Body, P.ParseStatement())
	}

	return program
}

// returns whether head of token array is EOF
func (P *Parser) NotEOF() bool {
	return P.tokens[0].Type != lexer.EOF
}

// parses Stmt into Expr
func (P *Parser) ParseStatement() Stmt {
	return P.ParseExpr()
}

// Parse Expr, starts parsing at highest implemented level of following
// orders of precedence:
//
//	AssignmentExpr
//	MemberExpr
//	FunctionCall
//	LogicalExpr
//	ComparisonExpr
//	AdditiveExpr
//	MultiplicativeExpr
//	UnaryExpr
//	PrimaryExpr
func (P *Parser) ParseExpr() Expr {
	return P.ParseAdditiveExpr()
}

// parse primary expression
func (P *Parser) ParsePrimaryExpr() Expr {
	token := P.at().Type

	switch token {
	case lexer.Null:
		P.eat()
		return NullLiteral{ExprStmt: ExprStmt{Kind: NullLiteralNode}, Value: "null"}
	case lexer.Number:
		val, _ := strconv.ParseFloat(P.eat().Value, 64)
		return NumericLiteral{Value: val, ExprStmt: ExprStmt{Kind: NumericLiteralNode}}
	case lexer.Identifier:
		return Ident{Symbol: P.eat().Value, ExprStmt: ExprStmt{Kind: IdentifierNode}}
	case lexer.OpenParen:
		P.eat() // eat the opening paren
		val := P.ParseExpr()
		P.eatExpected(lexer.CloseParen, "Missing close paren") // eat closing paren
		return val

	default:
		bytes, err := json.Marshal(P.at())
		if err == nil {
			fmt.Println(string(bytes))
		}
		// Do something better than panicking here
		panic("Unexpeceted token found during parsing")
	}
}

// Parses additive expressions with left to right precendence for order of operations.
// Also kicks off ParseMultiplicativeExpr()
func (P *Parser) ParseAdditiveExpr() Expr {
	left := P.ParseMultiplicativeExpr()

	for P.at().Value == "+" || P.at().Value == "-" {
		// recall that next pops the head off the tokens array of Parser
		operator := P.eat().Value
		right := P.ParseMultiplicativeExpr()

		// This bubbles up the expr
		left = BinaryExpr{ExprStmt: ExprStmt{Kind: BinaryExprNode}, Left: left, Right: right, Operator: operator}
	}
	return left
}

// Parses multiplicative expressions with left to right precendence for order of operations.
// Also kicks off ParsePrimaryExpr()
func (P *Parser) ParseMultiplicativeExpr() Expr {
	left := P.ParsePrimaryExpr()

	for P.at().Value == "*" || P.at().Value == "/" || P.at().Value == "%" {

		operator := P.eat().Value
		right := P.ParsePrimaryExpr()

		// This bubbles up the tree
		left = BinaryExpr{ExprStmt: ExprStmt{Kind: BinaryExprNode}, Left: left, Right: right, Operator: operator}
	}
	return left
}

func PrintAST(stmt Stmt) {
	bytes, err := json.MarshalIndent(stmt, "", "    ")
	if err != nil {
		return
	}
	str := string(bytes)
	str = strings.ReplaceAll(str, "\"Kind\": 1", "Program")
	str = strings.ReplaceAll(str, "\"Kind\": 2", "NumericLiteral")
	str = strings.ReplaceAll(str, "\"Kind\": 3", "Null")
	str = strings.ReplaceAll(str, "\"Kind\": 4", "Identifier")
	str = strings.ReplaceAll(str, "\"Kind\": 5", "BinaryExpr")

	fmt.Println(str)
}
