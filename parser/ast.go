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

	// Literals
	NumericLiteralNode
	NullLiteralNode
	IdentifierNode
	PropertyLiteralNode
	ObjectLiteralNode

	// Expressions
	BinaryExprNode
	AssignmentNode
	MemberExprNode
	FunctionCallExprNode
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

	ObjectLiteral struct {
		Kind       NodeType          `json:"kind"` // Type should always be ObjectLiteralNode
		Properties []PropertyLiteral `json:"properties"`
	}

	PropertyLiteral struct {
		Kind  NodeType `json:"kind"` // Type should always be PropertyLiteralNode
		Key   string   `json:"key"`
		Value *Expr    `json:"value"` // Pointer so it can be nil
	}

	MemberExpr struct {
		Kind     NodeType `json:"kind"` // Type should always be MemberExprNode
		Object   Expr     `json:"object"`
		Field    Expr     `json:"property"`
		Computed bool     `json:"computed"`
	}

	FunctionCallExpr struct {
		Kind   NodeType `json:"kind"` // Type should always be FunctionCallExprNode
		Args   []Expr   `json:"args"`
		Caller Expr     `json:"caller"`
	}
)

// Implement Node methods
func (e ExprStmt) GetKind() NodeType {
	return e.Kind
}

func (b BinaryExpr) GetKind() NodeType {
	return BinaryExprNode
}

func (i Ident) GetKind() NodeType {
	return IdentifierNode
}

func (n NumericLiteral) GetKind() NodeType {
	return NumericLiteralNode
}

func (p Program) GetKind() NodeType {
	return ProgramNode
}

func (v VarDeclaration) GetKind() NodeType {
	return VarDeclarationNode
}

func (v VarAssignemntExpr) GetKind() NodeType {
	return AssignmentNode
}

func (o ObjectLiteral) GetKind() NodeType {
	return ObjectLiteralNode
}

func (p PropertyLiteral) GetKind() NodeType {
	return PropertyLiteralNode
}

func (m MemberExpr) GetKind() NodeType {
	return MemberExprNode
}

func (f FunctionCallExpr) GetKind() NodeType {
	return FunctionCallExprNode
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

func (o ObjectLiteral) expressionNode() {}
func (o ObjectLiteral) statementNode()  {}

func (p PropertyLiteral) expressionNode() {}
func (p PropertyLiteral) statementNode()  {}

func (m MemberExpr) expressionNode() {}
func (m MemberExpr) statementNode()  {}

func (f FunctionCallExpr) expressionNode() {}
func (f FunctionCallExpr) statementNode()  {}

func PrintAST(stmt Stmt) {
	bytes, err := json.MarshalIndent(stmt, "", "    ")
	if err != nil {
		return
	}
	str := string(bytes)
	str = strings.ReplaceAll(str, "\"Kind\": 10", "MemberExpr")
	str = strings.ReplaceAll(str, "\"kind\": 10", "Kind: MemberExpr")
	str = strings.ReplaceAll(str, "\"Kind\": 11", "FunctionCallExpr")
	str = strings.ReplaceAll(str, "\"kind\": 11", "Kind: FunctionCallExpr")
	str = strings.ReplaceAll(str, "\"Kind\": 1", "Program")
	str = strings.ReplaceAll(str, "\"kind\": 1", "Kind: Program")
	str = strings.ReplaceAll(str, "\"Kind\": 2", "VarDeclaration")
	str = strings.ReplaceAll(str, "\"kind\": 2", "Kind: VarDeclaration")
	str = strings.ReplaceAll(str, "\"Kind\": 3", "NumericLiteral")
	str = strings.ReplaceAll(str, "\"kind\": 3", "Kind: NumericLiteral")
	str = strings.ReplaceAll(str, "\"Kind\": 4", "NullLiteral")
	str = strings.ReplaceAll(str, "\"kind\": 4", "Kind: NullLiteral")
	str = strings.ReplaceAll(str, "\"Kind\": 5", "Identifier")
	str = strings.ReplaceAll(str, "\"kind\": 5", "Kind: Indentifier")
	str = strings.ReplaceAll(str, "\"Kind\": 6", "PropertyLiteral")
	str = strings.ReplaceAll(str, "\"kind\": 6", "Kind: PropertyLiteral")
	str = strings.ReplaceAll(str, "\"Kind\": 7", "ObjectLiteral")
	str = strings.ReplaceAll(str, "\"kind\": 7", "Kind: ObjectLiteral")
	str = strings.ReplaceAll(str, "\"Kind\": 8", "BinaryExpr")
	str = strings.ReplaceAll(str, "\"kind\": 8", "Kind: BinaryExpr")
	str = strings.ReplaceAll(str, "\"Kind\": 9", "AssignmentExpr")
	str = strings.ReplaceAll(str, "\"kind\": 9", "Kind: AssingmentExpr")

	fmt.Println(str)
}
