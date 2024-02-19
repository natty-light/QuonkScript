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
	case parser.ProgramNode:
		return evalProgram(astNode.(parser.Program), scope)
	default:
		parser.PrintAST(astNode)
		panic("This NodeType has not been implemented")
	}
}

func evalBinaryExpr(expr parser.BinaryExpr, scope *Scope) RuntimeValue {
	leftHandSide := Evaluate(expr.Left, scope)
	rightHandSide := Evaluate(expr.Right, scope)

	if leftHandSide.GetType() == NumberValueType && rightHandSide.GetType() == NumberValueType {
		return evalNumericBinaryExpr(leftHandSide.(NumberValue), rightHandSide.(NumberValue), expr.Operator)
	} else {
		return NullValue{TypedValue: TypedValue{Type: NullValueType}, Value: nil}
	}
}

func evalNumericBinaryExpr(left NumberValue, right NumberValue, operator string) RuntimeValue {
	num := 0.0

	if operator == "+" {
		num = left.Value + right.Value
	} else if operator == "-" {
		num = left.Value - right.Value
	} else if operator == "*" {
		num = left.Value * right.Value
	} else if operator == "/" {
		// TODO: implement zero check
		num = left.Value / right.Value
	} else if operator == "%" {
		// cast to int for mod
		num = float64(int(left.Value) % int(right.Value))
	}

	return MakeNumber(num)
}

func evalIdentifier(ident parser.Ident, scope *Scope) RuntimeValue {
	val := scope.LookupVariable(ident.Symbol)
	return val
}

func evalProgram(prog parser.Program, scope *Scope) RuntimeValue {
	var lastEvaluated RuntimeValue = MakeNull()

	for _, stmt := range prog.Body {
		lastEvaluated = Evaluate(stmt, scope)
	}

	return lastEvaluated
}
