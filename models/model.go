package models

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	Url   *url.URL
	Alive bool
	Mux   sync.Mutex

	Reverse_proxy *httputil.ReverseProxy
}

type Serverpool struct {
	Backend []*Backend
	Mux     sync.Mutex
	Current int
}
