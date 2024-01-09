package middleware

import (
	"net/http"
	"strings"

	"github.com/Klapstuhl/fusp/pkg/config"
	"github.com/sirupsen/logrus"
)

type AddPrefix struct {
	name   string
	prefix string
	next   http.Handler
}

func (a *AddPrefix) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := logrus.WithFields(logrus.Fields{"middleware": a.name, "type": "AddPrefix"})
	oldURLPath := req.URL.Path

	newURlPath := a.prefix + oldURLPath
	if !strings.HasPrefix(newURlPath, "/") {
		newURlPath = "/" + newURlPath
	}
	req.URL.Path = newURlPath
	logger.Debug("added prefix: ", a.prefix, " to ", oldURLPath)

	if req.URL.RawPath != "" {
		oldURLRawPath := req.URL.RawPath
		newURLRawPath := a.prefix + oldURLRawPath
		if !strings.HasPrefix(newURLRawPath, "/") {
			newURLRawPath = "/" + newURLRawPath
		}
		req.URL.RawPath = newURLRawPath
	}

	req.RequestURI = req.URL.RequestURI()
	a.next.ServeHTTP(w, req)
}

func NewAddPrefix(cfg config.AddPrefix, name string, next http.Handler) (http.Handler, error) {
	return &AddPrefix{name: name, prefix: cfg.Prefix, next: next}, nil
}
