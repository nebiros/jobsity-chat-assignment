package httpext

import (
	"net/http"
	"time"
)

func NewServer(address string, mux http.Handler) *http.Server {
	return &http.Server{
		Addr:    address,
		Handler: mux,
		// SEE: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/#server-timeouts
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
