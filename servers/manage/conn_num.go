package manage

import (
	"log"
	"strconv"

	"github.com/Evi1/awsl/tools"
)

// NewConnectionCount get new connection
func NewConnectionCount(id int, tag string) {
	if ServerConnectionNumber == nil {
		ServerConnectionNumber = make(map[int]connNum)
	}
	cn, ok := ServerConnectionNumber[id]
	if !ok {
		ServerConnectionNumber[id] = connNum{Tag: tag, Counter: tools.NewCounter()}
		cn = ServerConnectionNumber[id]
	}
	cn.Add(1)
}

// ConnectionCloseCount close
func ConnectionCloseCount(id int) {
	if ServerConnectionNumber == nil {
		log.Panic("connection close count nil map")
	}
	cn, ok := ServerConnectionNumber[id]
	if !ok {
		log.Panic("connection close count id error" + strconv.Itoa(id))
	}
	cn.Add(-1)
}

type connNum struct {
	Tag string
	*tools.Counter
}

// ServerConnectionNumber ServerConnectionNumber map
var ServerConnectionNumber map[int]connNum
