package adialer

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/global"
	"github.com/rikaaa0928/awsl/utils/ctxdatamap"
	"golang.org/x/net/websocket"
)

var AWSLConf = struct {
	sync.RWMutex
	Inited     bool
	remoteHost string
	remotePort string
	auth       string
	uri        string
	skipVerify bool
	wsConfig   *websocket.Config
}{}

func NewAWSL(conf map[string]interface{}) ADialer {
	AWSLConf.RLock()
	if !AWSLConf.Inited {
		AWSLConf.RUnlock()
		AWSLConf.Lock()
		if !AWSLConf.Inited {
			if skip, ok := conf["skipVerify"].(bool); ok && skip {
				AWSLConf.skipVerify = true
			}
			var err error
			AWSLConf.remoteHost = conf["host"].(string)
			AWSLConf.remotePort = strconv.Itoa(int(conf["port"].(float64)))
			AWSLConf.auth = conf["auth"].(string)
			AWSLConf.uri = conf["uri"].(string)
			AWSLConf.wsConfig, err = websocket.NewConfig("wss://"+net.JoinHostPort(AWSLConf.remoteHost, AWSLConf.remotePort)+"/"+AWSLConf.uri,
				"https://"+net.JoinHostPort(AWSLConf.remoteHost, AWSLConf.remotePort)+"/")
			if err != nil {
				panic(err)
			}
			AWSLConf.Inited = true
		}
		AWSLConf.Unlock()
	} else {
		AWSLConf.RUnlock()
	}
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		tcpConn, err := net.Dial("tcp", net.JoinHostPort(AWSLConf.remoteHost, AWSLConf.remotePort))
		if err != nil {
			return ctx, nil, err
		}
		tlsClient := tls.Client(tcpConn, &tls.Config{InsecureSkipVerify: AWSLConf.skipVerify, ServerName: AWSLConf.remoteHost})
		ws, err := websocket.NewClient(AWSLConf.wsConfig, tlsClient)
		if err != nil {
			log.Println("awsl client new client", err)
			return ctx, nil, err
		}
		conn := aconn.NewAConn(ws)
		conn.SetEndAddr(addr)
		//ctx = context.WithValue(ctx, global.CTXSendAuth, AWSLConf.auth)
		ctx = ctxdatamap.Set(ctx, global.TransferAuth, AWSLConf.auth)
		return ctx, conn, nil
	}
}
