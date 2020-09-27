package adialer

func DefaultDialMids(d ADialer, ty string) ADialer {
	switch ty {
	case "awsl", "tcp":
		d = NewAddrDataMid(d)
		d = NewAuthDataMid(d)
	default:
	}
	d = NewSendDataMid(d)
	return d
}
