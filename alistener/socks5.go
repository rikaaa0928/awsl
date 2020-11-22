package alistener

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"runtime"
	"strconv"
	"time"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils"
	"github.com/rikaaa0928/awsl/utils/superlib"
)

// ErrUDP ErrUDP
var ErrUDP = errors.New("udp error")

func NewSocksAcceptMid(ctx context.Context, inTag string, conf map[string]interface{}) AcceptMid {
	ch := make(chan aconn.AConn, 2*runtime.NumCPU())
	if udpI, ok := conf["udp"]; ok {
		if udpB, ok := udpI.(bool); ok && udpB {
			go func() {
				closed := false
				udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(conf["host"].(string), strconv.Itoa(int(conf["port"].(float64)))))
				if err != nil {
					panic(err)
				}
				l, err := net.ListenUDP("udp", udpAddr)
				if err != nil {
					panic(err)
				}
				log.Println("udp listen ", udpAddr)
				go func(closed *bool) {
					select {
					case <-ctx.Done():
						*closed = true
						l.Close()
						close(ch)
					}
				}(&closed)
				buf := utils.GetMem(65536)
				defer utils.PutMem(buf)
				for !closed {
					n, srcAddr, err := l.ReadFromUDP(buf)
					if err != nil {
						log.Println("ReadFromUDP error: ", err)
						continue
					}
					go func(addr *net.UDPAddr, b []byte) {
						defer func() {
							recover()
						}()
						u, err := superlib.NewUDPMSG(b, srcAddr)
						if err != nil {
							log.Println(err)
							return
						}
						endAddr, err := net.ResolveUDPAddr("udp", u.DstStr)
						if err != nil {
							log.Println(err)
							return
						}
						conn := &udpConnWrapper{UDPConn: l, srcAddr: srcAddr, endAddr: endAddr, readChan: make(chan []byte, 1)}
						go func() {
							defer func() {
								recover()
							}()
							t := time.NewTimer(time.Minute * 10)
							select {
							case <-t.C:
								close(conn.readChan)
							case conn.readChan <- u.Data:
							}
						}()
						ch <- conn
					}(srcAddr, buf[0:n])
				}
			}()
		}
	}
	return func(next Acceptor) Acceptor {
		return func(ctx context.Context) (context.Context, aconn.AConn, error) {
			tch := make(chan struct {
				ctx  context.Context
				conn aconn.AConn
				err  error
			}, 2*runtime.NumCPU())
			go func() {
				ctx, conn, err := next(ctx)
				tch <- struct {
					ctx  context.Context
					conn aconn.AConn
					err  error
				}{ctx, conn, err}
			}()
			var conn aconn.AConn
			var err error
			ctx = context.WithValue(ctx, consts.CTXInTag, inTag)
			select {
			case conn, ok := <-ch:
				if !ok {
					log.Println("udp chan closed")
					return ctx, nil, ErrUDP
				}
				// ctx = context.WithValue(ctx, consts.CTXSuperType, msg.T)
				// ctx = context.WithValue(ctx, consts.CTXSuperData, msg.MSG)
				// //ctx = superlib.SetID(ctx, msg.ID)
				return ctx, conn, nil
			case tm := <-tch:
				ctx = tm.ctx
				conn = tm.conn
				err = tm.err
			}
			//ctx, conn, err := next(ctx)
			if err != nil {
				return ctx, nil, err
			}
			buf := utils.GetMem(65536)
			defer utils.PutMem(buf)
			n, err := conn.Read(buf)
			if err != nil {
				conn.Close()
				return ctx, nil, err
			}
			if n < 1 {
				conn.Close()
				return ctx, nil, errors.New("invalid length")
			}
			switch buf[0] {
			case 5:
				_, err = conn.Write([]byte("\x05\x00"))
				if err != nil {
					conn.Close()
					return ctx, nil, err
				}
				return socks5Stage2(ctx, conn, buf)
			case 4:
				return socks4(ctx, conn, buf, n)
			default:
				conn.Close()
				return ctx, nil, errors.New("unsuported type")
			}
		}
	}
}

func socks4(ctx context.Context, conn aconn.AConn, buf []byte, n int) (context.Context, aconn.AConn, error) {
	if n < 8 {
		conn.Close()
		return ctx, nil, errors.New("invalid length" + strconv.Itoa(n))
	}
	portBytes := buf[2:4]
	ipBytes := buf[4:8]
	host := strconv.Itoa(int(ipBytes[0])) + "." + strconv.Itoa(int(ipBytes[1])) + "." + strconv.Itoa(int(ipBytes[2])) + "." + strconv.Itoa(int(ipBytes[3]))
	var port int16
	tb := bytes.NewBuffer(portBytes)
	err := binary.Read(tb, binary.BigEndian, &port)
	if err != nil {
		conn.Close()
		return ctx, nil, err
	}
	buf[0] = 0
	buf[1] = 90
	_, err = conn.Write(buf[:8])
	if err != nil {
		conn.Close()
		return ctx, nil, err
	}
	conn.SetEndAddr(aconn.NewAddr(host, int(port), "tcp"))
	return ctx, conn, nil
}

func socks5Stage2(ctx context.Context, conn aconn.AConn, buf []byte) (context.Context, aconn.AConn, error) {

	n, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		return ctx, nil, err
	}
	if n < 2 {
		conn.Close()
		return ctx, nil, errors.New("invalid length")
	}
	switch buf[1] {
	case 1:
		host, _ := getRemoteHost5(buf[:n])
		port := getRemotePort5(buf[:n])
		_, err = conn.Write([]byte("\x05\x00\x00\x01\x00\x00\x00\x00\xff\xff"))
		if err != nil {
			conn.Close()
			return ctx, nil, err
		}
		conn.SetEndAddr(aconn.NewAddr(host, port, "tcp"))
		return ctx, conn, nil
	case 3:
		host, _ := getRemoteHost5(buf[:n])
		port := getRemotePort5(buf[:n])
		log.Println("udp from ", host, port, " to udp ", conn.LocalAddr().String())
		ctx, err = udp(ctx, conn)
		if err != nil {
			log.Println("udp error: ", err)
		}
		return ctx, nil, ErrUDP
	default:
		conn.Close()
		return ctx, nil, errors.New("unsuported or invalid cmd : " + strconv.Itoa(int(buf[1])))
	}
}

