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
