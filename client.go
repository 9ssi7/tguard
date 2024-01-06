package tguard

import (
	"context"
	"time"
)

type guard[T any] struct {
	checker     IdentityChecker[T]
	fallback    FallbackFunc[T]
	ttl         time.Duration
	interval    time.Duration
	storage     Storage
	storeKey    string
	standartTTL bool
}

func (g *guard[T]) Start(ctx context.Context, data T, ttls ...time.Duration) error {
	ttl := g.ttl
	if len(ttls) > 0 && !g.standartTTL {
		ttl = ttls[0]
	}
	cacheItem := Data[T]{
		Original:   data,
		ExpireTime: time.Now().Add(ttl).Unix(),
	}
	current, err := g.getData(ctx)
	if err != nil {
		return err
	}
	current = append(current, cacheItem)
	return g.saveData(ctx, current)
}

func (g *guard[T]) Cancel(ctx context.Context, id string) error {
	current, err := g.getData(ctx)
	if err != nil {
		return err
	}
	for i, v := range current {
		if g.checker(id, v.Original) {
			if i == len(current)-1 {
				current = current[:i]
			} else {
				current = append(current[:i], current[i+1:]...)
			}
			return g.saveData(ctx, current)
		}
	}
	return nil
}

func (g *guard[T]) check(ctx context.Context) {
	current, err := g.getData(ctx)
	if err != nil {
		return
	}
	isChanged := false
	now := time.Now().Unix()
	for i, v := range current {
		if v.ExpireTime < now {
			isChanged = true
			g.fallback(v.Original)
			current = append(current[:i], current[i+1:]...)
		}
	}
	if isChanged {
		_ = g.saveData(ctx, current)
	}
}

func (g *guard[T]) Connect(ctx context.Context) {
	ticker := time.NewTicker(g.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			g.check(ctx)
		}
	}
}
