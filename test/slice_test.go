package test

import "testing"

func TestSlice(t *testing.T) {
	s := make([]int, 10)
	s[0] = 1
	t.Logf("%p,%+v\n", s, s)
	t.Log(cap(s), len(s))
	c := func(l []int) []int {
		return l[:1]
	}(s)
	t.Logf("%p,%+v\n", c, c)
	t.Log(cap(c), len(c))
}
