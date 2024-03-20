package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Println("Server is listening on port 8080")

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		fmt.Println("Incoming connection on port 8080")

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		fmt.Println("Incoming data on port 8080")

		if err != nil {
			fmt.Printf("ERROR: Incoming data %s\n", err.Error())
			return
		}

		outputPayload := ""
		inputPayload := string(buf[:n])
		fmt.Printf("Incoming data on port 8080 : %s\n", inputPayload)

		if inputPayload == "*1\r\n$4\r\nping\r\n" || inputPayload == "*1\r\n$4\r\nPING\r\n" {
			outputPayload = "+PONG\r\n"
		} else {
			// outputPayload = fmt.Sprintf("-ERR unknown command '%s'\r\n", inputPayload)
			outputPayload = inputPayload
		}

		outputBytes := []byte(outputPayload)
		conn.Write(outputBytes)
	}
}
