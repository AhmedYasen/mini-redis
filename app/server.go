package main

import (
	"fmt"
	"strings"

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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		buffer := make([]byte, 1000)

		if _, e := conn.Read(buffer); e == nil {
			go func(cmd string) {

				cmd = strings.ToLower(cmd[:strings.Index(cmd, "\n")])
				switch cmd {
				case "ping":
					{
						conn.Write([]byte(fmt.Sprint("+PONG\r\n")))
					}
				default:
					{

					}
				}

			}(fmt.Sprintf("%s", buffer))

		} else {
			break
		}

	}

}
