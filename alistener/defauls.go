package alistener

func DefaultAcceptMids(l AcceptMidor, ty, tag string) {
	switch ty {
	case "socks", "socks5", "socks4":
		l.RegisterAcceptor(NewSocksAcceptMid(tag))
	default:
	}
}
