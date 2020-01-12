package test

import (
	"testing"

	"github.com/Evi1/awsl/servers"
)

func TestSocks5(t *testing.T) {
	/*go func() {
		time.Sleep(time.Second)
		c, err := net.Dial("tcp", "127.0.0.1:48888")
		if err != nil {
			t.Error(err)
			return
		}
		//c.Write([]byte("abcd"))
		c.Close()
	}()*/
	t.Log("a")
	s := servers.Socke5Server{IP: "0.0.0.0", Port: "48888"}
	l := s.Listen()
	c, e := l.Accept()
	t.Log("b")
	if e != nil {
		t.Error(e)
	}
	b := make([]byte, 5)
	n, e := c.Read(b)
	t.Log(n, e, b)
	//t.Log(c.Read(make([]byte, 5)))
	if e != nil {
		t.Error(e)
	}

}
