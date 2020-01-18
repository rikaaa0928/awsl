package servers

// AWSL AWSL
type AWSL struct {
	IP   string
	Port string
}

// Listen server
func (s AWSL) Listen() net.Listener {
	return nil
}

// ReadRemote server
func (s AWSL) ReadRemote(c net.Conn) (ANetAddr, error) {
	return nil, nil
}
