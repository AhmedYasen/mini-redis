package resp

import (
	"fmt"
	"reflect"
)

func (r *RedisValue) Encode() (ret string, err error) {
	switch r.Kind {
	case SimpleString:
		{
			return r.encodeSimpleString()
		}
	case Error:
		{
			return r.encodeError()
		}
	case Integer:
		{
			return r.encodeInteger()
		}
	case BulkString:
		{
			return r.encodeBulkString()
		}
	case Array:
		{
			return r.encodeArray()
		}
	case NullBulkString:
		{
			ret = "$-1\r\n"
			return
		}
	case NullArray:
		{
			ret = "*-1\r\n"
			return
		}

	}

	return
}

func (r *RedisValue) encodeArray() (ret string, err error) {
	arr := reflect.ValueOf(r.V)

	if arr.Kind() != reflect.Array {
		err = fmt.Errorf("The value '%s' is not of kind '%v'", r.V, r.Kind)
		return
	}
	ret = fmt.Sprintf("*%d\r\n", arr.Len())
	for loopIndex := 0; loopIndex < arr.Len(); loopIndex++ {
		rval, ok := arr.Index(loopIndex).Interface().(RedisValue)
		if ok == false {
			err = fmt.Errorf("The value '%s' is not of kind '%v'", rval.V, rval.Kind)
			return
		}

		r, e := rval.Encode()
		if e != nil {
			err = e
			return
		}

		ret += r
	}

	return
}

func (r *RedisValue) encodeSimpleString() (ret string, err error) {
	ss := reflect.ValueOf(r.V)
	if ss.Kind() != reflect.String {
		err = fmt.Errorf("The value '%s' is not of kind '%v'", r.V, r.Kind)
		return
	}

	ret = fmt.Sprintf("+%s\r\n", r.V)

	return
}

func (r *RedisValue) encodeError() (ret string, err error) {
	ss := reflect.ValueOf(r.V)
	if ss.Kind() != reflect.String {
		err = fmt.Errorf("The value '%s' is not of kind '%v'", r.V, r.Kind)
		return
	}

	ret = fmt.Sprintf("-%s\r\n", r.V)

	return
}

func (r *RedisValue) encodeInteger() (ret string, err error) {
	ss := reflect.ValueOf(r.V)
	if ss.Kind() != reflect.Int {
		err = fmt.Errorf("The value '%s' is not of kind '%v'", r.V, r.Kind)
		return
	}

	ret = fmt.Sprintf(":%s\r\n", r.V)

	return
}

func (r *RedisValue) encodeBulkString() (ret string, err error) {
	ss := reflect.ValueOf(r.V)
	if ss.Kind() != reflect.String {
		err = fmt.Errorf("The value '%s' is not of kind '%v'", r.V, r.Kind)
		return
	}

	val := fmt.Sprint(r.V)
	ret = fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)

	return
}
