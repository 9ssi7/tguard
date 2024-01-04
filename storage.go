package tguard

import (
	"context"
)

type memoryStorage struct {
	data map[string]interface{}
}

func (m *memoryStorage) Get(ctx context.Context, key string) (string, error) {
	if m.data == nil {
		return "", nil
	}
	if _, ok := m.data[key]; ok {
		isByteArr, ok := m.data[key].([]byte)
		if ok {
			return string(isByteArr), nil
		}
		return "", nil
	}
	return "", nil
}

func (m *memoryStorage) Set(ctx context.Context, key string, value interface{}) error {
	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	m.data[key] = value
	return nil
}

func (m *memoryStorage) Exist(ctx context.Context, key string) (bool, error) {
	if m.data == nil {
		return false, nil
	}
	if _, ok := m.data[key]; ok {
		return true, nil
	}
	return false, nil
}
