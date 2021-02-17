package adialer

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"strconv"
	"sync"

	"golang.org/x/net/websocket"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/global"
	"github.com/rikaaa0928/awsl/utils/ctxdatamap"
)

type sAwslConf struct {
	remoteHost string
	remotePort string
	auth       string
	uri        string
	skipVerify bool
	wsConfig   *websocket.Config
}

var AWSLConfs = map[string]*sAwslConf{}

var awslLock sync.Mutex

func NewAWSL(tag string, conf map[string]interface{}) ADialer {
	awslLock.Lock()
	awslConf, ok := AWSLConfs[tag]
	if !ok {
		awslConf = &sAwslConf{}
		if skip, ok := conf["skipVerify"].(bool); ok && skip {
			awslConf.skipVerify = true
		}
		var err error
		awslConf.remoteHost = conf["host"].(string)
		awslConf.remotePort = strconv.Itoa(int(conf["port"].(float64)))
		awslConf.auth = conf["auth"].(string)
		awslConf.uri = conf["uri"].(string)
		awslConf.wsConfig, err = websocket.NewConfig("wss://"+net.JoinHostPort(awslConf.remoteHost, awslConf.remotePort)+"/"+awslConf.uri,
			"https://"+net.JoinHostPort(awslConf.remoteHost, awslConf.remotePort)+"/")
		if err != nil {
			panic(err)
		}
		AWSLConfs[tag] = awslConf
	}
	awslLock.Unlock()
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		tcpConn, err := net.Dial("tcp", net.JoinHostPort(awslConf.remoteHost, awslConf.remotePort))
		if err != nil {
			return ctx, nil, err
		}
		tlsClient := tls.Client(tcpConn, &tls.Config{InsecureSkipVerify: awslConf.skipVerify, ServerName: awslConf.remoteHost})
		ws, err := websocket.NewClient(awslConf.wsConfig, tlsClient)
		if err != nil {
			log.Println("awsl client new client", err)
			return ctx, nil, err
		}
		conn := aconn.NewAConn(ws)
		conn.SetEndAddr(addr)
		//ctx = context.WithValue(ctx, global.CTXSendAuth, AWSLConf.auth)
		ctx = ctxdatamap.Set(ctx, global.TransferAuth, awslConf.auth)
		ctx = context.WithValue(ctx, global.CTXOutType, "awsl")
		return ctx, conn, nil
	}
}
