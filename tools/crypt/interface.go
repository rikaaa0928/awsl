package crypt

// Cryptor Cryptor
type Cryptor interface {
	Encrypt(data []byte)
	Decrypt(data []byte)
}
