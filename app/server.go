package main

import (
	"fmt"
	"net"
	"os"
	"strings"
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
		input, err := DecodeMessage(inputPayload)
		if err != nil {
			fmt.Printf("ERROR: Parsing message %s\n", err.Error())
			return
		}
		if strings.ToLower(input.([]interface{})[0].(string)) == "ping" {
			outputPayload = "+PONG\r\n"
		} else if strings.ToLower(input.([]interface{})[0].(string)) == "echo" {
			outputPayload = EncodeBulkString(input.([]interface{})[1].(string))
		} else {
			outputPayload = inputPayload
		}

		outputBytes := []byte(outputPayload)
		conn.Write(outputBytes)
	}
}
