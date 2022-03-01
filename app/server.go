package main

import (
	"fmt"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, 1000)
	for {

		if _, e := conn.Read(buffer); e == nil {
			// go func(cmd string) {
			// 	switch strings.ToLower(cmd) {
			// 	case "ping":
			// 		{
			// 			conn.Write([]byte(fmt.Sprint("+PONG\r\n")))
			// 		}
			// 	default:
			// 		{

			// 		}
			// 	}
			// }(fmt.Sprintf("%s", buffer))
			conn.Write([]byte(fmt.Sprint("+PONG\r\n")))

		} else {
			break
		}

	}

}
