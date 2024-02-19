package parser

type NodeType int

const (
	ProgramNode NodeType = iota + 1
	NumericLiteralNode
	IdentifierNode
	BinaryExprNode
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

// Implement expression
func (i Ident) expressionNode() {}
func (i Ident) statementNode()  {}

func (e ExprStmt) expressionNode() {}
func (e ExprStmt) statementNode()  {}

func (b BinaryExpr) expressionNode() {}
func (b BinaryExpr) statementNode()  {}

func (n NumericLiteral) expressionNode() {}
func (n NumericLiteral) statementNode()  {}

type Program struct {
	Kind NodeType // Type should always be ProgramNode but I don't know how to do that in Go
	Body []Stmt
}
