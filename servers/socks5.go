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
	return l
}

// ReadRemote server
func (s Socke5Server) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	buf := tools.MemPool.Get(65536)
	defer func() {
		tools.MemPool.Put(buf)
	}()
	n, d, err := tools.Receive(c, buf)
	if err != nil {
		return model.ANetAddr{}, err
	}
	if !bytes.Equal([]byte{d[0]}, []byte("\x05")) {
		return model.ANetAddr{}, errors.New("not socks5: " + string(d[:n]))
	}
	//stage1 respons
	d = []byte("\x05\x00")
	_, err = c.Write(d)
	if err != nil {
		return model.ANetAddr{}, err
	}
	//stage2 receive
	_, d, err = tools.Receive(c, buf)
	if err != nil {
		return model.ANetAddr{}, err
	}
	addr := model.ANetAddr{}
	addr.Host, addr.Typ = getRemoteHost(d)
	remotePort := getRemotePort(d)
	addr.Port = remotePort
	//stage2 respons
	d = []byte("\x05\x00\x00\x01\x00\x00\x00\x00\xff\xff")
	_, err = c.Write(d)
	if err != nil {
		return model.ANetAddr{}, err
	}
	return addr, nil
}

func getRemoteHost(data []byte) (s string, t int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%+v\n", err)
		}
	}()
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
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%+v\n", err)
		}
	}()
	tt := data[len(data)-2:]
	t := []byte{0x00, 0x00}
	t = append(t, tt...)
	tb := bytes.NewBuffer(t)
	var y int32
	binary.Read(tb, binary.BigEndian, &y)
	x = int(y)
	return
}
