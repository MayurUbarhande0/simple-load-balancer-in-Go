package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/MayurUbarhande0/reverseproxy/helper"
	"github.com/MayurUbarhande0/reverseproxy/models"
	"github.com/MayurUbarhande0/reverseproxy/routes"
)

func main() {
	// Define your backend targets
	serverList := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}

	pool := &models.Serverpool{}

	for _, s := range serverList {
		serverUrl, err := url.Parse(s)
		if err != nil {
			fmt.Printf("Error parsing %s: %s\n", s, err)
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)

		// Add to your Backend slice
		pool.Backend = append(pool.Backend, &models.Backend{
			Url:           serverUrl,
			Alive:         true,
			Reverse_proxy: proxy,
		})
	}

	go helper.Healthcheck(pool)

	http.HandleFunc("/", routes.LBHandler(pool))

	fmt.Println("Load Balancer active on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
