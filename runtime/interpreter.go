package runtime

import (
	"QuonkScript/parser"
)

// Typecasts used in ths function should be safe since we are careful about how we assign node types
func Evaluate(astNode parser.Stmt, scope *Scope) RuntimeValue {
	switch astNode.GetKind() {
	case parser.NumericLiteralNode:
		return MakeNumber(astNode.(parser.NumericLiteral).Value)
	// Return a null by default
	case parser.NullLiteralNode:
		return MakeNull()
	case parser.BooleanLiteralNode:
		return MakeBoolean(astNode.(parser.BooleanLiteral).Value)
	case parser.BinaryExprNode:
		return evalBinaryExpr(astNode.(parser.BinaryExpr), scope)
	case parser.IdentifierNode:
		return evalIdentifier(astNode.(parser.Ident), scope)
	// handle statements
	case parser.ProgramNode:
		return evalProgram(astNode.(parser.Program), scope)
	case parser.VarDeclarationNode:
		return evalVarDeclaration(astNode.(parser.VarDeclaration), scope)
	case parser.AssignmentNode:
		return evalAssignmentExpr(astNode.(parser.VarAssignemntExpr), scope)
	case parser.ObjectLiteralNode:
		return evalObjectExpr(astNode.(parser.ObjectLiteral), scope)
	case parser.InternalFunctionCallExprNode:
		return evalInternalFuncCallExpr(astNode.(parser.InternalFunctionCallExpr), scope)
	case parser.ComparisonExprNode:
		return evalComparisonExpr(astNode.(parser.ComparisonExpr), scope)
	default:
		parser.PrintAST(astNode)
		panic("This NodeType has not been implemented")
	}
}
