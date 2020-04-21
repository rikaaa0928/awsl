package crypt

// Cryptor Cryptor
type Cryptor interface {
	Encrypt(data []byte, n int)
	Decrypt(data []byte, n int)
}
