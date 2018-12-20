package bitset

import (
	"errors"
	"math/bits"
	"reflect"
	"unsafe"
	"encoding/binary"
)

const (
	wordBytes           = 8
	wordBits            = 64
	log2WordSize        = 6
	mask00111000        = 0x0000000000000038
	mask00000111        = 0x0000000000000007
	allBits      uint64 = 0xffffffffffffffff
)

// errors
var (
	ErrInvalidEndianness = errors.New("unsupported endianness")
	ErrInvalidLength     = errors.New("len(buffer) for zcbit must be N * 8")
	ErrUnsupportedArch   = errors.New("unsupported host endianness")
)

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

// BitSet is bit vector component
type BitSet struct {
	vec    []uint64
	orig   []byte
	swap   bool
	extend bool
}

// New create *BitSet
func New(b []byte, order binary.ByteOrder, extend bool) (*BitSet, error) {
	if len(b)%8 != 0 {
		return nil, ErrInvalidLength
	} else if order != binary.LittleEndian && order != binary.BigEndian {
		return nil, ErrInvalidEndianness
	} else if hostEndian != binary.LittleEndian && hostEndian != binary.BigEndian {
		return nil, ErrUnsupportedArch
	}
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	header.Len /= wordBytes
	header.Cap /= wordBytes

	return &BitSet{
		vec:    *(*[]uint64)(unsafe.Pointer(&header)),
		orig:   b, // refrain GC
		swap:   order != hostEndian,
		extend: extend,
	}, nil
}

// Get checks the bit is set.
func (b *BitSet) Get(i uint) bool {
	idx := i >> log2WordSize
	if int(idx) >= len(b.vec) {
		return false
	}
	if b.swap {
		return b.vec[idx]&(1<<(wordBits-(i&mask00111000)-8)<<(i&mask00000111)) != 0
	} else {
		return b.vec[idx]&(1<<(i&(wordBits-1))) != 0
	}
}

func (c *BitSet) extendVec(size int) bool {
	if len(c.vec) >= size {
		// do nothing
	} else if cap(c.vec) >= size {
		c.vec = c.vec[:size]
	} else {
		nextcap := size + size
		if size > 1024 {
			nextcap = size + size/4
		}
		if nextcap <= 0 {
			// overflow
			return false
		}
		newvec := make([]uint64, size, nextcap)
		copy(newvec, c.vec)
		c.vec = newvec
		// remove original reference to GC
		c.orig = nil
	}
	return true
}

// Set sets 1 to bit
func (b *BitSet) Set(i uint) bool {
	idx := i >> log2WordSize
	if int(idx) >= len(b.vec) {
		if !b.extend || !b.extendVec(int(idx +1)) {
			return false
		}
	}
	if b.swap {
		b.vec[idx] |= 1 << (wordBits - (i & mask00111000) - 8) << (i & mask00000111)
	} else {
		b.vec[idx] |= 1 << (i & (wordBits - 1))
	}
	return true
}

// Unset sets 0 to bit
func (b *BitSet) Unset(i uint) bool {
	idx := i >> log2WordSize
	if int(idx) >= len(b.vec) {
		return false
	}
	if b.swap {
		b.vec[idx] &^= 1 << (wordBits - (i & mask00111000) - 8) << (i & mask00000111)
	} else {
		b.vec[idx] &^= 1 << (i & (wordBits - 1))
	}
	return true
}

// FindFirstOne returns first 1 bit index and true.
// if not found then returns false
func (b *BitSet) FindFirstOne(i uint) (uint, bool) {
	idx := int(i >> log2WordSize)
	if idx >= len(b.vec) {
		return 0, false
	}
	if b.swap {
		v := swapUint64(b.vec[idx])
		v = v >> (i & (wordBits - 1))
		if v != 0 {
			return i + uint(bits.TrailingZeros64(v)), true
		}
		for idx++; idx < len(b.vec); idx++ {
			if b.vec[idx] != 0 {
				return uint(idx)*wordBits + uint(bits.TrailingZeros64(swapUint64(b.vec[idx]))), true
			}
		}
	} else {
		v := b.vec[idx] >> (i & (wordBits - 1))
		if v != 0 {
			return i + uint(bits.TrailingZeros64(v)), true
		}
		for idx++; idx < len(b.vec); idx++ {
			if b.vec[idx] != 0 {
				return uint(idx)*wordBits + uint(bits.TrailingZeros64(b.vec[idx])), true
			}
		}
	}
	return 0, false
}

// FindFirstZero returns first 0 bit index and true.
// if not found then returns false
// TODO: set tail
func (b *BitSet) FindFirstZero(i uint) (uint, bool) {
	idx := int(i >> log2WordSize)
	if idx >= len(b.vec) {
		return 0, false
	}
	if b.swap {
		offset := (i & (wordBits - 1))
		v := swapUint64(b.vec[idx]) >> offset
		trail := uint(bits.TrailingZeros64(^v))
		if trail < wordBits-offset {
			return i + trail, true
		}
		for idx++; idx < len(b.vec); idx++ {
			if b.vec[idx] != allBits {
				return uint(idx)*wordBits + uint(bits.TrailingZeros64(^swapUint64(b.vec[idx]))), true
			}
		}
	} else {
		offset := (i & (wordBits - 1))
		v := b.vec[idx] >> offset
		trail := uint(bits.TrailingZeros64(^v))
		if trail < wordBits-offset {
			return i + trail, true
		}
		for idx++; idx < len(b.vec); idx++ {
			if b.vec[idx] != allBits {
				return uint(idx)*wordBits + uint(bits.TrailingZeros64(^b.vec[idx])), true
			}
		}
	}
	return 0, false
}

// FindLastOne returns last 1 bit index and true.
// if not found then returns false
func (b *BitSet) FindLastOne() (uint, bool) {
	if b.swap {
		for i := len(b.vec); i > 0; i-- {
			v := b.vec[i-1]
			if v != 0 {
				v = swapUint64(v)
				return uint(i*wordBits - bits.LeadingZeros64(v) - 1), true
			}
		}
	} else {
		for i := len(b.vec); i > 0; i-- {
			v := b.vec[i-1]
			if v != 0 {
				return uint(i*wordBits - bits.LeadingZeros64(v) - 1), true
			}
		}
	}
	return 0, false
}

// Count returns the number of set bit
func (b *BitSet) Count() uint {
	var cnt uint
	for _, v := range b.vec {
		cnt += uint(bits.OnesCount64(v))
	}
	return cnt
}
