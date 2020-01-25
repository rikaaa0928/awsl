package model

// ANetAddr addr
type ANetAddr struct {
	Typ  int //4 6 1
	Host string
	Port int
}

// AddrWithAuth addr
type AddrWithAuth struct {
	ANetAddr
	Auth string
}
