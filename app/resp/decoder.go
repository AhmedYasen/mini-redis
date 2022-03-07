package resp

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// func (r *RedisValue) Decoder() (ret []RedisValue, err error) {
func (r *RedisValue) Decoder() (ret []RedisValue, err error) {

	if r.Kind != Request {
		err = fmt.Errorf("You must pass a Request Kind NOT '%v'", r.Kind)
		return
	}

	arrStr := reflect.ValueOf(r.V)

	if arrStr.Kind() != reflect.String {
		err = fmt.Errorf("You must pass a string NOT '%s'", arrStr.Kind())
		return
	}

	req := fmt.Sprint(r.V)

	arr_len_re := regexp.MustCompile("\\*\\d+\r\n")
	arr_req_heads := arr_len_re.FindStringSubmatch(req)
	if len(arr_req_heads) == 0 {
		err = fmt.Errorf("request err: Wrong head format")
		return
	}
	arr_req_head := arr_req_heads[0]
	arr_len_str := strings.TrimSpace(arr_req_head)[1:]
	arr_len, err := strconv.Atoi(arr_len_str)
	if err != nil {
		err = fmt.Errorf("request err: %s", err)
		return
	}

	req = strings.TrimLeft(req, arr_req_head)

	for arr_len > 0 {
		arr_len--
		switch req[0:1] {
		case "$":
			{
				bulk_re := regexp.MustCompile("\\$[[:digit:]]+\r\n[[:alnum:]]+\r\n")
				bulk_str := bulk_re.FindStringSubmatch(req)[0]
				str := strings.TrimSpace(bulk_str[strings.Index(bulk_str, "\n"):])
				ret = append(ret, RedisValue{BulkString, str})
				req = req[len(bulk_str):]
			}
		case ":":
			{
				int_re := regexp.MustCompile(":[[:digit:]]+\r\n")
				int_str := int_re.FindStringSubmatch(req)[0][1:]
				str := strings.TrimSpace(int_str)
				num, _ := strconv.Atoi(str)
				ret = append(ret, RedisValue{Integer, num})
				req = req[len(int_str):]
			}
		default:
			{

			}
		}
	}

	return
}
