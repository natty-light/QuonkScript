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

type NullValue struct {
	TypedValue        // Type will be NullValueType
	Value      string // Value will be "null"
}

type NumberValue struct {
	TypedValue // Type will be NumberValueType
	Value      float64
}

func (n NullValue) GetType() ValueType {
	return n.Type
}

func (n NumberValue) GetType() ValueType {
	return n.Type
}
