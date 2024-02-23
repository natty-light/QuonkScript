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
	var value RuntimeValue = MakeNull()
	if declaration.Value == nil {
		return scope.DeclareVariable(declaration.Identifier, value, declaration.Constant)
	} else {
		exprVal := *declaration.Value
		if exprVal.GetKind() == parser.IdentifierNode {
			value = Evaluate(exprVal, scope)
		} else {
			fn := exprVal.(parser.FunctionDeclaration)
			value = FunctionValue{Name: fn.Name, Params: fn.Params, DeclarationScope: scope, Body: fn.Body, TypedValue: TypedValue{Type: FunctionValueType}}
		}
	}
	return scope.DeclareVariable(declaration.Identifier, value, declaration.Constant)

}

func evalFunctionDeclaration(declaration parser.FunctionDeclaration, scope *Scope) RuntimeValue {
	function := FunctionValue{Name: declaration.Name, Params: declaration.Params, DeclarationScope: scope, Body: declaration.Body, TypedValue: TypedValue{Type: FunctionValueType}} // intializes with zero value for all fields

	return scope.DeclareVariable(function.Name, function, true)
}
