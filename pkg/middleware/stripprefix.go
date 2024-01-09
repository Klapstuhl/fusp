package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Klapstuhl/fusp/pkg/config"
	"github.com/sirupsen/logrus"
)

type StripPrefix struct {
	name  string
	regex []*regexp.Regexp
	next  http.Handler
}

func (s *StripPrefix) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := logrus.WithFields(logrus.Fields{"middleware": s.name, "type": "StripPrefix"})

	for _, exp := range s.regex {
		match := exp.FindStringSubmatch(req.URL.Path)
		if len(match) > 0 && len(match[0]) > 0 {
			prefix := match[0]
			if !strings.HasPrefix(req.URL.Path, prefix) {
				continue
			}

			logger.Info("stripping prefix ", prefix, "from ", req.URL.Path)

			req.URL.Path = strings.Replace(req.URL.Path, prefix, "", 1)
			if !strings.HasPrefix(req.URL.Path, "/") {
				req.URL.Path = "/" + req.URL.Path
			}
			req.RequestURI = req.URL.RequestURI()
			s.next.ServeHTTP(w, req)
		}
	}
	s.next.ServeHTTP(w, req)
}

func NewStripPrefix(cfg config.StripPrefix, name string, next http.Handler) (http.Handler, error) {
	re := []*regexp.Regexp{}

	for _, str := range cfg.Prefix {
		exp, err := regexp.Compile(str)
		if err != nil {
			return nil, fmt.Errorf("middleware: %s: invalid regexp, %v", name, err)
		}
		re = append(re, exp)
	}

	return &StripPrefix{name: name, regex: re, next: next}, nil
}
