package runtime

import "fmt"

func Print(Args []RuntimeValue, scope *Scope) RuntimeValue {
	for _, arg := range Args {
		str := printRuntimeValue(arg)
		fmt.Printf("%s ", str)
	}
	fmt.Println()
	return MakeNull()
}

func printRuntimeValue(val RuntimeValue) string {
	switch val.GetType() {
	case NullValueType:
		return "null"
	case BooleanValueType:
		return fmt.Sprintf("%t", val.(BooleanValue).GetValue())
	case NumberValueType:
		return fmt.Sprintf("%v", val.(NumberValue).GetValue())
	case ObjectValueType:
		obj := val.(ObjectValue)
		asStr := "{"
		for i, key := range obj.Keys() {
			asStr += fmt.Sprintf("\"%s\": %s", key, printRuntimeValue(obj.Get(key)))
			if i != len(obj.Keys())-1 {
				asStr += ", "
			}
		}
		asStr += "}"
		return asStr
	// If the runtime val is a variable, call SprintfRuntimeValue recursively on Value field
	case VariableValueType:
		return printRuntimeValue(val.(VariableValue).GetValue())
	case FunctionValueType:
		function := val.(FunctionValue)
		asStr := fmt.Sprintf("[Function: %s(", function.Name)
		for i, param := range function.Params {
			asStr += param
			if i != len(function.Params)-1 {
				asStr += ", "
			}
		}
		asStr += ")]"
		return asStr
	}
	return ""
}
