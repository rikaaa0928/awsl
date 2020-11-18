package adialer

func DefaultDialMids(d ADialer, ty string) ADialer {
	switch ty {
	case "awsl", "tcp":
		d = NewAddrDataMid(d)
		d = NewSendDataMid(d)
	default:
	}
	return d
}
