package adialer

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/utils/rwconn"
)

var client *http.Client

func init() {
	client = &http.Client{}
	client.Transport = &http.Transport{
		Proxy: nil,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func NewH2C(conf map[string]interface{}) ADialer {
	if skip, ok := conf["skipVerify"].(bool); ok && skip {
		client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return func(ctx context.Context, addr net.Addr) (aconn.AConn, error) {
		pr, pw := io.Pipe()
		req, err := http.NewRequest(http.MethodPut,
			"https://"+conf["host"].(string)+":"+strconv.Itoa(int(conf["port"].(float64)))+"/"+conf["uri"].(string),
			ioutil.NopCloser(pr))
		if err != nil {
			return nil, err
		}
		addrInfo, ok := addr.(aconn.AddrInfo)
		if !ok {
			log.Println("addr not aconn.AddrInfo")
			err = (&addrInfo).Parse(addr.Network(), addr.String())
			if err != nil {
				return nil, err
			}
		}
		addrBytes, err := json.Marshal(addrInfo)
		if err != nil {
			return nil, err
		}
		req.AddCookie(&http.Cookie{Name: "auth", Value: conf["auth"].(string)})
		req.AddCookie(&http.Cookie{Name: "addr", Value: url.QueryEscape(string(addrBytes))})
		// Send the request
		//resp, err := http.DefaultClient.Do(req)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			b, _ := ioutil.ReadAll(resp.Body)
			return nil, errors.New(strconv.Itoa(resp.StatusCode) + ". err body = " + string(b))
		}
		conn := aconn.NewAConn(rwconn.NewRWConn(pw, resp.Body))
		conn.SetEndAddr(addr)
		return conn, nil
	}
}
