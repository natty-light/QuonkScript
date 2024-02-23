package parser

import (
	"QuonkScript/lexer"
	"encoding/json"
	"fmt"
	"strconv"
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
	switch P.at().Type {
	case lexer.Mut:
		fallthrough
	case lexer.Const:
		return P.ParseVarDeclaration()
	case lexer.If:
		return P.ParseBranchStmt()
	default:
		return P.ParseExpr()
	}
}

// Parse Expr, starts parsing at highest implemented level of following
// orders of precedence:
//
//		AssignmentExpr
//		ObjectExpr
//		ChainedComparisonExpr
//		ComparisonExpr
//		AdditiveExpr
//		MultiplicativeExpr
//	 	FunctionCallExpr
//		MemberExpr
//		PrimaryExpr
//
// Kicks off ParseAssignmentExpr()
func (P *Parser) ParseExpr() Expr {
	return P.ParseAssignmentExpr()
}

// Parses Assignment expressions with left to right precedence
// Also kicks off ParseObjectExpr()
func (P *Parser) ParseAssignmentExpr() Expr {
	left := P.ParseObjectExpr() // this will be swapped for objects

	if P.at().Type == lexer.Equals {
		P.eat()                          // advance past equals token
		value := P.ParseAssignmentExpr() // we want to allow chaining so we must call recursively

		return VarAssignemntExpr{Value: value, Assignee: left, Kind: AssignmentExprNode}
	}

	return left
}

// Parses object expressions with left to right precedence
// Also kicks off ParseAdditiveExpr()
func (P *Parser) ParseObjectExpr() Expr {
	if P.at().Type != lexer.OpenCurlyBracket {
		return P.ParseChainedLogicalExpr() // If we do not find an open brace, proceed on
	}

	P.eat() // advance past open brace
	properties := make([]PropertyLiteral, 0)

	for P.NotEOF() && P.at().Type != lexer.CloseCurlyBracket {
		key := P.eatExpected(lexer.Identifier, "Honk! Expected field name following bracket in object literal").Value

		switch P.at().Type {
		// Allow shorthand {key, }
		case lexer.Comma:
			P.eat() // Advance past comma
			// append property with no value
			properties = append(properties, PropertyLiteral{Value: nil, Kind: PropertyLiteralNode, Key: key})
		// Allow shorthand { key }
		case lexer.CloseCurlyBracket:
			// append property with no value
			properties = append(properties, PropertyLiteral{Value: nil, Kind: PropertyLiteralNode, Key: key})
		// { key: value }
		default:
			P.eatExpected(lexer.Colon, "Honk! Expected colon following property name in object literal")
			val := P.ParseExpr() // Allow any expression
			// Append property with value
			properties = append(properties, PropertyLiteral{Value: &val, Kind: PropertyLiteralNode, Key: key})

			if P.at().Type != lexer.CloseCurlyBracket {
				P.eatExpected(lexer.Comma, "Honk! Expected comma or closing brace at end of object literal")
			}
		}

	}
	P.eatExpected(lexer.CloseCurlyBracket, "Honk! Expected closing bracket for object literal")
	return ObjectLiteral{Kind: ObjectLiteralNode, Properties: properties}
}

func (P *Parser) ParseComparisonExpr() Expr {
	left := P.ParseAdditiveExpr()

	if P.at().Type == lexer.Equality || P.at().Type == lexer.NotEqual || P.at().Type == lexer.GreaterThan || P.at().Type == lexer.LessThan || P.at().Type == lexer.GreaterEqualTo || P.at().Type == lexer.LessEqualTo {
		operator := P.eat().Value
		right := P.ParseAdditiveExpr()

		left = ComparisonExpr{Kind: ComparisonExprNode, Left: left, Right: right, Operator: operator}
	}

	return left
}

