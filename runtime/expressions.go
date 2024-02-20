package runtime

import "QuonkScript/parser"

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
