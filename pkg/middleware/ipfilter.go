package middleware

import (
	"net/http"

	"github.com/Klapstuhl/fusp/pkg/config"
	"github.com/Klapstuhl/fusp/pkg/ip"
	"github.com/sirupsen/logrus"
)

type IPFilter struct {
	name      string
	strategy  ip.Strategy
	ipList    *ip.List
	allowlist bool
	next      http.Handler
}

func (f *IPFilter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ip := f.strategy.GetIP(req)
	match := f.ipList.Contains(ip)

	if (f.allowlist && match) || (!f.allowlist && !match) {
		f.next.ServeHTTP(w, req)
	} else {
		f.block(w, req, ip)
	}
}

func (f *IPFilter) block(w http.ResponseWriter, req *http.Request, ip string) {
	logger := logrus.WithFields(logrus.Fields{"middleware": f.name, "type": "IPFilter"})
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(http.StatusText(http.StatusForbidden)))

	logger.Infof("Blocked request to %v, %s from %v", req.Host, req.URL.String(), ip)
}

func NewIPFilter(cfg config.IPFilter, name string, next http.Handler) (http.Handler, error) {
	header := cfg.Header
	if cfg.Header == "" {
		header = "X-FORWARDED-FOR"
	}

	strategy, err := ip.NewStrategy(header, cfg.Depth, cfg.ExcludedIPs)
	if err != nil {
		return nil, err					// TODO: error handling
	}

	ipList, err := ip.NewList(cfg.Addresses)
	if err != nil {	
		return nil, err					// TODO: error handling
	}

	return &IPFilter{
		name:      name,
		strategy:  strategy,
		ipList:    ipList,
		allowlist: cfg.Allowlist,
		next:      next,
	}, nil
}
