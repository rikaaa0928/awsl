package model

const (
	IPV4ADDR = 4
	IPV6ADDR = 6
	RAWADDR  = 1
	TCP      = 0
	UDP      = 1
)

// ANetAddr addr
type ANetAddr struct {
	Typ  int //4 6 1
	Host string
	Port int
	CMD  int // 0 tcp 1 udp
}

// AddrWithAuth addr
type AddrWithAuth struct {
	ANetAddr
	Auth string
}
