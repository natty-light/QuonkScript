package runtime

import (
	"QuonkScript/parser"
)

func Evaluate(astNode parser.Stmt, scope *Scope) RuntimeValue {
	switch astNode.GetKind() {
	case parser.NumericLiteralNode:
		return MakeNumber(astNode.(parser.NumericLiteral).Value)
	// Return a null by default
	case parser.NullLiteralNode:
		return MakeNull()
	case parser.BinaryExprNode:
		// This typecast is safe
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
	default:
		parser.PrintAST(astNode)
		panic("This NodeType has not been implemented")
	}
}
