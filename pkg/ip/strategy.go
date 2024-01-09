package ip

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

type Strategy interface {
	GetIP(r *http.Request) string
}

func NewStrategy(header string, depth int, excludedIPs []string) (Strategy, error) {
	if depth > 0 {
		return &DepthStrategy{header: header, depth: depth}, nil
	}

	if len(excludedIPs) > 0 {
		list, err := NewList(excludedIPs)
		if err != nil {
			return nil, fmt.Errorf("ExcludedStrategy %s", err.Error())
		}
		return &ExcludedStrategy{header: header, list: list}, nil
	}

	return &RemoteAddrStrategy{}, nil
}

type RemoteAddrStrategy struct{}

func (s *RemoteAddrStrategy) GetIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

type DepthStrategy struct {
	header string
	depth  int
}

func (s *DepthStrategy) GetIP(r *http.Request) string {
	addresses := strings.Split(r.Header.Get(s.header), ",")

	if len(addresses) < s.depth {
		return ""
	}
	return strings.TrimSpace(addresses[len(addresses)-s.depth])
}

type ExcludedStrategy struct {
	header string
	list   *List
}

func (s *ExcludedStrategy) GetIP(r *http.Request) string {
	if s.list == nil {
		return ""
	}

	addresses := strings.Split(r.Header.Get(s.header), ",")

	for i := len(addresses) - 1; i >= 0; i-- {
		address := strings.TrimSpace(addresses[i])
		if !s.list.Contains(address) {
			return address
		}
	}
	return ""
}
