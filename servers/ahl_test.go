package servers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools"
)

func Test_ahlserver(t *testing.T) {
	ahl := &AHL{IP: "127.0.0.1", Port: "1928", URI: "test/", Auth: "123", connPool: make(map[uint64]*ahlConn), Conns: make(chan net.Conn)}
	ahl.Listen()
	go func() {
		conn, err := ahl.Accept()
		if err != nil {
			t.Error(err)
		}
		t.Log(ahl.ReadRemote(conn))
		buf := tools.MemPool.Get(65536)
		defer tools.MemPool.Put(buf)
		n, err := conn.Read(buf)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(string(buf[:n]))
		_, err = conn.Write([]byte("hello back"))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = conn.Write([]byte("hello back2"))
		if err != nil {
			t.Error(err)
			return
		}
	}()
	time.Sleep(time.Second)
	targetAddr := model.ANetAddr{Host: "bilibili.network", Port: 443, Typ: model.TCP}
	jsonB, _ := json.Marshal(targetAddr)
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:1928/test/new?action=connect&addr="+string(jsonB), nil)
	req.AddCookie(&http.Cookie{Name: "pw", Value: "123"})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
		return
	}
	b, _ := ioutil.ReadAll(resp.Body)
	t.Log(string(b))
	req, err = http.NewRequest(http.MethodPost, "http://127.0.0.1:1928/test/"+string(b), bytes.NewBuffer([]byte("send hello")))
	if err != nil {
		t.Error(err)
		return
	}
	req.AddCookie(&http.Cookie{Name: "pw", Value: "123"})
	resp2, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		t.Error(resp2.StatusCode)
		return
	}
	b2, _ := ioutil.ReadAll(resp2.Body)
	num := resp2.Header.Get("Num")
	t.Log(num)
	t.Log(string(b2))
	for num != "0" {
		req, err = http.NewRequest(http.MethodPost, "http://127.0.0.1:1928/test/"+string(b), nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.AddCookie(&http.Cookie{Name: "pw", Value: "123"})
		resp3, err := client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp3.Body.Close()

		if resp3.StatusCode != http.StatusOK && resp3.StatusCode != http.StatusAccepted {
			t.Error(resp3.StatusCode)
			return
		}
		b2, _ = ioutil.ReadAll(resp3.Body)
		num = resp3.Header.Get("Num")
		t.Log(num)
		t.Log(string(b2))
	}
}

func TestAhlConn(t *testing.T) {
	c := &ahlConn{ReadBuf: make([]connData, 0), WriteBuf: make([]connData, 0)}
	b := make([]byte, 65536)
	copy(b, []byte("123"))
	t.Log(len(b))
	c.ReadBuf = append(c.ReadBuf, connData{data: b, n: 3})
	buf := make([]byte, 65536)
	n, _ := c.Read(buf)
	t.Log(string(buf[:n]))
	c.Write([]byte("123"))
	d := c.WriteBuf[0]
	t.Log(string(d.data[:d.n]))
}
