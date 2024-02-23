package runtime

import (
	"QuonkScript/parser"
	"fmt"
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

func evalCallExpr(call parser.InternalFunctionCallExpr, scope *Scope) RuntimeValue {
	args := make([]RuntimeValue, 0)
	for _, arg := range call.Args {
		// Evaluate all args
		args = append(args, Evaluate(arg, scope))
	}
	// Get runtime value for caller function
	fn := Evaluate(call.Caller, scope)

	if fn.GetType() == InternalFunctionValueType {
		// Call function
		fn.(InternalFunctionValue).Func(args, scope)
		return MakeNull()
	} else if fn.GetType() == FunctionValueType {
		function := fn.(FunctionValue)
		// Inherits from function
		functionScope := &Scope{Parent: function.DeclarationScope, Variables: map[string]Variable{}}

		if len(args) != len(function.Params) {
			panic(fmt.Sprintf("Honk! Too few arguments for call of function %s", function.Name))
		}

		// Populate scope
		for i, param := range function.Params {
			// we have already created the runtime value
			functionScope.DeclareVariable(param, args[i], false)
		}

		// I am not sure about this
		// What about early returns
		var result RuntimeValue = MakeNull()

		// Evaluate each statement -- Check for return here maybe? How to handle breaking out of all calls
		for _, stmt := range function.Body {
			result = Evaluate(stmt, functionScope)
		}

		return result
	} else {
		panic("Honk! Cannot call non-function value")
	}

}

func evalComparisonExpr(expr parser.ComparisonExpr, scope *Scope) RuntimeValue {
	left := Evaluate(expr.Left, scope)
	right := Evaluate(expr.Right, scope)

	if left.GetType() == BooleanValueType && right.GetType() == BooleanValueType {
		return evalBooleanComparisonExpr(left.(BooleanValue), right.(BooleanValue), expr.Operator)
	} else if left.GetType() == NumberValueType && right.GetType() == NumberValueType {
		return evalNumericComparisonExpr(left.(NumberValue), right.(NumberValue), expr.Operator)
	} else {
		return MakeNull()
	}

}

func evalBooleanComparisonExpr(left BooleanValue, right BooleanValue, operator string) RuntimeValue {
	result := false
	leftVal, rightVal := left.GetValue(), right.GetValue()

	if operator == "==" {
		result = leftVal == rightVal
	} else if operator == "!=" {
		result = leftVal != rightVal
	} else if operator == "&&" {
		result = leftVal && rightVal
	} else if operator == "||" {
		result = leftVal || rightVal
	}

	return MakeBoolean(result)
}

func evalNumericComparisonExpr(left NumberValue, right NumberValue, operator string) RuntimeValue {
	result := false

	if operator == "==" {
		result = left.Value == right.Value
	} else if operator == "!=" {
		result = left.Value != right.Value
	} else if operator == ">=" {
		result = left.Value >= right.Value
	} else if operator == "<=" {
		result = left.Value <= right.Value
	} else if operator == ">" {
		result = left.Value > right.Value
	} else if operator == "<" {
		result = left.Value < right.Value
	}

	return MakeBoolean(result)
}
