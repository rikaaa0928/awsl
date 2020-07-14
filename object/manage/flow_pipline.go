package manage

import (
	"net"

	"github.com/Evi1/awsl/tools"
)

// PipLine PipLine
var PipLine sPipLine

func init() {
	PipLine.manager = OFlowManager

}

type sPipLine struct {
	manager *SFlowManager
}

func (p *sPipLine) PipeThenClose(src, dst net.Conn, srcIn bool, id int, host string) {
	defer dst.Close()
	defer src.Close()
	buf := tools.MemPool.Get(65536)
	defer tools.MemPool.Put(buf)

	// io.CopyBuffer(dst, src, buf)
	for {
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
		if p.manager != nil && srcIn {
			p.manager.AddIn(id, host, int64(n))
		} else if p.manager != nil {
			p.manager.AddOut(id, host, int64(n))
		}
	}
}
