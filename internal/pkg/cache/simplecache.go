package cache

import (
	"sync"
	"time"
)

type KV interface {
	Get(key string) (any, bool)
	Set(key string, v any)
}

type item struct {
	v   any
	exp int64
}

type store struct {
	mu  sync.RWMutex
	m   map[string]item
	ttl time.Duration
	max int
}

type Config struct {
	TTL           int
	SweepInterval int
	MaxEntries    int
}

func NewTTLCache(cfg Config) KV {
	s := &store{m: make(map[string]item), ttl: time.Duration(cfg.TTL) * time.Second, max: cfg.MaxEntries}
	go func(interval time.Duration) {
		tick := time.NewTicker(interval)
		defer tick.Stop()
		for range tick.C {
			s.sweep()
		}
	}(time.Duration(cfg.SweepInterval) * time.Second)
	return s
}

func (s *store) Get(key string) (any, bool) {
	now := time.Now().Unix()
	s.mu.RLock()
	it, ok := s.m[key]
	s.mu.RUnlock()
	if !ok || (it.exp > 0 && it.exp < now) {
		return nil, false
	}
	return it.v, true
}

func (s *store) Set(key string, v any) {
	s.mu.Lock()
	if s.max > 0 && len(s.m) >= s.max {
		s.evictOne()
	}
	s.m[key] = item{v: v, exp: expiry(time.Now(), s.ttl)}
	s.mu.Unlock()
}

func (s *store) sweep() {
	now := time.Now().Unix()
	s.mu.Lock()
	for k, it := range s.m {
		if it.exp > 0 && it.exp < now {
			delete(s.m, k)
		}
	}
	s.mu.Unlock()
}

func (s *store) evictOne() {
	for k := range s.m {
		delete(s.m, k)
		return
	}
}

func expiry(t time.Time, ttl time.Duration) int64 {
	if ttl <= 0 {
		return 0
	}
	return t.Add(ttl).Unix()
}
