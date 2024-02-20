package runtime

import (
	"QuonkScript/set"
	"fmt"
)

type Scope struct {
	Parent    *Scope // pointer to env so it can be null
	Variables map[string]RuntimeValue
	Constants *set.Set
}

func (s *Scope) DeclareVariable(varname string, value RuntimeValue, constant bool) RuntimeValue {
	if s.Variables[varname] != nil {
		panic(fmt.Sprintf("Cannot redeclare variable %s", varname))
	}

	s.Variables[varname] = value
	if constant {
		s.Constants.Add(varname)
	}

	return value
}

func (e *Scope) AssignVariable(varname string, value RuntimeValue) RuntimeValue {
	scope := e.Resolve(varname)

	if scope.Constants.Includes(varname) {
		panic(fmt.Sprintf("Cannot reassign constant variable %s", varname))
	}

	scope.Variables[varname] = value
	return value
}

func (e *Scope) LookupVariable(varname string) RuntimeValue {
	scope := e.Resolve(varname)
	return scope.Variables[varname]
}

func (s *Scope) Resolve(varname string) *Scope {
	if s.Variables[varname] != nil {
		return s
	}

	if s.Parent == nil {
		panic(fmt.Sprintf("Cannot resolve variable %s", varname))
	}

	// since Parent is a pointer to allow for nil, Scope will always be a pointer
	return s.Parent.Resolve(varname)
}
