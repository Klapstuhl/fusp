package ip

import (
	"errors"
	"net"
)

type List struct {
	ips    []*net.IP
	ipNets []*net.IPNet
}

func NewList(addresses []string) (*List, error) {
	if len(addresses) == 0 {
		return nil, errors.New("empty address slice")
	}

	list := &List{}

	for _, address := range addresses {
		if ipAddr := net.ParseIP(address); ipAddr != nil {
			list.ips = append(list.ips, &ipAddr)
			continue
		}

		_, ipNet, err := net.ParseCIDR(address)
		if err != nil {
			return nil, err
		}
		list.ipNets = append(list.ipNets, ipNet)
	}

	return list, nil
}

func (l *List) Contains(address string) bool {
	ipAddr := net.ParseIP(address)
	if ipAddr == nil {
		return false
	}

	for _, ip := range l.ips {
		if ip.Equal(ipAddr) {
			return true
		}
	}

	for _, net := range l.ipNets {
		if net.Contains(ipAddr) {
			return true
		}
	}
	
	return false
}