package alistener

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"runtime"
	"strconv"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils"
)

// ErrUDP ErrUDP
var ErrUDP = errors.New("udp error")

// datagram is the UDP packet
type datagram struct {
	Rsv     []byte // 0x00 0x00
	Frag    byte
	Atyp    byte
	DstAddr []byte
	DstPort []byte // 2 bytes
	Data    []byte
}

func NewDatagramFromBytes(bb []byte) (*datagram, error) {
	n := len(bb)
	minl := 4
	if n < minl {
		return nil, fmt.Errorf("wrong udp data: %v", bb)
	}
	var addr []byte
	if bb[3] == 1 {
		minl += 4
		if n < minl {
			return nil, fmt.Errorf("wrong udp data: %v", bb)
		}
		addr = bb[minl-4 : minl]
	} else if bb[3] == 4 {
		minl += 16
		if n < minl {
			return nil, fmt.Errorf("wrong udp data: %v", bb)
		}
		addr = bb[minl-16 : minl]
	} else if bb[3] == 3 {
		minl += 1
		if n < minl {
			return nil, fmt.Errorf("wrong udp data: %v", bb)
		}
		l := bb[4]
		if l == 0 {
			return nil, fmt.Errorf("wrong udp data: %v", bb)
		}
		minl += int(l)
		if n < minl {
			return nil, fmt.Errorf("wrong udp data: %v", bb)
		}
		addr = bb[minl-int(l) : minl]
		addr = append([]byte{l}, addr...)
	} else {
		return nil, fmt.Errorf("wrong udp data: %v", bb)
	}
	minl += 2
	if n <= minl {
		return nil, fmt.Errorf("wrong udp data: %v", bb)
	}
	port := bb[minl-2 : minl]
	data := bb[minl:]
	d := &datagram{
		Rsv:     bb[0:2],
		Frag:    bb[2],
		Atyp:    bb[3],
		DstAddr: addr,
		DstPort: port,
		Data:    data,
	}
	return d, nil
}

func NewSocksAcceptMid(ctx context.Context, inTag string, conf map[string]interface{}) AcceptMid {
	ch := make(chan aconn.AConn, 2*runtime.NumCPU())
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
			}
		}(&closed)
		for !closed {
			buf := utils.GetMem(65536)
			defer utils.PutMem(buf)
			n, addr, err := l.ReadFromUDP(buf)
			if err != nil {
				log.Println("ReadFromUDP error: ", err)
				continue
			}
			go func(addr *net.UDPAddr, b []byte) {
				d, err := NewDatagramFromBytes(b)
				if err != nil {
					log.Println(err)
					return
				}
				if d.Frag != 0x00 {
					log.Println("Ignore frag", d.Frag)
					return
				}
				log.Println("udp data come", d)
			}(addr, buf[0:n])
		}
	}()
	return func(next Acceptor) Acceptor {
		return func(ctx context.Context) (context.Context, aconn.AConn, error) {
			select {
			case conn := <-ch:
				return ctx, conn, nil
			default:
			}
			ctx, conn, err := next(ctx)
			ctx = context.WithValue(ctx, consts.CTXInTag, inTag)
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
	io.Copy(ioutil.Discard, conn)
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
