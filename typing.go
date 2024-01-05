package tguard

import (
	"context"
	"time"
)

// Data is the fundamental structure for managing data over time.
type Data[T interface{}] struct {
	Original   T     `json:"original"`    // The original data.
	ExpireTime int64 `json:"expire_time"` // The expiration time for the data.
}

// FallbackFunc defines the fallback mechanism for timed-out data.
type FallbackFunc[T interface{}] func(data T)

// IdentityChecker is a function type used to authenticate data.
type IdentityChecker[T interface{}] func(id string, data T) bool

// Service is the interface that defines the core functionality.
type Service[T interface{}] interface {
	Start(ctx context.Context, data T, ttl ...time.Duration) error
	Cancel(ctx context.Context, id string) error
	Connect(ctx context.Context)
}

// Storage defines where the data will be stored.
type Storage interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}) error
	Exist(ctx context.Context, key string) (bool, error)
}

// Config defines the configuration for the service.
type Config[T interface{}] struct {
	Fallback        FallbackFunc[T]    // The fallback function for timed-out data.
	IdentityChecker IdentityChecker[T] // The function to check the identity of data.
	DefaultTTL      time.Duration      // The default time-to-live for data. default: 5 minutes. (time.Minute * 5)
	Interval        time.Duration      // The interval for checking timeouts. default: 1 minute. (time.Minute * 1)
	Storage         Storage            // The storage backend. default: memoryStorage
	StorageKey      string             // The key under which data will be stored. default: "tguard_default_key"
	WithStandardTTL bool               // Whether to use all data with the same expiration time. default: true
}
