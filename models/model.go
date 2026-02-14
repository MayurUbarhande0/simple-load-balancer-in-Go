package models

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

// Backend represents a backend server
type Backend struct {
	URL          *url.URL
	Alive        bool
	Mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// SetAlive sets the backend's alive status thread-safely
func (b *Backend) SetAlive(alive bool) {
	b.Mux.Lock()
	defer b.Mux.Unlock()
	b.Alive = alive
}

// IsAlive returns the backend's alive status thread-safely
func (b *Backend) IsAlive() bool {
	b.Mux.RLock()
	defer b.Mux.RUnlock()
	return b.Alive
}

// ServerPool holds information about reachable backends
type ServerPool struct {
	Backends []*Backend
	current  uint64
}

// NextIndex atomically increases the counter and returns an index
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.Backends)))
}
