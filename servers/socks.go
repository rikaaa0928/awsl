package servers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools"
)

// NewSocks NewSocks
func NewSocks(listenHost, listenPort, tag string, id int) Socks {
	return Socks{
		IP:   listenHost,
		Port: listenPort,
		id:   id,
		tag:  tag,
	}
}

// Socks socks5
type Socks struct {
	IP   string
	Port string
	tag  string
	id   int
}

// Listen server
func (s Socks) Listen() net.Listener {
	l, e := net.Listen("tcp", s.IP+":"+s.Port)
	if e != nil {
		panic(e)
	}
	sl := socksListenner{Listener: l, IP: s.IP, udpConn: make(chan *socksUDPConn)}
	return sl
}

// ReadRemote server
func (s Socks) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	sc, ok := c.(socksConn)
	if !ok {
		return model.ANetAddr{}, errors.New("invalid connection type")
	}
	return sc.GetRemote(), nil
}

// IDTag id and tag
func (s Socks) IDTag() (int, string) {
	return s.id, s.tag
}

type socksListenner struct {
	udpConn chan *socksUDPConn
	net.Listener
	IP string
}

func (l socksListenner) Accept() (net.Conn, error) {
	select {
	case conn := <-l.udpConn:
		return conn, nil
	default:
	}
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
		return l.socks5Stage2(conn, buf)
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
	return socksTCPConn{Conn: conn, remoteAddr: model.ANetAddr{Typ: model.IPV4ADDR, CMD: model.TCP, Host: ipStr, Port: int(port)}}, nil
}

func (l socksListenner) socks5Stage2(conn net.Conn, buf []byte) (net.Conn, error) {
	addr := model.ANetAddr{}
	n, err := conn.Read(buf)
	if err != nil {
		return conn, err
	}
	if n < 2 {
		return conn, errors.New("invalid length")
	}
	addr.Host, addr.Typ = getRemoteHost5(buf[:n])
	remotePort := getRemotePort5(buf[:n])
	addr.Port = remotePort

	switch buf[1] {
	case 1:
		sc := socksTCPConn{}
		addr.CMD = model.TCP
		sc.remoteAddr = addr
		sc.Conn = conn
		_, err = conn.Write([]byte("\x05\x00\x00\x01\x00\x00\x00\x00\xff\xff"))
		if err != nil {
			return sc, err
		}
		return sc, nil
	case 3:
		if config.UDP {
			go func() {
				err = l.newSocksUDPConn(conn)
				if err != nil {
					log.Println("udp new connction err.", err)
					conn.Close()
				}
			}()
		} else {
			conn.Close()
		}
		return nil, ErrUDP
	default:
		return conn, errors.New("unsuported or invalid cmd : " + strconv.Itoa(int(buf[1])))
	}
}

type socksConn interface {
	GetRemote() model.ANetAddr
}

type socksTCPConn struct {
	net.Conn
	remoteAddr model.ANetAddr
}

func (c socksTCPConn) GetRemote() model.ANetAddr {
	return c.remoteAddr
}

func (l socksListenner) newSocksUDPConn(conn net.Conn) error {
	listenAddr, err := net.ResolveUDPAddr("udp", l.IP+":0")
	if err != nil {
		return err
	}
	udpConn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return err
	}
	log.Println("udp listen", udpConn.LocalAddr(), udpConn.RemoteAddr(), "\n udp with tcp", conn.LocalAddr(), conn.RemoteAddr())

	ipport := strings.Split(udpConn.LocalAddr().String(), ":")
	if len(ipport) != 2 {
		return errors.New("split length error : " + strconv.Itoa(len(ipport)))
	}
	rep := []byte("\x05\x00\x00\x01\x00\x00\x00\x00\xff\xff")
	copy(rep[4:8], []byte(net.ParseIP(ipport[0]).To4()))
	portBuf := new(bytes.Buffer)
	port, err := strconv.Atoi(ipport[1])
	if err != nil {
		return err
	}
	err = binary.Write(portBuf, binary.LittleEndian, uint16(port))
	if err != nil {
		return err
	}
	copy(rep[8:10], portBuf.Bytes())
	_, err = conn.Write(rep)
	if err != nil {
		return err
	}

	buf := tools.MemPool.Get(65536)
	defer tools.MemPool.Put(buf)
	sUDP := &socksUDPConn{UDPConn: udpConn, close: false}
	go func() {
		for {
			tinyBuf := tools.MemPool.Get(10)
			defer tools.MemPool.Put(tinyBuf)
			n, err := conn.Read(tinyBuf)
			if n != 0 {
				log.Printf("udp : tcp conn read %d bytes \n", n)
			}
			if err != nil {
				ne, ok := err.(net.Error)
				if ok && ne.Timeout() {
					continue
				}
				break
			}
		}
		sUDP.close = true
		sUDP.Close()
	}()
	for !sUDP.close {
		n, err := udpConn.Read(buf)
		if config.Debug {
			log.Printf("udp init read %d bytes, % x \n", n, buf[:n])
		}
		if err != nil {
			if config.Debug {
				log.Println("udp init read error : ", err)
			}
			return err
		}
		sUDP.remoteAddr = model.ANetAddr{CMD: model.UDP}
		l.udpConn <- sUDP
	}
	return nil
}

type socksUDPConn struct {
	*net.UDPConn
	reqAddr    net.Addr
	close      bool
	remoteAddr model.ANetAddr
}

func (c *socksUDPConn) Read(b []byte) (int, error) {
	if c.UDPConn == nil {
		return 0, errors.New("udp read : nil udpConn")
	}
	n, addr, e := c.ReadFrom(b)
	c.reqAddr = addr
	return n, e
}

func (c *socksUDPConn) Write(b []byte) (int, error) {
	if c.UDPConn == nil {
		return 0, errors.New("udp write : nil udpConn")
	}
	return c.WriteTo(b, c.reqAddr)
}

func (c *socksUDPConn) Close() error {
	c.close = true
	if c.UDPConn != nil {
		return c.UDPConn.Close()
	}
	return nil
}

func (c *socksUDPConn) GetRemote() model.ANetAddr {
	return c.remoteAddr
}
