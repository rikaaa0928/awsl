package tools

import (
	"io"
	"log"
	"net"
	"time"
)

// SetReadTimeout set
func SetReadTimeout(c net.Conn, readTimeout time.Duration) {
	if readTimeout != 0 {
		c.SetReadDeadline(time.Now().Add(readTimeout))
	}
}

// PipeThenClose pip
func PipeThenClose(src, dst net.Conn) {
	defer dst.Close()
	buf := MemPool.Get(65536)
	defer func() {
		MemPool.Put(buf)
	}()
	SetReadTimeout(src, 60*time.Second)
	for {
		n, err := src.Read(buf)
		// read may return EOF with n > 0
		// should always process n > 0 bytes before handling error
		if n > 0 {
			// Note: avoid overwrite err returned by Read.
			if _, err := dst.Write(buf[0:n]); err != nil {
				log.Println("pip write:", err)
				break
			}
		}
		if err != nil {
			e, ok := err.(*net.OpError)
			if ok && !e.Temporary() {
				break
			}
			et, ok := err.(net.Error)
			if ok && et.Timeout() {
				break
			}
			if err == io.EOF {
				break
			}
			log.Println("pip read: " + err.Error())
			break
		}
	}
}
