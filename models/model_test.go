package models

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"testing"
)

func TestBackend_SetAlive(t *testing.T) {
	backend := &Backend{
		Alive: true,
	}

	backend.SetAlive(false)
	if backend.Alive != false {
		t.Error("Expected Alive to be false")
	}

	backend.SetAlive(true)
	if backend.Alive != true {
		t.Error("Expected Alive to be true")
	}
}

func TestBackend_IsAlive(t *testing.T) {
	backend := &Backend{
		Alive: true,
	}

	if !backend.IsAlive() {
		t.Error("Expected IsAlive to return true")
	}

	backend.Alive = false
	if backend.IsAlive() {
		t.Error("Expected IsAlive to return false")
	}
}

func TestBackend_ConcurrentAccess(t *testing.T) {
	backend := &Backend{
		Alive: true,
	}

	var wg sync.WaitGroup
	iterations := 100

	// Concurrent writes
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(val bool) {
			defer wg.Done()
			backend.SetAlive(val)
		}(i%2 == 0)
	}

	// Concurrent reads
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = backend.IsAlive()
		}()
	}

	wg.Wait()
}

func TestServerPool_NextIndex(t *testing.T) {
	url1, _ := url.Parse("http://localhost:8081")
	url2, _ := url.Parse("http://localhost:8082")

	pool := &ServerPool{
		Backends: []*Backend{
			{URL: url1, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(url1)},
			{URL: url2, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(url2)},
		},
	}

	// Test that NextIndex cycles through backends
	indices := make(map[int]bool)
	for i := 0; i < 10; i++ {
		idx := pool.NextIndex()
		if idx < 0 || idx >= len(pool.Backends) {
			t.Errorf("NextIndex returned out of range index: %d", idx)
		}
		indices[idx] = true
	}

	// Verify we hit different backends
	if len(indices) < 2 {
		t.Error("NextIndex should cycle through different backends")
	}
}

func TestServerPool_NextIndex_Concurrent(t *testing.T) {
	url1, _ := url.Parse("http://localhost:8081")
	url2, _ := url.Parse("http://localhost:8082")

	pool := &ServerPool{
		Backends: []*Backend{
			{URL: url1, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(url1)},
			{URL: url2, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(url2)},
		},
	}

	var wg sync.WaitGroup
	iterations := 1000

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			idx := pool.NextIndex()
			if idx < 0 || idx >= len(pool.Backends) {
				t.Errorf("NextIndex returned out of range index: %d", idx)
			}
		}()
	}

	wg.Wait()
}
