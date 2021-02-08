package rulelist

import (
	"net"
	"regexp"
	"strings"
)

type Rule interface {
	Include(string) bool
}

func New(str string) *list {
	_, cidr, _ := net.ParseCIDR(str)
	str3 := strings.Replace(str, ".", "\\.", -1)
	str3 = strings.Replace(str3, "*", ".*", -1)
	reg, _ := regexp.Compile(str3)
	return &list{str: str, reg: reg, cidr: cidr}
}

type list struct {
	str  string
	cidr *net.IPNet
	reg  *regexp.Regexp
}

func (l *list) Include(str string) bool {
	target := net.ParseIP(str)
	if target != nil && l.cidr != nil {
		if l.cidr.Contains(target) {
			return true
		}
	}
	if l.reg == nil {
		return l.str == str
	}
	return l.reg.Match([]byte(str))
}
