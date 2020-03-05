package dialer

import (
	"log"
	"net"
)

// MultiAddr MultiAddr
type MultiAddr struct {
	Hosts     map[string][]string
	HostInUse map[string]uint
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
	hostID, ok := d.HostInUse[host]
	if !ok {
		hostID = 0
		d.HostInUse[host] = hostID
		log.Println("addr choosing", hostID, hostList[hostID])
	}
	selectHost := hostList[hostID]
	conn, err := net.Dial(network, selectHost)
	if err != nil {
		d.HostInUse[host]++
		if int(d.HostInUse[host]) >= len(hostList) {
			d.HostInUse[host] = 0
		}
		log.Println("addr choosing", d.HostInUse[host], hostList[d.HostInUse[host]])
		return conn, err
	}
	return conn, nil
}
