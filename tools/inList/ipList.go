package inlist

import (
	"net"
)

// NewIPList NewIPList
func NewIPList(sl []string) IPList {
	l := IPList{List: make([]*net.IPNet, 0, len(sl))}
	for _, v := range sl {
		_, subnet, err := net.ParseCIDR(v)
		if err != nil {
			panic(err)
		}
		l.List = append(l.List, subnet)
	}
	return l
}

// IPList IPList
type IPList struct {
	List []*net.IPNet
}

// Include Include
func (l IPList) Include(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	for _, v := range l.List {
		if v.Contains(ip) {
			return true
		}
	}
	return false
}
