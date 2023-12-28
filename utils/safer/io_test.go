package safer_test

import (
	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/utils/safer"
	"testing"
)

func TestIO(t *testing.T) {
	tm := 256
	magic := uint32(tm)
	if magic == 0 {
		magic = 1
	}
	r := safer.IOSaferFactory(&magic, true)
	w := safer.IOSaferFactory(&magic, false)
	var mData []byte
	var wf aconn.IOer = func(bytes []byte) (int, error) {
		t.Logf("%v", bytes)
		mData = bytes
		return len(bytes), nil
	}

	wwf := w(wf)
	rwf := r(wf)
	data := []byte{1, 1, 2, 0, 255}
	n, _ := wwf(data)
	t.Log(n)
	n, _ = rwf(mData)
	t.Log(n)
}
