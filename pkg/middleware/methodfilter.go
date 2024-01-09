package middleware

import (
	"net/http"
	"strings"

	"github.com/Klapstuhl/fusp/pkg/config"
	"github.com/sirupsen/logrus"
)

type MethodFilter struct {
	name      string
	methods   []string
	blocklist bool
	next      http.Handler
}

func (f *MethodFilter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	match := func() bool {
		for _, method := range f.methods {
			if strings.ToUpper(method) == req.Method {
				return true
			}
		}
		return false
	}()

	if match {
		if f.blocklist {
			f.block(w, req)
		} else {
			f.next.ServeHTTP(w, req)
		}
	} else {
		if f.blocklist {
			f.next.ServeHTTP(w, req)
		} else {
			f.block(w, req)
		}
	}
}

func (f *MethodFilter) block(w http.ResponseWriter, req *http.Request) {
	logger := logrus.WithFields(logrus.Fields{"middleware": f.name, "type": "MethodFilter"})
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	logger.Info("blocked ", req.Method, " to ", req.URL.String())
}

func NewMethodFilter(cfg config.MethodFilter, name string, next http.Handler) (http.Handler, error) {
	return &MethodFilter{name: name, methods: cfg.Methods, blocklist: cfg.Block, next: next}, nil
}
