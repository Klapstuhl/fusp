package middleware

import (
	"net/http"

	"github.com/Klapstuhl/fusp/pkg/config"
	"github.com/sirupsen/logrus"
)

type endpointFilter struct {
	name      string
	endpoints []string
	allowlist bool
	next      http.Handler
}

func (f *endpointFilter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	match := func() bool {
		for _, endpoint := range f.endpoints {
			if endpoint == req.URL.String() {
				return true
			}
		}
		return false
	}()

	if (f.allowlist && match) || (!f.allowlist && !match) {
		f.next.ServeHTTP(w, req)
	} else {
		f.block(w, req)
	}
}

func (f *endpointFilter) block(w http.ResponseWriter, req *http.Request) {
	logger := logrus.WithFields(logrus.Fields{"middleware": f.name, "type": "EndpointFilter"})

	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(http.StatusText(http.StatusForbidden)))
	logger.Debug("blocked ", req.Method, " to ", req.URL.String())
}

func NewEndpointFilter(cfg config.EndpointFilter, name string, next http.Handler) (http.Handler, error) {
	return &endpointFilter{name: name, endpoints: cfg.Endpoints, allowlist: cfg.Allowlist, next: next}, nil
}
