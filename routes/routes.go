package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MayurUbarhande0/reverseproxy/helper"
	"github.com/MayurUbarhande0/reverseproxy/models"
)

// LBHandler creates a load balancer handler
func LBHandler(pool *models.ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peer := helper.GetNextPeer(pool)

		if peer != nil {
			log.Printf("Routing request to %s\n", peer.URL)
			peer.ReverseProxy.ServeHTTP(w, r)
			return
		}

		log.Println("No healthy backends available")
		http.Error(w, "No healthy backends available", http.StatusServiceUnavailable)
	}
}

// HealthHandler returns the health status of all backends
func HealthHandler(pool *models.ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := helper.GetBackendStats(pool)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
