package servers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"strconv"

	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools"
)

// NewSocks NewSocks
func NewSocks(listenHost string, listenPort string) SockeServer {
	return SockeServer{
		IP:   listenHost,
		Port: listenPort,
	}
}

// SockeServer socks5
type SockeServer struct {
	IP   string
	Port string
}

// Listen server
func (s SockeServer) Listen() net.Listener {
	log.Println(s.IP + ":" + s.Port)
	l, e := net.Listen("tcp", s.IP+":"+s.Port)
	if e != nil {
		panic(e)
	}
	sl := socksListenner{Listener: l, IP: s.IP}
	return sl
}

// ReadRemote server
func (s SockeServer) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	sc, ok := c.(socksConn)
	if !ok {
		uc, ok := c.(udpConn)
		if !ok {
			return model.ANetAddr{}, errors.New("invalid connection type")
		}
		return uc.remoteAddr, nil
	}
	return sc.remoteAddr, nil
}

type socksListenner struct {
	net.Listener
	IP string
}

func (l socksListenner) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, err
	}
	buf := tools.MemPool.Get(65536)
	defer tools.MemPool.Put(buf)
	n, err := conn.Read(buf)
	if err != nil {
		return conn, err
	}
	if n < 1 {
		return conn, errors.New("invalid length")
	}
	switch buf[0] {
	case 5:
		_, err = conn.Write([]byte("\x05\x00"))
		if err != nil {
			return conn, err
		}
		return socks5Stage2(conn, buf, l.IP)
	case 4:
		return socks4(conn, buf, n)
	default:
		return conn, errors.New("unsuported type")
	}
}

func socks4(conn net.Conn, buf []byte, n int) (net.Conn, error) {
	if n < 8 {
		return conn, errors.New("invalid length" + strconv.Itoa(n))
	}
	portBytes := buf[2:4]
	ipBytes := buf[4:8]
	ipStr := strconv.Itoa(int(ipBytes[0])) + "." + strconv.Itoa(int(ipBytes[1])) + "." + strconv.Itoa(int(ipBytes[2])) + "." + strconv.Itoa(int(ipBytes[3]))
	var port int16
	tb := bytes.NewBuffer(portBytes)
	err := binary.Read(tb, binary.BigEndian, &port)
	if err != nil {
		return conn, err
	}
	buf[0] = 0
	buf[1] = 90
	_, err = conn.Write(buf[:8])
	if err != nil {
		return conn, err
	}
	return socksConn{Conn: conn, remoteAddr: model.ANetAddr{Typ: model.IPV4ADDR, CMD: model.TCP, Host: ipStr, Port: int(port)}}, nil
}