func udp(ctx context.Context, conn aconn.AConn) (context.Context, error) {
	var err error
	if err != nil {
		conn.Close()
		return ctx, err
	}
	a, addr, port, err := parseAddress(conn.LocalAddr().String())
	if err != nil {
		conn.Close()
		return ctx, err
	}
	rep := []byte("\x05\x00\x00")
	rep = append(rep, a)
	rep = append(rep, addr...)
	rep = append(rep, port...)
	_, err = conn.Write(rep)
	if err != nil {
		conn.Close()
		return ctx, err
	}
	go io.Copy(ioutil.Discard, conn)
	return ctx, err
}

func getRemoteHost5(data []byte) (s string, t int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("socks5 get remote host error: %+v\n", err)
		}
	}()
	if data[3] == byte(0x03) {
		s = string(data[5 : len(data)-2])
		t = 0
		return
	}
	if data[3] == byte(0x01) {
		t = 4
		// for i := 0; i < len(data)-6; i++ {
		// 	s += strconv.Itoa(int(data[4+i]))
		// 	if i != len(data)-7 {
		// 		s += "."
		// 	}
		// }
		s = net.IP(data[4 : len(data)-2]).String()
		return
	}
	t = 6
	s = net.IP(data[4 : len(data)-2]).String()
	// for i := 0; i < len(data)-7; i += 2 {
	// 	s += strconv.FormatInt(int64(data[4+i]), 16)
	// 	s += fmt.Sprintf("%02x", int(data[5+i]))
	// 	if i != len(data)-8 {
	// 		s += ":"
	// 	}
	// }
	return
}

func getRemotePort5(data []byte) (x int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("socks5 get remote port erro: r%+v\n", err)
		}
	}()
	tt := data[len(data)-2:]
	tb := bytes.NewBuffer(tt)
	var y uint16
	binary.Read(tb, binary.BigEndian, &y)
	x = int(y)
	return
}

// parseAddress format address x.x.x.x:xx to raw address.
// addr contains domain length
func parseAddress(address string) (a byte, addr []byte, port []byte, err error) {
	var h, p string
	h, p, err = net.SplitHostPort(address)
	if err != nil {
		return
	}
	ip := net.ParseIP(h)
	if ip4 := ip.To4(); ip4 != nil {
		a = 1
		addr = []byte(ip4)
	} else if ip6 := ip.To16(); ip6 != nil {
		a = 4
		addr = []byte(ip6)
	} else {
		a = 3
		addr = []byte{byte(len(h))}
		addr = append(addr, []byte(h)...)
	}
	i, _ := strconv.Atoi(p)
	port = make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(i))
	return
}

type udpConnWrapper struct {
	*net.UDPConn
	srcAddr  *net.UDPAddr
	endAddr  net.Addr
	readChan chan []byte
}

func (c *udpConnWrapper) Write(b []byte) (n int, err error) {
	// uMsg := superlib.UDPMSG{}
	// err = json.Unmarshal(b, &uMsg)
	// if err != nil {
	// 	log.Println("udp unmarshal error", err)
	// 	return
	// }
	a, addr, port, err := parseAddress(c.endAddr.String())
	if err != nil {
		log.Println(err)
		return
	}
	data2Write := newDatagram(a, addr, port, b)
	n, err = c.UDPConn.WriteToUDP(data2Write.Bytes(), c.srcAddr)
	if err == nil {
		log.Println("udp write to ", c.srcAddr, n, string(b))
	}
	return
}

func (c *udpConnWrapper) Read(b []byte) (n int, err error) {
	t := time.NewTimer(time.Minute * 10)
	select {
	case <-t.C:
		defer func() {
			recover()
		}()
		close(c.readChan)
		return 0, errors.New("time out")
	case d, ok := <-c.readChan:
		if !ok {
			return 0, errors.New("udp readChan closed")
		}
		n = copy(b, d)
		return
	}
}

func (c *udpConnWrapper) Close() error {
	defer func() {
		recover()
	}()
	close(c.readChan)
	return nil
}

func (c *udpConnWrapper) EndAddr() net.Addr {
	return c.endAddr
}

func (c *udpConnWrapper) SetEndAddr(addr net.Addr) {
	c.endAddr = addr
}

type datagram struct {
	Rsv     []byte // 0x00 0x00
	Frag    byte
	Atyp    byte
	DstAddr []byte
	DstPort []byte // 2 bytes
	Data    []byte
}

func newDatagram(atyp byte, dstaddr []byte, dstport []byte, data []byte) *datagram {
	if atyp == 3 {
		dstaddr = append([]byte{byte(len(dstaddr))}, dstaddr...)
	}
	return &datagram{
		Rsv:     []byte{0x00, 0x00},
		Frag:    0x00,
		Atyp:    atyp,
		DstAddr: dstaddr,
		DstPort: dstport,
		Data:    data,
	}
}

func (d *datagram) Bytes() []byte {
	b := make([]byte, 0)
	b = append(b, d.Rsv...)
	b = append(b, d.Frag)
	b = append(b, d.Atyp)
	b = append(b, d.DstAddr...)
	b = append(b, d.DstPort...)
	b = append(b, d.Data...)
	return b
}