func (P *Parser) ParseChainedLogicalExpr() Expr {
	left := P.ParseComparisonExpr()

	for P.at().Type == lexer.And || P.at().Type == lexer.Or {
		operator := P.eat().Value
		right := P.ParseComparisonExpr()

		left = ComparisonExpr{
			Kind:     ComparisonExprNode,
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}
	return left
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
// Function kicks off ParseCallMemberExpr
func (P *Parser) ParseMultiplicativeExpr() Expr {
	left := P.ParseCallMemberExpr()

	for P.at().Value == "*" || P.at().Value == "/" || P.at().Value == "%" {

		operator := P.eat().Value
		right := P.ParseCallMemberExpr()

		// This bubbles up the tree
		left = BinaryExpr{ExprStmt: ExprStmt{Kind: BinaryExprNode}, Left: left, Right: right, Operator: operator}
	}
	return left
}

// This function is different as it takes in an Expr argument
func (P *Parser) ParseFunctionCallExpr(caller Expr) Expr {
	args := P.ParseArguments()
	var callExpr Expr = InternalFunctionCallExpr{Kind: InternalFunctionCallExprNode, Caller: caller, Args: args} // no walrus here since we need callExpr to just be an Expr

	// This allows us to recursively chain function calls
	if P.at().Type == lexer.OpenParen {
		callExpr = P.ParseFunctionCallExpr(callExpr)
	}

	return callExpr
}

// This function parses object member expressions
// Also, it kicks of ParsePrimaryExpr
func (P *Parser) ParseMemberExpr() Expr {
	obj := P.ParsePrimaryExpr()

	for P.at().Type == lexer.Dot || P.at().Type == lexer.OpenSquareBracket {
		// Get . or [
		operator := P.eat()
		var field Expr
		var computed bool

		if operator.Type == lexer.Dot {
			// non computed branch obj.field
			computed = false
			// We expect P.at() to be an Identifier
			field = P.ParsePrimaryExpr()

			if field.GetKind() != IdentifierNode {
				panic("Honk! Attempt to reference object field with something other than an identifier")
			}

		} else {
			// Computed branch
			computed = true
			// this allows obj[computed]
			field = P.ParseExpr()
			P.eatExpected(lexer.CloseSquareBracket, "Honk! Expected closing bracket for object field access")
		}
		obj = MemberExpr{Kind: MemberExprNode, Field: field, Computed: computed, Object: obj}
	}
	return obj
}

// parse primary expression, bottom of call stack
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
	case lexer.True:
		P.eat()
		return BooleanLiteral{Value: true, ExprStmt: ExprStmt{Kind: BooleanLiteralNode}}
	case lexer.False:
		P.eat()
		return BooleanLiteral{Value: false, ExprStmt: ExprStmt{Kind: BooleanLiteralNode}}
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

// Parses Object Member expressions with left to right precedence, can parse members recursively
func (P *Parser) ParseCallMemberExpr() Expr {
	member := P.ParseMemberExpr() // Will fall through to parse primary

	if P.at().Type == lexer.OpenParen {
		return P.ParseFunctionCallExpr(member)
	}
	return member
}

// This function parses arguments for a function call
// arguments are not parameters, args are just expressions
func (P *Parser) ParseArguments() []Expr {
	// To get here, P.at() is an open paren but this doesn't hurt
	P.eatExpected(lexer.OpenParen, "Honk! Expected parenthesis after function call")
	args := make([]Expr, 0)
	// if the next token is anything other than
	if P.at().Type != lexer.CloseParen {
		args = P.ParseArgumentList()
	}
	P.eatExpected(lexer.CloseParen, "Honk! Expected closing parenthesis for function call")
	return args
}

func (P *Parser) ParseArgumentList() []Expr {
	// Parse first arg
	args := []Expr{P.ParseAssignmentExpr()}

	// In JS impl, this line is while P.at().Type == TokenType.Comma && P.eat(), but I don't know what the equivalent is in Go so lets try this
	for P.at().Type != lexer.EOF && P.at().Type == lexer.Comma {
		P.eat()
		args = append(args, P.ParseAssignmentExpr())
	}

	return args
	// No need to eatExpected() here as we do that in the calling function
}

// Parses variable declaration expr stmt
func (P *Parser) ParseVarDeclaration() Stmt {
	// eat advances
	isConstant := P.eat().Type == lexer.Const
	// eatExpected advances
	identifier := P.eatExpected(lexer.Identifier, "Expected variable name").Value

	if P.at().Type == lexer.Semicolon {
		P.eat() // Advance
		if isConstant {
			panic("Constant variables must be initialized")
		}
		// Mutable variable declaration Node
		return VarDeclaration{Kind: VarDeclarationNode, Identifier: identifier, Constant: false, Value: nil}
	}

	P.eatExpected(lexer.Equals, "Expected equals following variable name in declaration")
	value := P.ParseExpr()
	// is this pointer fucked?
	declaration := VarDeclaration{Kind: VarDeclarationNode, Value: &value, Constant: isConstant, Identifier: identifier}
	P.eatExpected(lexer.Semicolon, "Missing semicolon following variable declaration")
	return declaration
}

func (P *Parser) ParseBranchStmt() Stmt {
	P.eat() // move past if
	P.eatExpected(lexer.OpenParen, "Honk! Expected opening ( before condition of if statement")
	condition := P.ParseChainedLogicalExpr() // We want to be able to allow things like if (x + 6 > 8 * 2)
	P.eatExpected(lexer.CloseParen, "Honk!, Expected closing ) following condition of if statement")
	P.eatExpected(lexer.OpenCurlyBracket, "Honk! Expected opening { following if statement")

	body := make([]Stmt, 0)

	for P.at().Type != lexer.CloseCurlyBracket {
		body = append(body, P.ParseStatement())
	}
	P.eatExpected(lexer.CloseCurlyBracket, "Honk! Expected } after body of if statement")

	// TODO: Implement elseif?
	elseBody := make([]Stmt, 0)

	if P.at().Type == lexer.Else {
		P.eat() // advance past else
		P.eatExpected(lexer.OpenCurlyBracket, "Honk! Expected { after else")
		for P.at().Type != lexer.CloseCurlyBracket {
			elseBody = append(elseBody, P.ParseStatement())
		}
		P.eatExpected(lexer.CloseCurlyBracket, "Honk! Expected } after body of else")
	}

	return BranchStmt{Condition: condition, Else: elseBody, Body: body}
}
