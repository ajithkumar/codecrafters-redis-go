package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	// fmt.Println("Server is listening on port 8080")

	defer listener.Close()

	toWorker := make(chan interface{}, 100)
	fromWorker := make(chan string, 100)
	storage := NewStorage()
	go processMessageWorker(storage, toWorker, fromWorker)

	for {
		conn, err := listener.Accept()
		// fmt.Println("Incoming connection on port 8080")

		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}

		go handleClient(conn, toWorker, fromWorker)
	}
}

func processMessageWorker(storage *Storage, fromMain <-chan interface{}, toMain chan<- string) {
	for {
		input := <-fromMain
		command := command(input)
		params := params(input)
		outputPayload := ""

		if command == "get" {
			value, ok := storage.Get(params[0].(string))
			if ok {
				outputPayload = EncodeBulkString(value.value)
			} else {
				outputPayload = "$-1\r\n"
			}
		} else if command == "set" {
			expiryMillis := 0
			if strings.ToLower(params[2].(string)) == "px" {
				tmpExpiryMillis, err := strconv.Atoi(params[3].(string))
				if err != nil {
					expiryMillis = 0
				} else {
					expiryMillis = tmpExpiryMillis
				}
			}
			storage.Set(params[0].(string), params[1].(string), expiryMillis)
			outputPayload = "+OK\r\n"
		} else {
			outputPayload = "+OK\r\n"
		}
		toMain <- outputPayload
	}
}

func command(input interface{}) string {
	return strings.ToLower(input.([]interface{})[0].(string))
}

func params(input interface{}) []interface{} {
	return input.([]interface{})[1:]
}

func handleClient(conn net.Conn, toWorker chan<- interface{}, fromWorker <-chan string) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		// fmt.Println("Incoming data on port 8080")

		if err != nil {
			fmt.Printf("ERROR: Incoming data %s\n", err.Error())
			return
		}

		outputPayload := ""
		inputPayload := string(buf[:n])

		input, err := DecodeMessage(inputPayload)
		if err != nil {
			fmt.Printf("ERROR: Parsing message %s\n", err.Error())
			return
		}
		if command(input) == "ping" {
			outputPayload = "+PONG\r\n"
		} else if command(input) == "echo" {
			outputPayload = EncodeBulkString(input.([]interface{})[1].(string))
		} else {
			toWorker <- input
			outputPayload = <-fromWorker
		}

		outputBytes := []byte(outputPayload)
		conn.Write(outputBytes)
	}
}
