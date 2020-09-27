package alistener

func DefaultAcceptMids(l AcceptMidor, ty, tag string, conf map[string]interface{}) {
	switch ty {
	case "socks", "socks5", "socks4":
		l.RegisterAcceptor(NewSocksAcceptMid(tag))
	case "awsl", "tcp", "h2c":
		l.RegisterAcceptor(NewAddrAuthMid(conf))
	default:
	}
}
