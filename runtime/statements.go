package runtime

import "QuonkScript/parser"

func evalProgram(prog parser.Program, scope *Scope) RuntimeValue {
	var lastEvaluated RuntimeValue = MakeNull()

	for _, stmt := range prog.Body {
		lastEvaluated = Evaluate(stmt, scope)
	}

	return lastEvaluated
}

func evalVarDeclaration(declaration parser.VarDeclaration, scope *Scope) RuntimeValue {
	var value RuntimeValue

	if declaration.Value == nil {
		value = MakeNull()
	} else {
		value = Evaluate(*declaration.Value, scope)
	}

	return scope.DeclareVariable(declaration.Identifier, value, declaration.Constant)
}

func evalBranchStatement(stmt parser.BranchStmt, scope *Scope) RuntimeValue {
	var lastEvaluated RuntimeValue = MakeNull()

	childScope := &Scope{Parent: scope, Variables: map[string]Variable{}}

	condition := Evaluate(stmt.Condition, scope).(BooleanValue)

	if condition.Type == BooleanValueType {
		if condition.GetValue() {
			for _, bodyStmt := range stmt.Body {
				lastEvaluated = Evaluate(bodyStmt, childScope)
			}
		} else {
			for _, elseStmt := range stmt.Else {
				lastEvaluated = Evaluate(elseStmt, childScope)
			}
		}
	}

	return lastEvaluated
}

func evalFunctionDeclaration(declaration parser.FunctionDeclaration, scope *Scope) RuntimeValue {
	function := FunctionValue{Name: declaration.Name, Params: declaration.Params, DeclarationScope: scope, Body: declaration.Body, TypedValue: TypedValue{Type: FunctionValueType}} // intializes with zero value for all fields

	return scope.DeclareVariable(function.Name, function, true)
}
