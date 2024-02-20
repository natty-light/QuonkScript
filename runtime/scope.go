package runtime

import (
	"fmt"
)

type Scope struct {
	Parent    *Scope // pointer to env so it can be null
	Variables map[string]Variable
}

func (s *Scope) DeclareVariable(varname string, value RuntimeValue, constant bool) RuntimeValue {
	if s.Variables[varname] != nil {
		panic(fmt.Sprintf("Honk! Cannot redeclare variable %s", varname))
	}

	s.Variables[varname] = VariableValue{Value: &value, Constant: constant, Name: varname}

	return value
}

func (e *Scope) AssignVariable(varname string, value RuntimeValue) RuntimeValue {
	scope := e.Resolve(varname)

	variable := scope.Variables[varname]

	if variable != nil && variable.IsConstant() {
		panic(fmt.Sprintf("Honk! Cannot assign to constant variable %s", varname))
	}

	// If we make it past the panic, we know the variable will not be constant
	scope.Variables[varname] = VariableValue{Value: &value, Name: varname, Constant: false}
	return value
}

func (e *Scope) LookupVariable(varname string) RuntimeValue {
	scope := e.Resolve(varname)
	return scope.Variables[varname].GetValue()
}

func (s *Scope) Resolve(varname string) *Scope {
	if s.Variables[varname] != nil {
		return s
	}

	if s.Parent == nil {
		panic(fmt.Sprintf("Honk! Cannot resolve variable %s", varname))
	}

	// since Parent is a pointer to allow for nil, Scope will always be a pointer
	return s.Parent.Resolve(varname)
}
