package helper

import (
	"net"
	"net/url"
	"time"

	"github.com/MayurUbarhande0/reverseproxy/models"
)

func GetNextPeer(pool *models.Serverpool) *models.Backend {
	pool.Mux.Lock()
	pool.Current = (pool.Current + 1) % len(pool.Backend)
	defer pool.Mux.Unlock()
	for range len(pool.Backend) {
		if pool.Backend[pool.Current].Alive == true {
			return pool.Backend[pool.Current]
		} else {
			pool.Current = (pool.Current + 1) % len(pool.Backend)
			continue
		}
	}
	return nil
}
func Pingserver(url *url.URL) bool {
	conn, err := net.DialTimeout("tcp", url.Host, 5*time.Second)

	if err != nil {
		return false
	}
	conn.Close()
	return true
}
func Healthcheck(s *models.Serverpool) {
	for {
		time.Sleep(1 * time.Second)
		for _, b := range s.Backend {
			alive := Pingserver(b.Url)
			b.Mux.Lock()
			b.Alive = alive
			b.Mux.Unlock()
		}
	}
}
