package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	AssignmentExprNode
	MemberExprNode
	InternalFunctionCallExprNode
	ComparisonExprNode

	// These are here so i don't have to fix the printing
	BooleanLiteralNode
	BranchNode
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

	InternalFunctionCallExpr struct {
		Kind   NodeType `json:"kind"` // Type should always be FunctionCallExprNode
		Args   []Expr   `json:"args"`
		Caller Expr     `json:"caller"`
	}

	ComparisonExpr struct {
		Kind     NodeType `json:"kind"`
		Operator string   `json:"operator"`
		Left     Expr
		Right    Expr
	}

	BooleanLiteral struct {
		ExprStmt `json:"kind"` // Type should always be NumericLiteralNode
		Value    bool          `json:"value"`
	}

	BranchStmt struct {
		Kind      NodeType `json:"kind"`
		Condition Expr     `json:"condition"`
		Body      []Stmt   `json:"body"`
		Else      []Stmt   `json:"else"`
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
	return AssignmentExprNode
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

func (f InternalFunctionCallExpr) GetKind() NodeType {
	return InternalFunctionCallExprNode
}

func (c ComparisonExpr) GetKind() NodeType {
	return ComparisonExprNode
}

func (b BooleanLiteral) GetKind() NodeType {
	return BooleanLiteralNode
}

func (b BranchStmt) GetKind() NodeType {
	return BranchNode
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

func (f InternalFunctionCallExpr) expressionNode() {}
func (f InternalFunctionCallExpr) statementNode()  {}

func (c ComparisonExpr) expressionNode() {}
func (c ComparisonExpr) statementNode()  {}

func (b BooleanLiteral) expressionNode() {}
func (b BooleanLiteral) statementNode()  {}

func (b BranchStmt) statementNode() {}

func PrintAST(stmt Stmt) {
	bytes, err := json.MarshalIndent(stmt, "", "    ")
	if err != nil {
		return
	}
	str := string(bytes)

	str = replaceStringsCapital(MemberExprNode, "MemberExpr", str)
	str = replaceStringsLowercase(MemberExprNode, "MemberExpr", str)
	str = replaceStringsCapital(InternalFunctionCallExprNode, "InternalFunctionCallExpr", str)
	str = replaceStringsLowercase(InternalFunctionCallExprNode, "InternalFunctionCallExpr", str)
	str = replaceStringsCapital(ComparisonExprNode, "ComparisonExpr", str)
	str = replaceStringsLowercase(ComparisonExprNode, "ComparisonExpr", str)
	str = replaceStringsCapital(BooleanLiteralNode, "BooleanLiteral", str)
	str = replaceStringsLowercase(BooleanLiteralNode, "BooleanLiteral", str)
	str = replaceStringsCapital(BranchNode, "BranchStmt", str)
	str = replaceStringsLowercase(BranchNode, "BranchStmt", str)
	str = replaceStringsCapital(ProgramNode, "Program", str)
	str = replaceStringsLowercase(ProgramNode, "Program", str)
	str = replaceStringsCapital(VarDeclarationNode, "VarDeclaration", str)
	str = replaceStringsLowercase(VarDeclarationNode, "VarDeclaration", str)
	str = replaceStringsCapital(NumericLiteralNode, "NumericLiteral", str)
	str = replaceStringsLowercase(NumericLiteralNode, "NumericLiteral", str)
	str = replaceStringsCapital(NullLiteralNode, "NullLiteral", str)
	str = replaceStringsLowercase(NullLiteralNode, "NullLiteral", str)
	str = replaceStringsCapital(IdentifierNode, "Identifier", str)
	str = replaceStringsLowercase(IdentifierNode, "Identifier", str)
	str = replaceStringsCapital(PropertyLiteralNode, "PropertyLiteral", str)
	str = replaceStringsLowercase(PropertyLiteralNode, "PropertyLiteral", str)
	str = replaceStringsCapital(ObjectLiteralNode, "ObjectLiteral", str)
	str = replaceStringsLowercase(ObjectLiteralNode, "ObjectLiteral", str)
	str = replaceStringsCapital(BinaryExprNode, "BinaryExpr", str)
	str = replaceStringsLowercase(BinaryExprNode, "BinaryExpr", str)
	str = replaceStringsCapital(AssignmentExprNode, "AssignmentExpr", str)
	str = replaceStringsLowercase(AssignmentExprNode, "AssignmentExpr", str)

	fmt.Println(str)
}

func replaceStringsCapital(node NodeType, replacer string, str string) string {
	return strings.ReplaceAll(str, "\"Kind\": "+strconv.Itoa(int(node)), replacer)
}

func replaceStringsLowercase(node NodeType, replacer string, str string) string {
	return strings.ReplaceAll(str, "\"kind\": "+strconv.Itoa(int(node)), "Kind: "+replacer)
}
