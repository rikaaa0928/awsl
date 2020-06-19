package servers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools"
)

// NewSocks5TCP NewSocks5TCP
func NewSocks5TCP(listenHost, listenPort, tag string, id int) Socks5TCP {
	return Socks5TCP{
		IP:   listenHost,
		Port: listenPort,
		id:   id,
		tag:  tag,
	}
}

// Socks5TCP socks5
type Socks5TCP struct {
	IP   string
	Port string
	tag  string
	id   int
}

// Listen server
func (s Socks5TCP) Listen() net.Listener {
	log.Println(s.IP + ":" + s.Port)
	l, e := net.Listen("tcp", s.IP+":"+s.Port)
	if e != nil {
		panic(e)
	}
	return l
}

// ReadRemote server
func (s Socks5TCP) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	buf := tools.MemPool.Get(65536)
	defer func() {
		tools.MemPool.Put(buf)
	}()
	addr := model.ANetAddr{}
	n, d, err := tools.Receive(c, buf)
	if err != nil {
		return addr, err
	}
	if d[0] != 5 {
		return addr, errors.New("not socks5: " + string(d[:n]))
	}
	//stage1 respons
	d = []byte("\x05\x00")
	_, err = c.Write(d)
	if err != nil {
		return addr, err
	}
	//stage2 receive
	_, d, err = tools.Receive(c, buf)
	if err != nil {
		return addr, err
	}
	if d[1] != 1 {
		return addr, errors.New("socks5 connect only")
	}
	addr.Host, addr.Typ = getRemoteHost5(d)
	remotePort := getRemotePort5(d)
	addr.Port = remotePort
	//stage2 respons
	d = []byte("\x05\x00\x00\x01\x00\x00\x00\x00\xff\xff")
	_, err = c.Write(d)
	if err != nil {
		return addr, err
	}
	return addr, nil
}

// IDTag id and tag
func (s Socks5TCP) IDTag() (int, string) {
	return s.id, s.tag
}

func getRemoteHost5(data []byte) (s string, t int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%+v\n", err)
		}
	}()
	if data[3] == byte(0x03) {
		s = string(data[5 : len(data)-2])
		t = model.RAWADDR
		return
	}
	if data[3] == byte(0x01) {
		t = model.IPV4ADDR
		for i := 0; i < len(data)-6; i++ {
			s += strconv.Itoa(int(data[4+i]))
			if i != len(data)-7 {
				s += "."
			}
		}
		return
	}
	t = model.IPV6ADDR
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

func getRemotePort5(data []byte) (x int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%+v\n", err)
		}
	}()
	tt := data[len(data)-2:]
	tb := bytes.NewBuffer(tt)
	var y uint16
	binary.Read(tb, binary.BigEndian, &y)
	x = int(y)
	return
}
