package main

import (
	"reflect"
	"testing"
)

func TestDecodeMessage(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test Case 1: Array of bulk strings",
			args: args{message: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"},
			want: []interface{}{"hello", "world"},
		},
		{
			name: "Test Case 2: Array of integers",
			args: args{message: "*3\r\n:1\r\n:2\r\n:3\r\n"},
			want: []interface{}{1, 2, 3},
		},
		{
			name: "Test Case 3: Array of integers and bulk strings",
			args: args{message: "*5\r\n:1\r\n:2\r\n:3\r\n$5\r\nhello\r\n$5\r\nworld\r\n"},
			want: []interface{}{1, 2, 3, "hello", "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeMessage(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeBulkString(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case 1: bulk strings",
			args: args{value: "hello"},
			want: "$5\r\nhello\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeBulkString(tt.args.value); got != tt.want {
				t.Errorf("EncodeBulkString() = %v, want %v", got, tt.want)
			}
		})
	}
}
