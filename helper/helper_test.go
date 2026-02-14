package helper

import (
	"context"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/MayurUbarhande0/reverseproxy/models"
)

func TestGetNextPeer(t *testing.T) {
	url1, _ := url.Parse("http://localhost:8081")
	url2, _ := url.Parse("http://localhost:8082")

	pool := &models.ServerPool{
		Backends: []*models.Backend{
			{URL: url1, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(url1)},
			{URL: url2, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(url2)},
		},
	}

	// Test getting next peer
	peer := GetNextPeer(pool)
	if peer == nil {
		t.Fatal("Expected to get a peer, got nil")
	}
	if !peer.IsAlive() {
		t.Error("Expected peer to be alive")
	}
}

func TestGetNextPeer_AllDead(t *testing.T) {
	url1, _ := url.Parse("http://localhost:8081")
	url2, _ := url.Parse("http://localhost:8082")

	pool := &models.ServerPool{
		Backends: []*models.Backend{
			{URL: url1, Alive: false, ReverseProxy: httputil.NewSingleHostReverseProxy(url1)},
			{URL: url2, Alive: false, ReverseProxy: httputil.NewSingleHostReverseProxy(url2)},
		},
	}

	peer := GetNextPeer(pool)
	if peer != nil {
		t.Error("Expected nil peer when all backends are dead")
	}
}

func TestGetNextPeer_SkipsDeadBackends(t *testing.T) {
	url1, _ := url.Parse("http://localhost:8081")
	url2, _ := url.Parse("http://localhost:8082")
	url3, _ := url.Parse("http://localhost:8083")

	pool := &models.ServerPool{
		Backends: []*models.Backend{
			{URL: url1, Alive: false, ReverseProxy: httputil.NewSingleHostReverseProxy(url1)},
			{URL: url2, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(url2)},
			{URL: url3, Alive: false, ReverseProxy: httputil.NewSingleHostReverseProxy(url3)},
		},
	}

	// Should get the only alive backend
	peer := GetNextPeer(pool)
	if peer == nil {
		t.Fatal("Expected to get a peer, got nil")
	}
	if peer.URL.String() != url2.String() {
		t.Errorf("Expected to get backend 2, got %s", peer.URL.String())
	}
}

func TestPingServer(t *testing.T) {
	// Test with invalid URL (should fail)
	url, _ := url.Parse("http://localhost:99999")
	result := PingServer(url)
	if result {
		t.Error("Expected PingServer to return false for unreachable server")
	}
}

func TestGetBackendStats(t *testing.T) {
	url1, _ := url.Parse("http://localhost:8081")
	url2, _ := url.Parse("http://localhost:8082")
	url3, _ := url.Parse("http://localhost:8083")

	pool := &models.ServerPool{
		Backends: []*models.Backend{
			{URL: url1, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(url1)},
			{URL: url2, Alive: false, ReverseProxy: httputil.NewSingleHostReverseProxy(url2)},
			{URL: url3, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(url3)},
		},
	}

	stats := GetBackendStats(pool)

	if stats["total"] != 3 {
		t.Errorf("Expected total to be 3, got %v", stats["total"])
	}
	if stats["alive"] != 2 {
		t.Errorf("Expected alive to be 2, got %v", stats["alive"])
	}
	if stats["dead"] != 1 {
		t.Errorf("Expected dead to be 1, got %v", stats["dead"])
	}

	backends, ok := stats["backends"].([]map[string]interface{})
	if !ok {
		t.Fatal("Expected backends to be a slice")
	}
	if len(backends) != 3 {
		t.Errorf("Expected 3 backends in stats, got %d", len(backends))
	}
}

func TestHealthCheck_Cancellation(t *testing.T) {
	url1, _ := url.Parse("http://localhost:8081")

	pool := &models.ServerPool{
		Backends: []*models.Backend{
			{URL: url1, Alive: true, ReverseProxy: httputil.NewSingleHostReverseProxy(url1)},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	// Start health check
	go HealthCheck(ctx, pool, 100*time.Millisecond)
	
	// Let it run for a bit
	time.Sleep(250 * time.Millisecond)
	
	// Cancel context
	cancel()
	
	// Give it time to stop
	time.Sleep(200 * time.Millisecond)
	
	// If we get here without hanging, the test passes
}
