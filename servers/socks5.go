package servers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Evi1/awsl/tools"
	"log"
	"net"
	"strconv"
)

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
	return l
}

// ReadRemote server
func (s Socke5Server) ReadRemote(c net.Conn) (ANetAddr, error) {
	log.Printf("read")
	n, d := tools.Receive(c)
	log.Printf("read: %d : %#v\n", n, d[:n])
	if !bytes.Equal([]byte{d[0]}, []byte("\x05")) {
		return ANetAddr{}, errors.New("not socks5")
	}
	//stage1 respons
	d = []byte("\x05\x00")
	n = tools.Send(c, d)
	log.Printf("write:%#v:n="+strconv.Itoa(n)+"\n", d)
	//stage2 receive
	n, d = tools.Receive(c)
	log.Printf("read: %d : %#v\n", n, d[:n])
	addr := ANetAddr{}
	addr.Host, addr.Typ = getRemoteHost(d)
	remotePort := getRemotePort(d)
	log.Println("rh=" + addr.Host)
	log.Printf("rp=%v\n", remotePort)
	addr.Port = remotePort
	//stage2 respons
	d = []byte("\x05\x00\x00\x01\x00\x00\x00\x00\xff\xff")
	n = tools.Send(c, d)
	log.Printf("write:%#v:n="+strconv.Itoa(n)+"\n", d)
	return addr, nil
}

func getRemoteHost(data []byte) (s string, t int) {
	if data[3] == byte(0x03) {
		s = string(data[5 : len(data)-2])
		t = 1
		return
	}
	if data[3] == byte(0x01) {
		t = 4
		for i := 0; i < len(data)-6; i++ {
			t := data[4+i]
			x := int(t)
			s += strconv.Itoa(x)
			if i != len(data)-7 {
				s += "."
			}
		}
		return
	}
	t = 6
	s += "["
	for i := 0; i < len(data)-7; i += 2 {
		s += strconv.FormatInt(int64(data[4+i]), 16)
		s += fmt.Sprintf("%02x", int(data[5+i]))
		if i != len(data)-8 {
			s += ":"
		}
	}
	s += "]"
	return
}

func getRemotePort(data []byte) (x int) {
	tt := data[len(data)-2:]
	t := []byte{0x00, 0x00}
	t = append(t, tt...)
	tb := bytes.NewBuffer(t)
	var y int32
	binary.Read(tb, binary.BigEndian, &y)
	x = int(y)
	return
}
