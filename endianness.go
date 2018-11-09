package bitset

import (
	"fmt"
	"unsafe"
)

// Endianness is LittleEndian or BigEndian
type Endianness uint8

func (e Endianness) String() string {
	switch e {
	case LittleEndian:
		return "little endian"
	case BigEndian:
		return "big endian"
	default:
		return fmt.Sprintf("unknown endian(%d)", e)
	}
}

// Use LittleEndian or BigEndian for BitSet initialization
const (
	unknown Endianness = iota
	LittleEndian
	BigEndian
)

var hostEndian Endianness

func init() {
	var x uint64 = 0x0011223344556677
	b := *(*[8]byte)(unsafe.Pointer(&x))
	if b[0] == 0x77 && b[7] == 0x00 {
		hostEndian = LittleEndian
	} else if b[0] == 0x00 && b[7] == 0x77 {
		hostEndian = BigEndian
	} else {
		// zcbit does not support Middle Endian etc...
		hostEndian = unknown
	}
}
