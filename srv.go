package tguard

import (
	"fmt"
	"time"
)

// New creates a new Service instance based on the given configuration.
func New[T any](cnf Config[T]) Service[T] {
	if cnf.Fallback == nil {
		panic(fmt.Errorf("tguard: fallback func is nil"))
	}
	if cnf.IdentityChecker == nil {
		panic(fmt.Errorf("tguard: identity checker func is nil"))
	}
	if cnf.Storage == nil {
		cnf.Storage = NewMemoryStorage()
	}
	if cnf.StorageKey == "" {
		cnf.StorageKey = "tguard_default_key"
	}
	if cnf.DefaultTTL == 0 {
		cnf.DefaultTTL = time.Minute * 5
	}
	if cnf.Interval == 0 {
		cnf.Interval = time.Minute * 1
	}
	return &guard[T]{
		checker:  cnf.IdentityChecker,
		fallback: cnf.Fallback,
		ttl:      cnf.DefaultTTL,
		interval: cnf.Interval,
		storage:  cnf.Storage,
		storeKey: cnf.StorageKey,
	}
}

// NewMemoryStorage creates a new Storage instance for in-memory storage.
func NewMemoryStorage() Storage {
	return &memoryStorage{}
}
