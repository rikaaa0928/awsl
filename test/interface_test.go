package test

import (
	"fmt"
	"testing"
)

type AB interface {
	A(string) string
	B(int) string
}

type aStruct struct {
	a string
}

type bStruct struct {
	FA func(string) string
	FB func(int) string
}

func (a aStruct) A(s string) string {
	return a.a + ":" + s
}

func (a aStruct) B(s int) string {
	return fmt.Sprintf("%s:%d", a.a, s)
}

func (a bStruct) A(s string) string {
	return a.FA(s)
}

func (a bStruct) B(s int) string {
	return a.FB(s)
}

func TestAB(t *testing.T) {
	a := aStruct{"a"}
	b := aStruct{"b"}
	c := bStruct{a.A, b.B}
	t.Log(c.FA("c"))
	t.Log(c.FB(1))
}
