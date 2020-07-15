package tools

import (
	"io"
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
	defer src.Close()
	buf := MemPool.Get(65536)
	defer MemPool.Put(buf)

	io.CopyBuffer(dst, src, buf)
	/*for {
		// SetReadTimeout(src, 3*time.Second)
		n, err := src.Read(buf)
		if n > 0 {
			if _, wErr := dst.Write(buf[0:n]); wErr != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}*/
}
