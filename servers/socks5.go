package servers

/*
import (
	"errors"
	"log"
	"net"
	"strconv"

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

		return nil, ErrUDP
	default:
		return conn, errors.New("unsuported or invalid cmd : " + strconv.Itoa(int(buf[1])))
	}
}

type socksConn struct {
	net.Conn
	remoteAddr model.ANetAddr
}
*/
