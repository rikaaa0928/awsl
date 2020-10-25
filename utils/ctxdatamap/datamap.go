package ctxdatamap

import (
	"context"
	"encoding/json"
)

type constString string

const CTXMapData constString = "mapData"

func Set(ctx context.Context, key string, value interface{}) context.Context {
	d := ctx.Value(CTXMapData)
	if d == nil {
		d = make(map[string]interface{})
	}
	m := d.(map[string]interface{})
	m[key] = value
	ctx = context.WithValue(ctx, CTXMapData, m)
	return ctx
}

func Get(ctx context.Context, key string) interface{} {
	d := ctx.Value(CTXMapData)
	if d == nil {
		return nil
	}
	m := d.(map[string]interface{})
	return m[key]
}

func Parse(ctx context.Context, data []byte) context.Context {
	newMap := make(map[string]interface{})
	err := json.Unmarshal(data, &newMap)
	if err != nil {
		return ctx
	}
	d := ctx.Value(CTXMapData)
	if d == nil {
		d = make(map[string]interface{})
	}
	m := d.(map[string]interface{})
	for k, v := range newMap {
		m[k] = v
	}
	ctx = context.WithValue(ctx, CTXMapData, m)
	return ctx
}

func Bytes(ctx context.Context) []byte {
	d := ctx.Value(CTXMapData)
	if d == nil {
		return nil
	}
	b, err := json.Marshal(d)
	if err != nil {
		return nil
	}
	return b
}
