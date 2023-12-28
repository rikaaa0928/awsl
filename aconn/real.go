package aconn

import (
	"time"

	"github.com/rikaaa0928/awsl/global"
)

type Closer func() error

type IOer func([]byte) (int, error)

func CreateRealConn(c AConn) *RealConn {
	if c == nil {
		return nil
	}
	return &RealConn{
		AConn: c,
		close: c.Close,
		read:  c.Read,
		write: c.Write,
	}
}

type RealConn struct {
	AConn
	close Closer
	read  IOer
	write IOer
}

func (c *RealConn) Read(p []byte) (n int, err error) {
	if global.TimeOut > 0 {
		c.SetReadDeadline(time.Now().Add(time.Duration(global.TimeOut) * time.Second))
	}
	return c.read(p)
}

func (c *RealConn) Write(p []byte) (n int, err error) {
	return c.write(p)
}

func (c *RealConn) Close() error {
	return c.close()
}

func (c *RealConn) RegisterCloser(mid CLoserMid) {
	c.close = mid(c.close)
}

func (c *RealConn) RegisterReader(mid IOMid) {
	c.read = mid(c.read)
	//c.read = func(bytes []byte) (int, error) {
	//	return mid(c.read, bytes)(bytes)
	//}
}

func (c *RealConn) RegisterWriter(mid IOMid) {
	c.write = mid(c.write)
	//c.write = func(bytes []byte) (int, error) {
	//	return mid(c.write, bytes)(bytes)
	//}
}

type CLoserMid func(closer Closer) Closer

type IOMid func(io IOer) IOer
