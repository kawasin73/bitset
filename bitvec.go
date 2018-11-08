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
	ErrUnsupportedArch  = errors.New("unsupported host byte order")
)

// BitVec is bit vector component
type BitVec struct {
	vec  []uint64
	swap bool
}

// New create *BitVec
func New(b []byte, endian ByteOrder) (*BitVec, error) {
	if endian != LittleEndian && endian != BigEndian {
		return nil, ErrInvalidByteOrder
	} else if hostEndian != LittleEndian && hostEndian != BigEndian {
		return nil, ErrUnsupportedArch
	}
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	header.Len /= wordBytes
	header.Cap /= wordBytes

	return &BitVec{
		vec:  *(*[]uint64)(unsafe.Pointer(&header)),
		swap: endian != hostEndian,
	}, nil
}

// Test checks the bit is set.
func (b *BitVec) Test(i uint) bool {
	idx := i >> log2WordSize
	if int(idx) >= len(b.vec) {
		return false
	}
	if b.swap {
		v := swapUint64(b.vec[idx])
		return v&(1<<(i&(wordSize-1))) != 0
	} else {
		return b.vec[idx]&(1<<(i&(wordSize-1))) != 0
	}
}

// Set sets 1 to bit
func (b *BitVec) Set(i uint) bool {
	idx := i >> log2WordSize
	if int(idx) >= len(b.vec) {
		return false
	}
	if b.swap {
		v := swapUint64(b.vec[idx])
		v |= 1 << (i & (wordSize - 1))
		b.vec[idx] = swapUint64(v)
	} else {
		b.vec[idx] |= 1 << (i & (wordSize - 1))
	}
	return true
}

// Clear sets 0 to bit
func (b *BitVec) Clear(i uint) bool {
	idx := i >> log2WordSize
	if int(idx) >= len(b.vec) {
		return false
	}
	if b.swap {
		v := swapUint64(b.vec[idx])
		v &^= 1 << (i & (wordSize - 1))
		b.vec[idx] = swapUint64(v)
	} else {
		b.vec[idx] &^= 1 << (i & (wordSize - 1))
	}
	return true
}

func swapUint64(n uint64) uint64 {
	return ((n & 0x00000000000000FF) << 56) |
		((n & 0x000000000000FF00) << 40) |
		((n & 0x0000000000FF0000) << 24) |
		((n & 0x00000000FF000000) << 8) |
		((n & 0x000000FF00000000) >> 8) |
		((n & 0x0000FF0000000000) >> 24) |
		((n & 0x00FF000000000000) >> 40) |
		((n & 0xFF00000000000000) >> 56)
}
