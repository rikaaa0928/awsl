package test

import (
	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/servers"
	"golang.org/x/net/websocket"
	"log"
	"testing"
	"time"
)

func TestAWSL(t *testing.T) {
	c := clients.AWSL{}
	s := &servers.AWSL{
		IP:       "127.0.0.1",
		Port:     "1928",
		URI:      "wss",
		Listener: &servers.AWSListener{C: make(chan *websocket.Conn), IP: "127.0.0.1", Port: "1928"},
		Key:      GetTestPath() + "/server.key",
		Cert:     GetTestPath() + "/server.crt",
	}
	go func() {
		l := s.Listen()
		for {
			conn, _ := l.Accept()
			go func() {
				t.Log(conn.Write([]byte("qweasdzxc")))
			}()
		}
	}()
	log.Print("start")
	time.Sleep(time.Second)
	conn, err := c.Dial("127.0.0.1", "1928")
	if err != nil {
		t.Error(err)
	}
	t.Log(conn.LocalAddr(), conn.RemoteAddr())
	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b[:n]))
	b = make([]byte, 1024)
	n, err = conn.Read(b)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b[:n]))
	time.Sleep(5 * time.Second)
}
