package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/Klapstuhl/fusp/pkg/config"
	"github.com/sirupsen/logrus"
)

type ReplacePath struct {
	name        string
	regex       *regexp.Regexp
	replacement string
	next        http.Handler
}

func (r *ReplacePath) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := logrus.WithFields(logrus.Fields{"middleware": r.name, "type": "ReplacePath"})

	URLPath := req.URL.EscapedPath()

	if r.regex.MatchString(URLPath) {
		req.URL.RawPath = r.regex.ReplaceAllString(URLPath, r.replacement)

		var err error
		if req.URL.Path, err = url.PathUnescape(req.URL.RawPath); err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	r.next.ServeHTTP(w, req)
}

func NewReplacePath(cfg config.ReplacePath, name string, next http.Handler) (http.Handler, error) {
	re, err := regexp.Compile(cfg.Regex)
	if err != nil {
		return nil, fmt.Errorf("middleware: %s: invalid regexp, %v", name, err)
	}
	return &ReplacePath{name: name, regex: re, replacement: cfg.Replacement, next: next}, nil
}
