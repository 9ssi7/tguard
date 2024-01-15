package tguard

import (
	"context"
	"time"
)

type guard[T any] struct {
	checker     IdentityChecker[T]
	fallback    FallbackFunc[T]
	nowFunc     NowFunc
	ttl         time.Duration
	interval    time.Duration
	storage     Storage
	storeKey    string
	standartTTL bool
}

func (g *guard[T]) GetData(ctx context.Context) ([]Data[T], error) {
	return g.getData(ctx)
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
	maxIdx := len(current) - 1
	for i, v := range current {
		if g.checker(id, v.Original) {
			current, _ = removeSoftItem[T](current, i, maxIdx)
			g.saveData(ctx, current)
			break
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
	now := g.nowFunc().Unix()
	maxIdx := len(current) - 1
	for i, v := range current {
		if v.ExpireTime < now {
			isChanged = true
			g.fallback(v.Original)
			current, maxIdx = removeSoftItem[T](current, i, maxIdx)
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
