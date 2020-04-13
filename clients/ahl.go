package clients

// import (
// 	"bytes"
// 	"crypto/tls"
// 	"encoding/json"
// 	"errors"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net"
// 	"net/http"
// 	"strconv"
// 	"sync"
// 	"time"

// 	"github.com/Evi1/awsl/config"
// 	"github.com/Evi1/awsl/model"
// 	"github.com/Evi1/awsl/tools"
// )

// // NewAHL NewAHL
// func NewAHL(serverHost, serverPort, uri, auth string) *AHL {
// 	return &AHL{
// 		ServerHost: serverHost,
// 		ServerPort: serverPort,
// 		URI:        uri,
// 		Auth:       auth,
// 	}
// }

// // AHL AHL
// type AHL struct {
// 	ServerHost string
// 	ServerPort string
// 	URI        string
// 	Auth       string
// }

// // Dial Dial
// func (c *AHL) Dial(addr model.ANetAddr) (net.Conn, error) {
// 	transCfg := &http.Transport{
// 		TLSClientConfig: &tls.Config{
// 			InsecureSkipVerify: config.GetConf().NoVerify,
// 			ServerName:         c.ServerHost,
// 		},
// 	}
// 	conn := &ahlConn{Addr: addr,
// 		ServerURL: "https://" + c.ServerHost + ":" + c.ServerPort + "/" + c.URI + "/",
// 		Client:    &http.Client{Transport: transCfg},
// 		Cookie:    &http.Cookie{Name: "pw", Value: c.Auth},
// 		Data:      make(chan ahlData),
// 		CloseChan: make(chan int8),
// 		leftlock:  sync.Mutex{}}
// 	go func() {
// 		for {
// 			select {
// 			case <-conn.CloseChan:
// 				break
// 			default:
// 			}
// 			if conn.ReadLeft == 0 {
// 				time.Sleep(time.Second)
// 				continue
// 			}
// 			conn.readMore()
// 		}
// 	}()
// 	return conn, nil
// }

// // Verify Verify
// func (c *AHL) Verify(conn net.Conn) error {
// 	ac, ok := conn.(*ahlConn)
// 	if !ok {
// 		return errors.New("wrong type")
// 	}
// 	jsonB, err := json.Marshal(ac.Addr)
// 	if err != nil {
// 		return err
// 	}
// 	req, err := http.NewRequest(http.MethodGet, ac.ServerURL+"new?action=connect&addr="+string(jsonB), nil)
// 	if err != nil {
// 		return err
// 	}
// 	req.AddCookie(ac.Cookie)
// 	client := ac.Client
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return errors.New(strconv.Itoa(resp.StatusCode))
// 	}
// 	b, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	ac.ServerURL += string(b)
// 	log.Println("client : " + ac.ServerURL + " connected. addr : " + string(jsonB))
// 	return nil
// }

// type ahlConn struct {
// 	Addr      model.ANetAddr
// 	ServerURL string
// 	Client    *http.Client
// 	Cookie    *http.Cookie
// 	Data      chan ahlData
// 	CloseChan chan int8
// 	DataLeft  []byte
// 	leftlock  sync.Mutex
// 	LeftErr   error
// 	ReadLeft  uint64
// }

// type ahlData struct {
// 	data []byte
// 	err  error
// }

// func (c *ahlConn) Read(b []byte) (n int, err error) {
// 	defer func() { log.Println(c.ServerURL, "client read ", n, err) }()
// 	c.leftlock.Lock()
// 	if len(c.DataLeft) != 0 {
// 		log.Println(c.ServerURL, "read from left")
// 		if len(b) >= len(c.DataLeft) {
// 			copy(b, c.DataLeft)
// 			l := len(c.DataLeft)
// 			c.DataLeft = nil
// 			c.leftlock.Unlock()
// 			return l, c.LeftErr
// 		}
// 		copy(b, c.DataLeft[:len(b)])
// 		c.DataLeft = c.DataLeft[len(b):]
// 		c.leftlock.Unlock()
// 		return len(b), nil
// 	}
// 	c.leftlock.Unlock()
// 	select {
// 	case data := <-c.Data:
// 		if len(b) >= len(data.data) {
// 			copy(b, data.data)
// 			return len(data.data), data.err
// 		}
// 		copy(b, data.data)
// 		c.DataLeft = data.data[len(b):]
// 		c.LeftErr = data.err
// 		log.Println(c.ServerURL, "not read all")
// 		return len(b), nil
// 	case <-time.After(time.Minute):
// 		return 0, tools.ErrTimeout
// 	case <-c.CloseChan:
// 	}
// 	return 0, io.EOF
// }

// func (c *ahlConn) Write(b []byte) (n int, err error) {
// 	defer func() { log.Println(c.ServerURL, " client write ", n, err) }()
// 	req, err := http.NewRequest(http.MethodPost, c.ServerURL, bytes.NewBuffer(b))
// 	if err != nil {
// 		return 0, err
// 	}
// 	req.AddCookie(c.Cookie)
// 	resp, err := c.Client.Do(req)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode == http.StatusAccepted {
// 		return len(b), nil
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		return 0, errors.New(strconv.Itoa(resp.StatusCode))
// 	}
// 	respBytes, err := ioutil.ReadAll(resp.Body)
// 	num := resp.Header.Get("Num")
// 	go func() {
// 		select {
// 		case c.Data <- ahlData{data: respBytes, err: err}:
// 			log.Println(c.ServerURL, "client write get ", len(respBytes), err, num)
// 		case <-c.CloseChan:
// 		}
// 		if num != "0" {
// 			inum, err := strconv.ParseUint(num, 10, 64)
// 			if err != nil {
// 				return
// 			}
// 			c.ReadLeft = inum
// 		}
// 	}()
// 	return len(b), nil
// }

// func (c *ahlConn) readMore() {
// 	req, err := http.NewRequest(http.MethodPost, c.ServerURL, nil)
// 	if err != nil {
// 		return
// 	}
// 	req.AddCookie(c.Cookie)
// 	resp, err := c.Client.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode == http.StatusAccepted {
// 		return
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		return
// 	}
// 	respBytes, err := ioutil.ReadAll(resp.Body)
// 	num := resp.Header.Get("Num")
// 	go func() {
// 		select {
// 		case c.Data <- ahlData{data: respBytes, err: err}:
// 			log.Println(c.ServerURL, "client readmore get ", len(respBytes), err, num)
// 		case <-c.CloseChan:
// 		}
// 		if num != "0" {
// 			inum, err := strconv.ParseUint(num, 10, 64)
// 			if err != nil {
// 				return
// 			}
// 			c.ReadLeft = inum
// 		}
// 	}()
// 	return
// }

// func (c *ahlConn) Close() error {
// 	defer func() {
// 		recover()
// 	}()
// 	close(c.CloseChan)
// 	log.Println("client : " + c.ServerURL + " close")
// 	req, err := http.NewRequest(http.MethodGet, c.ServerURL+"new?action=close", nil)
// 	if err != nil {
// 		return err
// 	}
// 	req.AddCookie(c.Cookie)
// 	client := c.Client
// 	_, err = client.Do(req)
// 	return err
// }
// func (c *ahlConn) LocalAddr() net.Addr {
// 	return nil
// }
// func (c *ahlConn) RemoteAddr() net.Addr {
// 	return nil
// }
// func (c *ahlConn) SetDeadline(t time.Time) error {
// 	return nil
// }
// func (c *ahlConn) SetReadDeadline(t time.Time) error {
// 	return nil
// }
// func (c *ahlConn) SetWriteDeadline(t time.Time) error {
// 	return nil
// }
