package ast

type NodeType int

const (
	ProgramNode NodeType = iota + 1
	NumericLiteralNode
	IdentifierNode
	BinaryExprNode
)

// Statements will not return a value
type Statement struct {
	Type NodeType
}

type Program struct {
	Statement // Type should always be ProgramNode but I don't know how to do that in Go
	Body      []Statement
}

type Expression struct {
	Statement
}

type BinaryExpr struct {
	Expression // Type should always be BinaryExprNode
	left       Expression
	right      Expression
	operator   string
}

type Identifier struct {
	Expression // Type should always be IndentifierNode
	symbol     string
}

type NumericLiteral struct {
	Expression // Type should always be NumericLiteralNode
	value      int
}
