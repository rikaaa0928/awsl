package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/servers"
	"github.com/Evi1/awsl/tools/dialer"
)

func Test_TCP(t *testing.T) {
	go func() {
		s := servers.TCP{IP: "0.0.0.0", Port: "1234", Auth: "123"}
		l := s.Listen()
		fmt.Println("listen")
		conn, err := l.Accept()
		fmt.Println(s.ReadRemote(conn))
		fmt.Println("accept")
		if err != nil {
			t.Error(err)
		}
		buf := make([]byte, 2)
		n, err := conn.Read(buf)
		fmt.Println(n, err, string(buf[:n]))
		conn.Write([]byte("b"))
	}()
	time.Sleep(2 * time.Second)
	fmt.Println("client")
	m := make(map[string][]string)
	m["lo.bilibili.network:1234"] = make([]string, 1)
	m["lo.bilibili.network:1234"][0] = "lo.bilibili.network:1234"
	c := clients.TCP{ServerHost: "lo.bilibili.network", ServerPort: "1234", Auth: "123",
		Dialer: &dialer.MultiAddr{Hosts: m, HostInUse: make(map[string]uint)}}
	conn, err := c.Dial(model.ANetAddr{Host: "www.bilibili.network", Port: 443, Typ: 1})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("dialed")
	fmt.Println(c.Verify(conn))
	time.Sleep(time.Second)
	n, err := conn.Write([]byte("a"))
	fmt.Println(n, err)
	time.Sleep(2 * time.Second)
}
