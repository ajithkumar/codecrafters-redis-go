package main

import (
	"fmt"
	"strconv"
	"strings"
)

func EncodeBulkString(value string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
}

func DecodeMessage(message string) (interface{}, error) {
	list := strings.Split(message, "\r\n")
	response, _, err := Decode(list, 0)
	return response, err
}

func Decode(list []string, off int) (interface{}, int, error) {
	if string(list[off][0]) == "*" {
		array, nextOff, err := DecodeArray(list, off)
		if err != nil {
			return array, -1, err
		}
		return array, nextOff, nil
	} else if string(list[off][0]) == "$" {
		str, nextOff, err := DecodeBulkString(list, off)
		if err != nil {
			return str, -1, err
		}
		return str, nextOff, nil
	} else if string(list[off][0]) == ":" {
		str, nextOff, err := DecodeInteger(list, off)
		if err != nil {
			return str, -1, err
		}
		return str, nextOff, nil
	}
	return nil, -1, nil
}

func DecodeArray(list []string, off int) ([]interface{}, int, error) {
	len, err := strconv.Atoi(string(list[off][1]))
	if err != nil {
		return make([]interface{}, 0), -1, err
	}
	array := make([]interface{}, len)
	nextOff := off + 1
	for i := 0; i < len; i++ {
		array[i], nextOff, err = Decode(list, nextOff)
		if err != nil {
			return make([]interface{}, 0), -1, err
		}
	}
	return array, nextOff, nil
}

func DecodeBulkString(list []string, off int) (string, int, error) {
	//TODO: handle length zero
	str := list[off+1]
	return str, off + 2, nil
}

func DecodeInteger(list []string, off int) (int, int, error) {
	intValue, err := strconv.Atoi(string(list[off][1:]))
	if err != nil {
		return 0, -1, err
	}
	return intValue, off + 1, nil
}
