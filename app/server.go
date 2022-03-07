package main

import (
	"fmt"

	// Uncomment this block to pass the first stage
	"MREDIS/rcmd"
	"MREDIS/rdb"
	"MREDIS/resp"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	db := rdb.Db{Persistence: make(map[rdb.Key]rdb.Values)}

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	go rdb.DbElementTimeoutHandler(&db, 5)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handle_connections(conn, &db)

	}

}

func handle_connections(conn net.Conn, db *rdb.Db) {
	buffer := make([]byte, 1000)
	for {
		if _, e := conn.Read(buffer); e == nil {

			req := string(buffer)
			response := serve_request(req, db)

			conn.Write([]byte(response))
		}
	}
}

func serve_request(req string, db *rdb.Db) (response string) {
	rv := resp.RedisValue{Kind: resp.Request, V: req}
	r, e := rv.Decoder()

	if e != nil {
		fmt.Println(e)
		os.Exit(2)
	}

	return rcmd.HandleCommand(r, db)
}
