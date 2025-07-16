package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

var tmpDirectory = "/tmp"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	if len(os.Args) > 1 {
		tmpDirectory = os.Args[2]
	}

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// Accept any new incoming connection
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		fmt.Println()
		fmt.Println("New connection from:", conn.RemoteAddr())

		// Create a separate thread to process new incoming connection
		go handleConnection(conn)

	}
}
