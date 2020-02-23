package clients

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/model"
)

// NewH2C NewH2C
func NewH2C(serverHost, serverPort, uri, auth string) *H2C {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: config.GetConf().NoVerify}
	return &H2C{
		ServerHost: serverHost,
		ServerPort: serverPort,
		URI:        uri,
		Auth:       &http.Cookie{Name: "pw", Value: auth}}
}

// H2C H2C
type H2C struct {
	ServerHost string
	ServerPort string
	URI        string
	Auth       *http.Cookie
}

// Dial Dial
func (c *H2C) Dial(addr model.ANetAddr) (net.Conn, error) {
	pr, pw := io.Pipe()
	req, err := http.NewRequest(http.MethodPut, "https://"+c.ServerHost+":"+c.ServerPort+"/"+c.URI+"/", ioutil.NopCloser(pr))
	if err != nil {
		return nil, err
	}
	addrBytes, err := json.Marshal(addr)
	if err != nil {
		return nil, err
	}
	req.AddCookie(c.Auth)
	req.AddCookie(&http.Cookie{Name: "addr", Value: url.QueryEscape(string(addrBytes))})
	// Send the request

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(resp.StatusCode))
	}
	return &h2cConn{w: pw, r: resp.Body}, nil
}

// Verify Verify
func (c *H2C) Verify(conn net.Conn) error {
	return nil
}

type h2cConn struct {
	w *io.PipeWriter
	r io.ReadCloser
}

func (c *h2cConn) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}

func (c *h2cConn) Write(b []byte) (n int, err error) {
	return c.w.Write(b)
}

func (c *h2cConn) Close() error {
	return c.r.Close()
}
func (c *h2cConn) LocalAddr() net.Addr {
	return nil
}
func (c *h2cConn) RemoteAddr() net.Addr {
	return nil
}
func (c *h2cConn) SetDeadline(t time.Time) error {
	return nil
}
func (c *h2cConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *h2cConn) SetWriteDeadline(t time.Time) error {
	return nil
}
