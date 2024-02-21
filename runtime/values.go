package runtime

type ValueType int

const (
	NullValueType ValueType = iota + 1
	NumberValueType
	BooleanValueType
	ObjectValueType
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

func MakeVariable(varname string, v *RuntimeValue, constant bool) VariableValue {
	return VariableValue{Value: v, Constant: constant, Name: varname}
}

// Object

type Object interface {
	GetType()
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
