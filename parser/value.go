package parser

type ValueType int

const (
	VALUE_TYPE_UNKNOWN ValueType = iota
	VALUE_TYPE_INTEGER
	VALUE_TYPE_FLOAT
)

var VALUE_TYPE_MAP = map[ValueType]string{
	VALUE_TYPE_UNKNOWN: "VALUE_TYPE_UNKNOWN",
	VALUE_TYPE_INTEGER: "VALUE_TYPE_INTEGER",
	VALUE_TYPE_FLOAT:   "VALUE_TYPE_FLOAT",
}

func (vt ValueType) String() string {
	return VALUE_TYPE_MAP[vt]
}
