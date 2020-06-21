package manage

import (
	"log"
	"strconv"

	"github.com/Evi1/awsl/tools"
)

func ncc(cns map[int]connNum, id int, tag string) {
	cn, ok := cns[id]
	if !ok {
		cns[id] = connNum{Tag: tag, Counter: tools.NewCounter()}
		cn = cns[id]
	}
	cn.Add(1)
}

// NewConnectionCount get new connection
func NewConnectionCount(isServer bool, id int, tag string) {
	if isServer {
		ncc(ServerConnectionNumber, id, tag)
		return
	}
	ncc(ClientConnectionNumber, id, tag)
}

func ccc(cns map[int]connNum, id int) {
	cn, ok := cns[id]
	if !ok {
		log.Panic("connection close count id error. id:" + strconv.Itoa(id))
	}
	cn.Add(-1)
}

// ConnectionCloseCount close
func ConnectionCloseCount(isServer bool, id int) {
	if isServer {
		ccc(ServerConnectionNumber, id)
		return
	}
	ccc(ClientConnectionNumber, id)
}

type connNum struct {
	Tag string
	*tools.Counter
}

// ServerConnectionNumber ServerConnectionNumber map
var ServerConnectionNumber map[int]connNum

// ClientConnectionNumber ClientConnectionNumber map
var ClientConnectionNumber map[int]connNum

// RemoteConnectionNumber RemoteConnectionNumber map
var RemoteConnectionNumber map[string]connNum

func init() {
	ServerConnectionNumber = make(map[int]connNum)
	ClientConnectionNumber = make(map[int]connNum)
}
