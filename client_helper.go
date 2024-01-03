package tguard

import (
	"context"
	"encoding/json"
)

func (g *guard[T]) unmarshal(bytes string) ([]Data[T], error) {
	var t []Data[T]
	if bytes == "" {
		return make([]Data[T], 0), nil
	}
	err := json.Unmarshal([]byte(bytes), &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (g *guard[T]) marshal(t []Data[T]) interface{} {
	bytes, _ := json.Marshal(t)
	return bytes
}

func (g *guard[T]) getData(ctx context.Context) ([]Data[T], error) {
	bytes, err := g.storage.Get(ctx, g.storeKey)
	if err != nil {
		return nil, err
	}
	return g.unmarshal(bytes)
}

func (g *guard[T]) saveData(ctx context.Context, data []Data[T]) error {
	bytes := g.marshal(data)
	return g.storage.Set(ctx, g.storeKey, bytes)
}