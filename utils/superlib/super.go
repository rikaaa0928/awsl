package superlib

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/rikaaa0928/awsl/utils/bitmap"
)

type constString string

const CTXSuperID constString = "superID"

var num uint32
var l sync.Mutex
var bm *bitmap.BitMap

func init() {
	bm = bitmap.NewBitMap(0xffffffff)
}

type SuperMSG struct {
	T   string
	MSG string
	ID  uint32
}

type UDPMSG struct {
	DstStr string
	SrcStr string
	Data   []byte
}

func NewUDPMSG(bb []byte, srcAddr net.Addr) (m UDPMSG, err error) {
	n := len(bb)
	minl := 4
	if n < minl {
		return m, fmt.Errorf("wrong udp data: %v", bb)
	}
	if bb[2] != 0 {
		return m, fmt.Errorf("Ignore frag: %v", bb[2])
	}
	var addr []byte
	if bb[3] == 1 {
		minl += 4
		if n < minl {
			return m, fmt.Errorf("wrong udp data: %v", bb)
		}
		addr = bb[minl-4 : minl]
	} else if bb[3] == 4 {
		minl += 16
		if n < minl {
			return m, fmt.Errorf("wrong udp data: %v", bb)
		}
		addr = bb[minl-16 : minl]
	} else if bb[3] == 3 {
		minl += 1
		if n < minl {
			return m, fmt.Errorf("wrong udp data: %v", bb)
		}
		l := bb[4]
		if l == 0 {
			return m, fmt.Errorf("wrong udp data: %v", bb)
		}
		minl += int(l)
		if n < minl {
			return m, fmt.Errorf("wrong udp data: %v", bb)
		}
		addr = bb[minl-int(l) : minl]
		addr = append([]byte{l}, addr...)
	} else {
		return m, fmt.Errorf("wrong udp data: %v", bb)
	}
	minl += 2
	if n <= minl {
		return m, fmt.Errorf("wrong udp data: %v", bb)
	}
	port := bb[minl-2 : minl]
	m.Data = bb[minl:]
	var s string
	if bb[3] == 3 {
		s = bytes.NewBuffer(addr[1:]).String()
	} else {
		s = net.IP(addr).String()
	}
	p := strconv.Itoa(int(binary.BigEndian.Uint16(port)))
	m.DstStr = net.JoinHostPort(s, p)
	m.SrcStr = srcAddr.String()
	return
}

func GetID(ctx context.Context) string {
	v := ctx.Value(CTXSuperID)
	if v == nil {
		return "0"
	}
	return strconv.FormatUint(uint64(v.(uint32)), 10)
}

func SetID(ctx context.Context, v uint32) context.Context {
	return context.WithValue(ctx, CTXSuperID, v)
}

func NewID() uint32 {
	l.Lock()
	defer func() {
		num++
		l.Unlock()
	}()
	for !bm.Set(num) {
		num++
	}
	return num
}

func RestoreID(id uint32) {
	bm.Del(id)
}
