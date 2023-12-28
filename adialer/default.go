package adialer

func DefaultDialMids(d ADialer, ty string) ADialer {
	switch ty {
	case "awsl", "tcp", "quic":
		d = NewAddrDataMid(d)
		d = NewSendDataMid(d, ty)
	default:
	}
	return d
}
