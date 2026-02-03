package routes

import (
	"net/http"

	"github.com/MayurUbarhande0/reverseproxy/helper"
	"github.com/MayurUbarhande0/reverseproxy/models"
)

func LBHandler(pool *models.Serverpool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		peer := helper.GetNextPeer(pool)

		if peer != nil {

			peer.Reverse_proxy.ServeHTTP(w, r)
			return
		}

		http.Error(w, "No healthy backends available", http.StatusServiceUnavailable)
	}
}
