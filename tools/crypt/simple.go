package crypt

// Simple simple one
type Simple uint8

// Encrypt Encrypt
func (c Simple) Encrypt(data []byte, n int) {
	for i, v := range data {
		if i >= n {
			return
		}
		data[i] = v + byte(c)
	}
}

// Decrypt Decrypt
func (c Simple) Decrypt(data []byte, n int) {
	for i, v := range data {
		if i >= n {
			return
		}
		data[i] = v - byte(c)
	}
}
