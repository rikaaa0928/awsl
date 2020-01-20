package test

import (
	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/servers"
	"log"
	"net"
	"testing"
	"time"
)

func TestAWSL(t *testing.T) {
	c := clients.AWSL{
		ServerHost: "127.0.0.1",
		ServerPort: "1928",
	}
	s := &servers.AWSL{
		IP:       "127.0.0.1",
		Port:     "1928",
		URI:      "wss",
		Listener: &servers.AWSListener{C: make(chan net.Conn), IP: "127.0.0.1", Port: "1928"},
		Key:      GetTestPath() + "/server.key",
		Cert:     GetTestPath() + "/server.crt",
	}
	go func() {
		l := s.Listen()
		for {
			conn, _ := l.Accept()
			go func() {
				t.Log(s.ReadRemote(conn))

				t.Log(conn.Write([]byte("qweasdzxc")))

				b := make([]byte, 1024)
				n, err := conn.Read(b)
				if err != nil {
					t.Error(err)
				}
				t.Log(string(b[:n]))
				t.Log(conn.Close())
			}()
		}
	}()
	log.Print("start")
	time.Sleep(time.Second)
	conn, err := c.Dial(servers.ANetAddr{
		Typ:  1,
		Host: "bilibili.network",
		Port: 443,
	})
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b[:n]))
	t.Log(conn.Write([]byte("hahaha")))
	time.Sleep(3 * time.Second)
}
