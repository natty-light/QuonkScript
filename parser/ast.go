package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

type NodeType int

const (
	// Statements
	ProgramNode NodeType = iota + 1
	VarDeclarationNode

	// Expressions
	NumericLiteralNode
	NullLiteralNode
	IdentifierNode
	BinaryExprNode
	AssignmentNode
)

// Node Interfaces
type (
	Node interface {
		GetKind() NodeType
	}
	// Statements will not return a value
	Stmt interface {
		Node
		statementNode()
	}

	Expr interface {
		Stmt
		expressionNode()
	}
)

// Expressions
type (
	ExprStmt struct {
		Kind NodeType
	}

	BinaryExpr struct {
		ExprStmt `json:"kind"` // Type should always be BinaryExprNode
		Left     Expr          `json:"left"`
		Right    Expr          `json:"right"`
		Operator string        `json:"operator"`
	}

	Ident struct {
		ExprStmt `json:"kind"` // Type should always be IndentifierNode
		Symbol   string        `json:"symbol"`
	}

	NumericLiteral struct {
		ExprStmt `json:"kind"` // Type should always be NumericLiteralNode
		Value    float64       `json:"value"`
	}

	NullLiteral struct {
		ExprStmt `json:"kind"` // Type should always be NullLiteralNode
		Value    string        `json:"value"` // value should always be null
	}
	Program struct {
		Kind NodeType // Type should always be ProgramNode but I don't know how to do that in Go
		Body []Stmt
	}

	VarDeclaration struct {
		Kind       NodeType `json:"kind"` // Type should always be VarDeclarationNode but I don't know how to do that in Go
		Constant   bool     `json:"constant"`
		Identifier string   `json:"string"`
		Value      *Expr    `json:"value"` // Variables can be initialized without values
	}

	VarAssignemntExpr struct {
		Kind     NodeType // Type should always be AssignmentNode but I don't know how to do that in Go
		Assignee Expr     // This is important for the implementation of objects in supporting complex expressions
		Value    Expr
	}
)

// Implement Node methods
func (e ExprStmt) GetKind() NodeType {
	return e.Kind
}

func (b BinaryExpr) GetKind() NodeType {
	return b.Kind
}

func (i Ident) GetKind() NodeType {
	return i.Kind
}

func (n NumericLiteral) GetKind() NodeType {
	return n.Kind
}

func (p Program) GetKind() NodeType {
	return p.Kind
}

func (v VarDeclaration) GetKind() NodeType {
	return v.Kind
}

func (v VarAssignemntExpr) GetKind() NodeType {
	return v.Kind
}

// Implement expression and statements
func (i Ident) expressionNode() {}
func (i Ident) statementNode()  {}

func (e ExprStmt) expressionNode() {}
func (e ExprStmt) statementNode()  {}

func (b BinaryExpr) expressionNode() {}
func (b BinaryExpr) statementNode()  {}

func (n NumericLiteral) expressionNode() {}
func (n NumericLiteral) statementNode()  {}

func (n NullLiteral) expressionNode() {}
func (n NullLiteral) statementNode()  {}

func (p Program) statementNode() {}

func (v VarDeclaration) statementNode() {}

func (v VarAssignemntExpr) statementNode()  {}
func (v VarAssignemntExpr) expressionNode() {}

func PrintAST(stmt Stmt) {
	bytes, err := json.MarshalIndent(stmt, "", "    ")
	if err != nil {
		return
	}
	str := string(bytes)
	str = strings.ReplaceAll(str, "\"Kind\": 1", "Program")
	str = strings.ReplaceAll(str, "\"Kind\": 2", "VarDeclaration")
	str = strings.ReplaceAll(str, "\"Kind\": 3", "NumericLiteral")
	str = strings.ReplaceAll(str, "\"Kind\": 4", "Null")
	str = strings.ReplaceAll(str, "\"Kind\": 5", "Identifier")
	str = strings.ReplaceAll(str, "\"Kind\": 6", "BinaryExpr")
	str = strings.ReplaceAll(str, "\"Kind\": 7", "AssignmentExpr")

	fmt.Println(str)
}
