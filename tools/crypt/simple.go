package crypt

// Simple simple one
type Simple uint8

// Encrypt Encrypt
func (c Simple) Encrypt(data []byte) {
	for i, v := range data {
		data[i] = v + byte(c)
	}
}

// Decrypt Decrypt
func (c Simple) Decrypt(data []byte) {
	for i, v := range data {
		data[i] = v - byte(c)
	}
}
