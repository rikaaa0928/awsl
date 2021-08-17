package udplib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rikaaa0928/awsl/global"
)

type UDPMSG struct {
	DstStr string
	SrcStr string
	Data   []byte
}

func NewUDPMSG(bb []byte, srcAddr net.Addr) (m UDPMSG, err error) {
	n := len(bb)
	minl := 4
	if n < minl {
		return m, fmt.Errorf("wrong udp data: %v", bb)
	}
	if bb[2] != 0 {
		return m, fmt.Errorf("Ignore frag: %v", bb[2])
	}
	var addr []byte
	if bb[3] == 1 {
		minl += 4
		if n < minl {
			return m, fmt.Errorf("wrong udp data: %v", bb)
		}
		addr = bb[minl-4 : minl]
	} else if bb[3] == 4 {
		minl += 16
		if n < minl {
			return m, fmt.Errorf("wrong udp data: %v", bb)
		}
		addr = bb[minl-16 : minl]
	} else if bb[3] == 3 {
		minl += 1
		if n < minl {
			return m, fmt.Errorf("wrong udp data: %v", bb)
		}
		l := bb[4]
		if l == 0 {
			return m, fmt.Errorf("wrong udp data: %v", bb)
		}
		minl += int(l)
		if n < minl {
			return m, fmt.Errorf("wrong udp data: %v", bb)
		}
		addr = bb[minl-int(l) : minl]
		addr = append([]byte{l}, addr...)
	} else {
		return m, fmt.Errorf("wrong udp data: %v", bb)
	}
	minl += 2
	if n <= minl {
		return m, fmt.Errorf("wrong udp data: %v", bb)
	}
	port := bb[minl-2 : minl]
	m.Data = bb[minl:]
	var s string
	if bb[3] == 3 {
		s = bytes.NewBuffer(addr[1:]).String()
	} else {
		s = net.IP(addr).String()
	}
	p := strconv.Itoa(int(binary.BigEndian.Uint16(port)))
	m.DstStr = net.JoinHostPort(s, p)
	m.SrcStr = srcAddr.String()
	return
}

var m map[string]*UDPConn
var lock sync.Mutex
var UDPDial *prometheus.GaugeVec
var InMap = errors.New("in map")

func init() {
	m = make(map[string]*UDPConn)
	UDPDial = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "awsl",
		Subsystem: "udplib",
		Name:      "realtime_udp_dial",
		Help:      "Number of realtime udp dial.",
	}, []string{"key"})
	prometheus.MustRegister(UDPDial)
}

type UDPConn struct {
	*net.UDPConn
	key string
}

func DialUDP(network string, laddr, raddr *net.UDPAddr) (*UDPConn, error) {
	lock.Lock()
	defer lock.Unlock()
	key := laddr.String() + "-" + raddr.String()
	if mu, ok := m[key]; ok {
		return mu, InMap
	}
	u, err := net.DialUDP(network, nil, raddr)
	ru := &UDPConn{UDPConn: u, key: key}
	if err != nil {
		return ru, err
	}
	if global.MetricsPort > 0 {
		UDPDial.With(prometheus.Labels{"key": key}).Inc()
	}
	m[key] = ru
	return ru, err
}

func (c *UDPConn) Close() error {
	lock.Lock()
	defer lock.Unlock()
	if global.MetricsPort > 0 {
		UDPDial.With(prometheus.Labels{"key": c.key}).Dec()
	}
	delete(m, c.key)
	return c.UDPConn.Close()
}
