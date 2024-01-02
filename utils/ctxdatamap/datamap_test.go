package ctxdatamap_test

import (
	"context"
	"github.com/rikaaa0928/awsl/utils/ctxdatamap"
	"github.com/rikaaa0928/awsl/utils/safer"
	"testing"
)

func TestDataMap(t *testing.T) {
	ctx := context.Background()
	str := `{"x":"z"}`
	ctxdatamap.Set(ctx, "e", "f")
	ctxdatamap.Parse(ctx, []byte(str))
	bs := ctxdatamap.Bytes(ctx)
	t.Log(string(bs))
	ctx = ctxdatamap.Parse(ctx, []byte(str))
	bs = ctxdatamap.Bytes(ctx)
	t.Log(string(bs))

	ctx = ctxdatamap.Set(ctx, "a", "b")
	ctxdatamap.Set(ctx, "c", "d")
	bs = ctxdatamap.Bytes(ctx)
	t.Log(string(bs))
	str = `{"xx":"z"}`
	ctxdatamap.Parse(ctx, []byte(str))
	bs = ctxdatamap.Bytes(ctx)
	t.Log(string(bs))
}

func TestDataMapMagic(t *testing.T) {
	ctx := context.Background()
	ctx = ctxdatamap.Set(ctx, "a", "b")
	bs := ctxdatamap.Bytes(ctx)
	t.Log(string(bs))
	safer.Handle(bs, safer.Magic(byte(len(bs))), false)
	t.Log(string(bs))
	ctx2 := context.Background()
	safer.Handle(bs, safer.Magic(byte(len(bs))), true)
	ctx2 = ctxdatamap.Parse(ctx2, bs)
	bs = ctxdatamap.Bytes(ctx2)
	t.Log(string(bs))
}
