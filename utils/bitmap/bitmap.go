package bitmap

import (
	"sync"
)

type BitMap struct {
	m []uint64
	sync.Mutex
}

func NewBitMap(max uint32) *BitMap {
	l := max / 64
	return &BitMap{m: make([]uint64, l+1)}
}

func (b *BitMap) Set(i uint32) bool {
	b.Lock()
	defer b.Unlock()
	if i >= uint32(len(b.m))*64 {
		return true
	}
	index := i / 64
	offset := i % 64
	if b.get(index, offset) {
		return false
	}
	b.m[index] = b.m[index] | uint64(1)<<offset
	return true
}

func (b *BitMap) get(index, offset uint32) bool {
	return (b.m[index] & (uint64(1) << offset)) != 0
}

func (b *BitMap) Del(i uint32) {
	b.Lock()
	defer b.Unlock()
	if i >= uint32(len(b.m))*8 {
		return
	}
	index := i / 64
	offset := i % 64
	b.m[index] = b.m[index] & ((uint64(1) << offset) ^ 0xff)
}
