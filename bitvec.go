package zcbit

import (
	"errors"
	"reflect"
	"unsafe"
)

const (
	wordBytes    = 8
	wordSize     = 64
	log2WordSize = 6
)

// errors
var (
	ErrInvalidByteOrder = errors.New("unsupported byte order")
)

// BitVec is bit vector component
type BitVec struct {
	vec []uint64
}

// New create *BitVec
func New(b []byte, endian ByteOrder) (*BitVec, error) {
	if endian != LittleEndian {
		return nil, ErrInvalidByteOrder
	}
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	header.Len /= wordBytes
	header.Cap /= wordBytes

	return &BitVec{
		vec: *(*[]uint64)(unsafe.Pointer(&header)),
	}, nil
}

// Test checks the bit is set.
func (b *BitVec) Test(i uint) bool {
	idx := i >> log2WordSize
	if int(idx) >= len(b.vec) {
		return false
	}
	return b.vec[idx]&(1<<(i&(wordSize-1))) != 0
}

// Set sets 1 to bit
func (b *BitVec) Set(i uint) bool {
	idx := i >> log2WordSize
	if int(idx) >= len(b.vec) {
		return false
	}
	b.vec[idx] |= 1 << (i & (wordSize - 1))
	return true
}

// Clear sets 0 to bit
func (b *BitVec) Clear(i uint) bool {
	idx := i >> log2WordSize
	if int(idx) >= len(b.vec) {
		return false
	}
	b.vec[idx] &^= 1 << (i & (wordSize - 1))
	return true
}
