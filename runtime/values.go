package runtime

type ValueType int

const (
	NullValueType ValueType = iota + 1
	NumberValueType
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
	return n.Type
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
	return n.Type
}

func MakeNumber(n float64) NumberValue {
	return NumberValue{TypedValue: TypedValue{Type: NumberValueType}, Value: n}
}
