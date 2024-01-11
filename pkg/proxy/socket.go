package proxy

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
)

func NewSocket(path string) (http.Handler, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	var transport = &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", path)
		},
		DisableCompression: true,
	}

	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = " "
	}

	socket := httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		socket.ServeHTTP(w, req)
	})
	return handler, nil
}
