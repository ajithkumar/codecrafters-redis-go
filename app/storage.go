package main

type StorageValue struct {
	value string
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
	return value, ok
}

func (s *Storage) Set(key string, value string) {
	// expiryTime := time.Now().UTC().UnixMilli() + 100
	s.values[key] = StorageValue{value: value}
}
