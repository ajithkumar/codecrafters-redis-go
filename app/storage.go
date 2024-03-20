package main

import (
	"time"
)

type StorageValue struct {
	value      string
	expiryTime int
}

type Storage struct {
	values map[string]StorageValue
}

func NewStorage() *Storage {
	return &Storage{
		values: make(map[string]StorageValue),
	}
}

func (s *Storage) Get(key string) (StorageValue, bool) {
	// TODO: Add expiry check
	value, ok := s.values[key]
	if ok {
		if value.expiryTime > 0 && value.expiryTime < int(time.Now().UnixMilli()) {
			return StorageValue{value: "", expiryTime: 0}, false
		}
		return value, ok
	} else {
		return StorageValue{value: "", expiryTime: 0}, false
	}
}

func (s *Storage) Set(key string, value string, expiry int) {
	expiryTime := 0
	if expiry > 0 {
		expiryTime = int(time.Now().UnixMilli()) + expiry
	}
	s.values[key] = StorageValue{value: value, expiryTime: expiryTime}
}
