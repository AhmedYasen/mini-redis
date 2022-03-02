package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	db := Db{persistence: make(map[string]interface{})}
	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handle_connections(conn, &db)

	}

}

func handle_connections(conn net.Conn, db *Db) {
	buffer := make([]byte, 1000)
	for {
		if _, e := conn.Read(buffer); e == nil {

			req := string(buffer)
			cmds, err := parse_request(req)

			fmt.Println("CMDS: ", cmds)

			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			var command [3]string
			index := 0
			for _, cmd_part := range cmds {
				cmd_part = strings.ToLower(cmd_part)
				switch cmd_part {
				case "ping":
					{
						conn.Write([]byte("+PONG\r\n"))
					}
				default:
					{
						command[index] = cmd_part
						index++
						if index >= len(cmds) {
							fmt.Println(command)
							index = 0
							resp, err := handle_command(&command, db)
							if err != nil {
								conn.Write([]byte(fmt.Sprintf("- %s \r\n", err)))
							} else {
								conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(resp), resp)))
							}
						}

					}
				}
			}

		}
	}
}

func handle_command(command *[3]string, db *Db) (response string, err error) {
	switch command[0] {
	case "echo":
		{
			response = command[1]
		}
	case "set":
		{
			db.mu.Lock()
			db.Set(command[1], command[2])
			db.mu.Unlock()
			response = "OK"
		}
	case "get":
		{
			response = fmt.Sprint(db.Get(command[1]))
		}
	default:
		{
			err = fmt.Errorf("request err: unknown command  %s", command[1])
		}

	}
	return
}

func parse_request(req string) (ret []string, err error) {
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
		fmt.Println("REQ BEF SWTCH", req)
		switch req[0:1] {
		case "$":
			{
				bulk_re := regexp.MustCompile("\\$[[:digit:]]+\r\n[[:alnum:]]+\r\n")
				bulk_str := bulk_re.FindStringSubmatch(req)[0]
				str := strings.TrimSpace(bulk_str[strings.Index(bulk_str, "\n"):])
				ret = append(ret, str)
				req = req[len(bulk_str):]
			}
		case ":":
			{

			}
		default:
			{

			}
		}
	}

	return
}

type Db struct {
	mu          sync.Mutex
	persistence map[string]interface{}
}

func (d *Db) Set(key string, val interface{}) bool {
	d.persistence[key] = val
	return true
}

func (d *Db) Get(key string) (val interface{}) {
	return d.persistence[key]
}
