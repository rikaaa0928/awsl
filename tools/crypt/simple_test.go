package crypt

import "testing"

func Test_Simple(t *testing.T) {
	c := Simple(1)
	data := []byte{0xff, byte('a')}
	t.Log(data)
	c.Encrypt(data, len(data)-1)
	t.Log(data)
	c.Decrypt(data, len(data)-1)
	t.Log(data)
}
