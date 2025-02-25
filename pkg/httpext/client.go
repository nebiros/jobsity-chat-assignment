package httpext

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

func NewClient() (*http.Client, error) {
	// SEE: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/#client-timeouts
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			DualStack: true,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}

	if err := http2.ConfigureTransport(transport); err != nil {
		return nil, err
	}

	return &http.Client{
			Transport: transport,
			Timeout:   60 * time.Second,
		},
		nil
}
