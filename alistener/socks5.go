package alistener

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/utils"
)

// ErrUDP ErrUDP
var ErrUDP = errors.New("udp error")

func NewSocksAcceptMid(inTag string) AcceptMid {
	return func(next Accepter) Accepter {
		return func(ctx context.Context) (context.Context, aconn.AConn, error) {
			ctx, conn, err := next(ctx)
			ctx = context.WithValue(ctx, CTXIntag, inTag)
			if err != nil {
				return ctx, nil, err
			}
			buf := utils.GetMem(65536)
			defer utils.PutMem(buf)
			n, err := conn.Read(buf)
			if err != nil {
				conn.Close()
				return ctx, nil, err
			}
			if n < 1 {
				conn.Close()
				return ctx, nil, errors.New("invalid length")
			}
			switch buf[0] {
			case 5:
				_, err = conn.Write([]byte("\x05\x00"))
				if err != nil {
					conn.Close()
					return ctx, nil, err
				}
				return socks5Stage2(ctx, conn, buf)
			case 4:
				return socks4(ctx, conn, buf, n)
			default:
				conn.Close()
				return ctx, nil, errors.New("unsuported type")
			}
		}
	}
}

func socks4(ctx context.Context, conn aconn.AConn, buf []byte, n int) (context.Context, aconn.AConn, error) {
	if n < 8 {
		conn.Close()
		return ctx, nil, errors.New("invalid length" + strconv.Itoa(n))
	}
	portBytes := buf[2:4]
	ipBytes := buf[4:8]
	host := strconv.Itoa(int(ipBytes[0])) + "." + strconv.Itoa(int(ipBytes[1])) + "." + strconv.Itoa(int(ipBytes[2])) + "." + strconv.Itoa(int(ipBytes[3]))
	var port int16
	tb := bytes.NewBuffer(portBytes)
	err := binary.Read(tb, binary.BigEndian, &port)
	if err != nil {
		conn.Close()
		return ctx, nil, err
	}
	buf[0] = 0
	buf[1] = 90
	_, err = conn.Write(buf[:8])
	if err != nil {
		conn.Close()
		return ctx, nil, err
	}
	conn.SetEndAddr(host, int(port), "tcp")
	return ctx, conn, nil
}

func socks5Stage2(ctx context.Context, conn aconn.AConn, buf []byte) (context.Context, aconn.AConn, error) {

	n, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		return ctx, nil, err
	}
	if n < 2 {
		conn.Close()
		return ctx, nil, errors.New("invalid length")
	}
	host, _ := getRemoteHost5(buf[:n])
	remotePort := getRemotePort5(buf[:n])
	port := remotePort
	switch buf[1] {
	case 1:
		_, err = conn.Write([]byte("\x05\x00\x00\x01\x00\x00\x00\x00\xff\xff"))
		if err != nil {
			conn.Close()
			return ctx, nil, err
		}
		conn.SetEndAddr(host, int(port), "tcp")
		return ctx, conn, nil
	case 3:
		conn.Close()
		return ctx, nil, ErrUDP
	default:
		conn.Close()
		return ctx, nil, errors.New("unsuported or invalid cmd : " + strconv.Itoa(int(buf[1])))
	}
}

func getRemoteHost5(data []byte) (s string, t int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%+v\n", err)
		}
	}()
	if data[3] == byte(0x03) {
		s = string(data[5 : len(data)-2])
		t = 0
		return
	}
	if data[3] == byte(0x01) {
		t = 4
		for i := 0; i < len(data)-6; i++ {
			s += strconv.Itoa(int(data[4+i]))
			if i != len(data)-7 {
				s += "."
			}
		}
		return
	}
	t = 6
	s += "["
	for i := 0; i < len(data)-7; i += 2 {
		s += strconv.FormatInt(int64(data[4+i]), 16)
		s += fmt.Sprintf("%02x", int(data[5+i]))
		if i != len(data)-8 {
			s += ":"
		}
	}
	s += "]"
	return
}

func getRemotePort5(data []byte) (x int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%+v\n", err)
		}
	}()
	tt := data[len(data)-2:]
	tb := bytes.NewBuffer(tt)
	var y uint16
	binary.Read(tb, binary.BigEndian, &y)
	x = int(y)
	return
}
