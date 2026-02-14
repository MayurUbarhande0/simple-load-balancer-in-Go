package helper

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/MayurUbarhande0/reverseproxy/models"
)

// GetNextPeer returns next active peer using round-robin algorithm
func GetNextPeer(pool *models.ServerPool) *models.Backend {
	// Get next index
	next := pool.NextIndex()

	// Start from next and try each backend once
	l := len(pool.Backends)
	for i := 0; i < l; i++ {
		idx := (next + i) % l
		if pool.Backends[idx].IsAlive() {
			return pool.Backends[idx]
		}
	}
	return nil
}

// PingServer checks if backend is alive
func PingServer(url *url.URL) bool {
	timeout := 3 * time.Second
	conn, err := net.DialTimeout("tcp", url.Host, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// HealthCheck runs a routine for checking status of backends every t interval
func HealthCheck(ctx context.Context, pool *models.ServerPool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Health check stopped")
			return
		case <-ticker.C:
			for _, backend := range pool.Backends {
				alive := PingServer(backend.URL)
				backend.SetAlive(alive)
				status := "up"
				if !alive {
					status = "down"
				}
				log.Printf("Health check: %s [%s]\n", backend.URL, status)
			}
		}
	}
}

// GetBackendStats returns statistics about backend health
func GetBackendStats(pool *models.ServerPool) map[string]interface{} {
	total := len(pool.Backends)
	alive := 0

	backends := make([]map[string]interface{}, 0, total)
	for _, backend := range pool.Backends {
		isAlive := backend.IsAlive()
		if isAlive {
			alive++
		}
		backends = append(backends, map[string]interface{}{
			"url":   backend.URL.String(),
			"alive": isAlive,
		})
	}

	return map[string]interface{}{
		"total":    total,
		"alive":    alive,
		"dead":     total - alive,
		"backends": backends,
	}
}

// RetryableReverseProxy wraps the reverse proxy with retry logic
func RetryableReverseProxy(pool *models.ServerPool, maxRetries int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		attempts := 0
		for attempts < maxRetries {
			peer := GetNextPeer(pool)
			if peer == nil {
				http.Error(w, "No healthy backends available", http.StatusServiceUnavailable)
				return
			}

			// Try to serve the request
			peer.ReverseProxy.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Service unavailable after retries", http.StatusServiceUnavailable)
	}
}
