package rcmd

import (
	"MREDIS/rdb"
	"MREDIS/resp"
	"fmt"
	"strings"
	"time"
)

func HandleCommand(command []resp.RedisValue, db *rdb.Db) (response string) {

	cmd := strings.ToLower(command[0].V.(string))
	switch cmd {
	case "ping":
		{
			ret := resp.RedisValue{Kind: resp.SimpleString, V: "PONG"}
			response, _ = ret.Encode()
		}
	case "echo":
		{
			ret := resp.RedisValue{Kind: resp.BulkString, V: command[1].V.(string)}
			response, _ = ret.Encode()
		}
	case "set":
		{
			timeout := -1

			if len(command) > 4 {

				if command[3].V.(string) != "px" {
					ret := resp.RedisValue{Kind: resp.Error, V: fmt.Sprintf("wrong '%v' argument", command[3].V)}
					response, _ = ret.Encode()
					return
				}

				timeout = command[4].V.(int)
			}

			db.Mu.Lock()
			db.Set(rdb.Key(command[1].V.(string)), command[2].V.(string), int(time.Now().Unix()*1000)+timeout)
			db.Mu.Unlock()
			fmt.Println("DB: ", db.Persistence)

			ret := resp.RedisValue{Kind: resp.SimpleString, V: "Ok"}
			response, _ = ret.Encode()
		}
	case "get":
		{
			val := db.Get(rdb.Key(command[1].V.(string)))

			if val.Timeout <= int(time.Now().Unix()*1000) {
				db.Remove(rdb.Key(command[1].V.(string)))
				response = "<nil>"
				ret := resp.RedisValue{Kind: resp.NullBulkString, V: ""}
				response, _ = ret.Encode()
			} else {
				val := fmt.Sprint(db.Get(rdb.Key(command[1].V.(string))).Val)
				ret := resp.RedisValue{Kind: resp.BulkString, V: val}
				response, _ = ret.Encode()
			}

		}
	default:
		{
			err := fmt.Sprintf("request err: unknown command  %s", command[0].V)
			ret := resp.RedisValue{Kind: resp.Error, V: err}
			response, _ = ret.Encode()
		}

	}
	return
}
