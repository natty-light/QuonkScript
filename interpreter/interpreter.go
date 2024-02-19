package interpreter

import (
	"QuonkScript/parser"
)

func Evaluate(astNode parser.Stmt) RuntimeValue {
	switch astNode.GetKind() {
	case parser.NumericLiteralNode:
		return NumberValue{Value: astNode.(parser.NumericLiteral).Value, TypedValue: TypedValue{Type: NumberValueType}}
	// Return a null by default
	case parser.NullLiteralNode:
		return NullValue{Value: "null", TypedValue: TypedValue{Type: NullValueType}}
	case parser.BinaryExprNode:
		// This typecast is save
		return evalBinaryExpr(astNode.(parser.BinaryExpr))
	case parser.ProgramNode:
		return evalProgram(astNode.(parser.Program))
	default:
		parser.PrintAST(astNode)
		panic("This NodeType has not been implemented")
	}
}

func evalBinaryExpr(expr parser.BinaryExpr) RuntimeValue {
	leftHandSide := Evaluate(expr.Left)
	rightHandSide := Evaluate(expr.Right)

	if leftHandSide.GetType() == NumberValueType && rightHandSide.GetType() == NumberValueType {
		return evalNumericBinaryExpr(leftHandSide.(NumberValue), rightHandSide.(NumberValue), expr.Operator)
	} else {
		return NullValue{TypedValue: TypedValue{Type: NullValueType}, Value: "null"}
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

	return NumberValue{TypedValue: TypedValue{Type: NumberValueType}, Value: num}
}

func evalProgram(prog parser.Program) RuntimeValue {
	var lastEvaluated RuntimeValue = NullValue{TypedValue: TypedValue{Type: NullValueType}, Value: "null"}

	for _, stmt := range prog.Body {
		lastEvaluated = Evaluate(stmt)
	}

	return lastEvaluated
}
