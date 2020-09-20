package test

import (
	"encoding/json"
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

func TestJson(t *testing.T) {
	var a map[string]json.RawMessage
	err := json.Unmarshal([]byte(`{"a":"21"}`), &a)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a["a"])
	err = json.Unmarshal([]byte(`{"a":21}`), &a)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a["a"])
}

func TestJson2(t *testing.T) {
	var a []map[string]json.RawMessage
	err := json.Unmarshal([]byte(`[{"a":"21"}]`), &a)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a[0]["a"])

}

func TestMap(t *testing.T) {
	m := map[string]string{"a": "a"}
	var b map[string]interface{}
	var c interface{}
	var ok bool
	c = m
	b, ok = c.(map[string]interface{})
	if !ok {
		t.Fatal(ok)
	}
	t.Log(b["a"])
}

func TestJson3(t *testing.T) {
	var a []string
	err := json.Unmarshal([]byte(`["1","2","3"]`), &a)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a)

}
