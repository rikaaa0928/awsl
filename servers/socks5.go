package servers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"strconv"

	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools"
)

var udpPort int

func init() {
	udpPort = 65535
}

// NewSocks5 NewSocks5
func NewSocks5(listenHost string, listenPort string) Socke5Server {
	return Socke5Server{
		IP:   listenHost,
		Port: listenPort,
	}
}

// Socke5Server socks5
type Socke5Server struct {
	IP   string
	Port string
}

// Listen server
func (s Socke5Server) Listen() net.Listener {
	log.Println(s.IP + ":" + s.Port)
	l, e := net.Listen("tcp", s.IP+":"+s.Port)
	if e != nil {
		panic(e)
	}
	sl := socks5Listenner{Listener: l, listenIP: s.IP}
	return sl
}

// ReadRemote server
func (s Socke5Server) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	sc, ok := c.(socksConn)
	if !ok {
		uc, ok := c.(*udpConn)
		if !ok {
			return model.ANetAddr{}, errors.New("invalid connection type")
		}
		return uc.remoteAddr, nil
	}
	return sc.remoteAddr, nil
}

type socks5Listenner struct {
	net.Listener
	listenIP string
}

func (l socks5Listenner) Accept() (net.Conn, error) {
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
	if buf[0] != 5 {
		return conn, errors.New("unsuported type")
	}
	_, err = conn.Write([]byte("\x05\x00"))
	if err != nil {
		return conn, err
	}

	// stage2
	return socks5Stage2(conn, buf, l.listenIP)
}

func socks5Stage2(conn net.Conn, buf []byte, listenIP string) (net.Conn, error) {
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
		sc := socksConn{}
		addr.CMD = model.TCP
		sc.remoteAddr = addr
		sc.Conn = conn
		_, err = conn.Write([]byte("\x05\x00\x00\x01\x00\x00\x00\x00\xff\xff"))
		if err != nil {
			return sc, err
		}
		return sc, nil
	case 3:
		if config.Debug {
			log.Println("udp")
		}
		uc := &udpConn{}
		uc.ip = listenIP
		addr.CMD = model.UDP
		uc.remoteAddr = addr
		lp, err := uc.HandleUDP(conn)
		if err != nil {
			conn.Close()
			return nil, err
		}
		d := []byte("\x05\x00\x00\x01\x00\x00\x00\x00\xff\xff")
		copy(d[4:8], []byte(net.ParseIP(listenIP).To4()))
		bufer := new(bytes.Buffer)
		err = binary.Write(bufer, binary.LittleEndian, int16(lp))
		if err != nil {
			return uc, err
		}
		copy(d[8:], bufer.Bytes())
		_, err = conn.Write(d)
		if err != nil {
			return uc, err
		}
		return uc, nil
	default:
		return conn, errors.New("unsuported or invalid cmd : " + strconv.Itoa(int(buf[1])))
	}
}

type socksConn struct {
	net.Conn
	remoteAddr model.ANetAddr
}

type udpConn struct {
	*net.UDPConn
	remoteAddr model.ANetAddr
	tcpCon     net.Conn
	ip         string
	addr       net.Addr
	//udpListener net.Listener
}

func (c *udpConn) HandleUDP(conn net.Conn) (int, error) {
	listened := false
	c.tcpCon = conn
	p := 0
	if udpPort <= 1024 {
		udpPort = 65535
	}
	for times := 0; !listened && times < 10; times++ {
		//c.udpListener, err = net.Listen("udp", c.ip+":"+strconv.Itoa(p))
		addr, err := net.ResolveUDPAddr("udp", c.ip+":"+strconv.Itoa(udpPort))
		if err != nil {
			udpPort--
			return p, err
		}
		/*c.udpListener, err = net.ListenUDP("udp", addr)
		if err != nil {
			if config.Debug {
				log.Println("udp listen err port : " + strconv.Itoa(p) + " err : " + err.Error())
			}
			continue
		}*/
		udpConn, err := net.ListenUDP("udp", addr)
		if err != nil {
			if config.Debug {
				log.Println("udp listen err. addr : " + addr.String() + ". err : " + err.Error())
			}
			udpPort--
			continue
		}
		c.UDPConn = udpConn
		p = udpPort
		listened = true
		if config.Debug {
			log.Println("udp:listen : ", udpConn.LocalAddr(), udpConn.RemoteAddr(), udpPort)
		}
		udpPort--
	}
	if !listened {
		return p, errors.New("failed to find udp listenning port")
	}
	go func() {
		buf := tools.MemPool.Get(65535)
		defer tools.MemPool.Put(buf)
		for {
			_, err := conn.Read(buf)
			if err != nil {
				nerr, ok := err.(net.Error)
				if ok && nerr.Timeout() {
					continue
				}
				break
			}
		}
		//c.udpListener.Close()
		c.UDPConn.Close()
	}()
	return p, nil
}

func (c *udpConn) Read(b []byte) (n int, err error) {
	if config.Debug {
		log.Println("udp", "start read", c.addr, c.UDPConn.LocalAddr())
	}
	n, c.addr, err = c.UDPConn.ReadFrom(b)
	if config.Debug {
		log.Println("udp", "read", n, c.addr, err, c.UDPConn.LocalAddr())
	}
	return
}

func (c *udpConn) Write(b []byte) (n int, err error) {
	if config.Debug {
		log.Println("udp", "start write", c.addr, c.UDPConn.LocalAddr())
	}
	n, err = c.UDPConn.WriteTo(b, c.addr)
	if config.Debug {
		log.Println("udp", "write", n, c.addr, err, c.UDPConn.LocalAddr())
	}
	return
}

func (c *udpConn) Close() error {
	if config.Debug {
		log.Println("udp", "close", c.addr, c.UDPConn.LocalAddr())
	}
	if c.tcpCon != nil {
		c.tcpCon.Close()
	}
	/*if c.udpListener != nil {
		c.udpListener.Close()
	}*/
	if c.UDPConn != nil {
		return c.UDPConn.Close()
	}
	return nil
}
