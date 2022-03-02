package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	db := Db{persistence: make(map[cmd_name]cmd_params)}
	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	waker_ch := make(chan bool)

	// go command_timeout_handler(&db, waker_ch)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handle_connections(conn, &db, waker_ch)

	}

}

func handle_connections(conn net.Conn, db *Db, waker_ch chan<- bool) {
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
			var command []string
			index := 0
			for _, cmd_part := range cmds {
				cmd_part = strings.ToLower(cmd_part)
				command = append(command, cmd_part)
				index++
				if index >= len(cmds) {
					index = 0
					resp, err := handle_command(command, db, waker_ch)
					fmt.Println("RESP: ", resp)
					if err != nil {
						conn.Write([]byte(fmt.Sprintf("- %s \r\n", err)))
					} else if resp == "<nil>" {
						conn.Write([]byte(fmt.Sprintf("$-1\r\n")))
					} else {
						conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(resp), resp)))
					}
				}
			}

		}
	}
}

func command_timeout_handler(db *Db, waker <-chan bool) {
	var key_of_least_timeout cmd_name
	for {
		least_timeout := int(^uint(0) >> 1)
		if len(db.persistence) < 1 {
			<-waker
		}
		for key, val := range db.persistence {
			if val.timeout > 0 && val.timeout < least_timeout {
				key_of_least_timeout = key
				least_timeout = val.timeout
			}

		}
		ms := least_timeout - int(time.Now().Unix()*1000)
		fmt.Println("MS: ", ms)
		select {
		case <-waker:
			{

			}
		case <-time.After(time.Duration(ms) * time.Millisecond):
			{
				fmt.Println("Deleting")
				db.mu.Lock()
				db.Remove(key_of_least_timeout)
				db.mu.Unlock()
			}
		}

	}
}

func handle_command(command []string, db *Db, waker_ch chan<- bool) (response string, err error) {
	switch command[0] {
	case "ping":
		{
			response = "PONG"
		}
	case "echo":
		{
			response = command[1]
		}
	case "set":
		{
			timeout := -1

			if len(command) > 4 {

				if command[3] != "px" {
					err = fmt.Errorf("wrong '%s' argument", command[3])
					return
				}

				timeout, err = strconv.Atoi(command[4])

				if err != nil {
					err = fmt.Errorf("cannot convert timeout because: %s", err)
					return
				}

			}

			db.mu.Lock()
			// db.Set(cmd_name(command[1]), command[2], int(time.Now().Unix()*1000)+timeout)
			db.Set(cmd_name(command[1]), command[2], timeout)
			db.mu.Unlock()

			if timeout > 0 {
				// waker_ch <- true
				go func(db *Db, timeout int64, key string) {
					time.Sleep(time.Duration(timeout))
					fmt.Println("Deleting")
					db.mu.Lock()
					db.Remove(cmd_name(key))
					db.mu.Unlock()
				}(db, int64(timeout), command[1])
			}

			fmt.Println("DB: ", db.persistence)

			response = "OK"
		}
	case "get":
		{

			response = fmt.Sprint(db.Get(cmd_name(command[1])).val)
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
				int_re := regexp.MustCompile(":[[:digit:]]+\r\n")
				int_str := int_re.FindStringSubmatch(req)[0][1:]
				ret = append(ret, strings.TrimSpace(int_str))
				req = req[len(int_str):]
			}
		default:
			{

			}
		}
	}

	return
}

type cmd_name string

type cmd_params struct {
	val     interface{}
	timeout int
}

type Db struct {
	mu          sync.Mutex
	persistence map[cmd_name]cmd_params
}

func (d *Db) Set(key cmd_name, val interface{}, timeout_ms int) bool {
	d.persistence[key] = cmd_params{val: val, timeout: timeout_ms}
	return true
}

func (d *Db) Get(key cmd_name) cmd_params {
	return d.persistence[key]
}

func (d *Db) Remove(key cmd_name) {
	delete(d.persistence, "key")
}
