package runtime

import "QuonkScript/parser"

type ValueType int

const (
	NullValueType ValueType = iota + 1
	NumberValueType
	BooleanValueType
	ObjectValueType
	InternalFunctionValueType
	VariableValueType
	FunctionValueType
)

type RuntimeValue interface {
	GetType() ValueType
}

type TypedValue struct {
	Type ValueType
}

// Null
type NullValue struct {
	TypedValue         // Type will be NullValueType
	Value      *string // Value will be nil, pointer is to allow that
}

func (n NullValue) GetType() ValueType {
	return NullValueType
}

func (n NullValue) GetValue() *string {
	return nil
}

func MakeNull() NullValue {
	return NullValue{TypedValue: TypedValue{Type: NullValueType}, Value: nil}
}

// Number
type NumberValue struct {
	TypedValue // Type will be NumberValueType
	Value      float64
}

func (n NumberValue) GetType() ValueType {
	return NumberValueType
}

func (n NumberValue) GetValue() float64 {
	return n.Value
}

func MakeNumber(n float64) NumberValue {
	return NumberValue{TypedValue: TypedValue{Type: NumberValueType}, Value: n}
}

// Boolean

type BooleanValue struct {
	TypedValue // Type will be BoolValueType
	Value      bool
}

func (b BooleanValue) GetType() ValueType {
	return BooleanValueType
}

func (b BooleanValue) GetValue() bool {
	return b.Value
}

func MakeBoolean(b bool) BooleanValue {
	return BooleanValue{TypedValue: TypedValue{Type: BooleanValueType}, Value: b}
}

// Variable

type Variable interface {
	GetName() string
	GetValue() RuntimeValue
	IsConstant() bool
}
type VariableValue struct {
	Value    *RuntimeValue
	Constant bool
	Name     string
}

func (v VariableValue) GetName() string {
	return v.Name
}

func (v VariableValue) GetValue() RuntimeValue {
	return *v.Value
}

func (v VariableValue) IsConstant() bool {
	return v.Constant
}

func (v VariableValue) GetType() ValueType {
	return VariableValueType
}

func MakeVariable(varname string, v *RuntimeValue, constant bool) VariableValue {
	return VariableValue{Value: v, Constant: constant, Name: varname}
}

// Object

type Object interface {
	RuntimeValue
	Keys() []string
	Get(name string) RuntimeValue
	Set(name string, value RuntimeValue) RuntimeValue
}

type ObjectValue struct {
	TypedValue
	Properties map[string]RuntimeValue
}

func (o ObjectValue) GetType() ValueType {
	return o.Type
}

func (o ObjectValue) Keys() []string {
	keys := make([]string, 0)

	for k := range o.Properties {
		keys = append(keys, k)
	}
	return keys
}

func (o ObjectValue) Get(name string) RuntimeValue {
	return o.Properties[name]
}

func (o ObjectValue) Set(name string, value RuntimeValue) RuntimeValue {
	o.Properties[name] = value

	return value
}

// Functions (I am not going to distinguish from native and user defined functions)

// This is cool
type InternalFunctionCall func(Args []RuntimeValue, scope *Scope) RuntimeValue

type InternalFunction interface {
	RuntimeValue
	Call(Args []RuntimeValue, scope *Scope) RuntimeValue
}

type InternalFunctionValue struct {
	TypedValue
	Func InternalFunctionCall
}

func (f InternalFunctionValue) GetType() ValueType {
	return InternalFunctionValueType
}

func MakeFunction(call InternalFunctionCall) InternalFunctionValue {
	return InternalFunctionValue{TypedValue: TypedValue{Type: InternalFunctionValueType}, Func: call}
}

type Function interface {
	RuntimeValue
}

type FunctionValue struct {
	TypedValue
	Name             string
	Params           []string
	DeclarationScope *Scope
	Body             []parser.Stmt
}

func (f FunctionValue) GetType() ValueType {
	return FunctionValueType
}
