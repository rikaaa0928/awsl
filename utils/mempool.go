package utils

import "sync"

var poolMap map[int]*sync.Pool

func init() {
	poolMap = make(map[int]*sync.Pool)
}

func GetMem(size int) []byte {
	pool, ok := poolMap[size]
	if !ok {
		poolMap[size] = &sync.Pool{New: func() interface{} {
			return make([]byte, size)
		}}
		pool = poolMap[size]
	}
	return pool.Get().([]byte)
}

func PutMem(b []byte) {
	size := cap(b)
	pool, ok := poolMap[size]
	if !ok {
		return
	}
	pool.Put(b)
}
