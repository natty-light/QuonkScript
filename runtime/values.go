package runtime

type ValueType int

const (
	NullValueType ValueType = iota + 1
	NumberValueType
	BooleanValueType
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
