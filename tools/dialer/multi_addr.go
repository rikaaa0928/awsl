package dialer

import (
	"log"
	"net"
	"sync"
)

// MultiAddr MultiAddr
type MultiAddr struct {
	Hosts     map[string][]string
	HostInUse map[string]uint
	lock      sync.Mutex
}

// Dial Dial
func (d *MultiAddr) Dial(network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	host = net.JoinHostPort(host, port)
	hostList, ok := d.Hosts[host]
	if !ok || len(hostList) == 0 {
		log.Println("addr not in map", addr)
		return net.Dial(network, addr)
	}
	d.lock.Lock()
	hostID, ok := d.HostInUse[host]
	if !ok {
		hostID = 0
		d.HostInUse[host] = hostID
		log.Println("addr choosing", hostID, hostList[hostID])
	}
	d.lock.Unlock()
	selectHost := hostList[hostID]
	conn, err := net.Dial(network, selectHost)
	if err != nil {
		d.lock.Lock()
		defer d.lock.Unlock()
		if hostID == d.HostInUse[host] {
			d.HostInUse[host]++
			log.Println("addr choosing", d.HostInUse[host], hostList[d.HostInUse[host]])
		}
		if int(d.HostInUse[host]) >= len(hostList) {
			d.HostInUse[host] = 0
		}
		return conn, err
	}
	return conn, nil
}
