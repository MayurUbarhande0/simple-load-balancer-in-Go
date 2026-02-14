package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := flag.Int("port", 8081, "Port to listen on")
	flag.Parse()

	hostname, _ := os.Hostname()
	address := fmt.Sprintf(":%d", *port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%d] %s %s from %s\n", *port, r.Method, r.URL.Path, r.RemoteAddr)
		response := fmt.Sprintf("Response from backend server\n"+
			"Hostname: %s\n"+
			"Port: %d\n"+
			"Path: %s\n"+
			"Time: %s\n",
			hostname, *port, r.URL.Path, time.Now().Format(time.RFC3339))
		fmt.Fprint(w, response)
	})

	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%d] Slow request received\n", *port)
		time.Sleep(5 * time.Second)
		fmt.Fprintf(w, "Slow response from backend on port %d\n", *port)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	log.Printf("Backend server starting on %s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
