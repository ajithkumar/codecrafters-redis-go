package main

import (
	"fmt"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Println("Server is listening on port 8080")

	defer listener.Close()

	for {
		// Block until we receive an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection
		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	// Ensure we close the connection after we're done
	defer conn.Close()

	// Read data
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	outputPayload := ""
	inputPayload := string(buf[:n])

	if inputPayload == "*1\r\n$4\r\nping\r\n" || inputPayload == "*1\r\n$4\r\nPING\r\n" {
		outputPayload = "+PONG\r\n"
	} else {
		// outputPayload = fmt.Sprintf("-ERR unknown command '%s'\r\n", inputPayload)
		outputPayload = inputPayload
	}

	outputBytes := []byte(outputPayload)
	conn.Write(outputBytes)
}
