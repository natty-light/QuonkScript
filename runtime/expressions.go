package runtime

import (
	"QuonkScript/parser"
)

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

func evalAssignmentExpr(expr parser.VarAssignemntExpr, scope *Scope) RuntimeValue {
	if expr.Assignee.GetKind() != parser.IdentifierNode { // will check for objects in the future
		panic("Honk! Attempt to assign value to something other than an identifier")
	}
	varname := expr.Assignee.(parser.Ident).Symbol
	// This cast is safe as we are not implementing objects yet
	return scope.AssignVariable(varname, Evaluate(expr.Value, scope))
}

func evalObjectExpr(object parser.ObjectLiteral, scope *Scope) RuntimeValue {
	obj := ObjectValue{TypedValue: TypedValue{Type: ObjectValueType}, Properties: make(map[string]RuntimeValue)}
	var val RuntimeValue
	for _, propertyLiteral := range object.Properties {
		key := propertyLiteral.Key
		value := propertyLiteral.Value
		// { key }
		if value == nil {
			val = scope.LookupVariable(key)
		} else {
			// Dereference pointer
			val = Evaluate(*value, scope)
		}
		obj.Set(key, val)
	}
	return obj
}

func evalInternalFuncCallExpr(call parser.InternalFunctionCallExpr, scope *Scope) RuntimeValue {
	args := make([]RuntimeValue, 0)
	for _, arg := range call.Args {
		// Evaluate all args
		args = append(args, Evaluate(arg, scope))
	}
	// Get runtime value for caller function
	fn := Evaluate(call.Caller, scope).(InternalFunctionValue)

	if fn.GetType() != InternalFunctionValueType {
		panic("Honk! Cannot call non-function value")
	}

	// Call function
	fn.Func(args, scope)
	return MakeNull()
}
