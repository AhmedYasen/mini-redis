package resp

type Kind uint

const (
	Request = iota
	SimpleString
	Error
	Integer
	BulkString
	Array
	NullBulkString
	NullArray
)

type RedisValue struct {
	Kind
	V interface{}
}
