package zcbit

import (
	"fmt"
	"unsafe"
)

// ByteOrder is LittleEndian or BigEndian
type ByteOrder uint8

func (bo ByteOrder) String() string {
	switch bo {
	case LittleEndian:
		return "little endian"
	case BigEndian:
		return "big endian"
	default:
		return fmt.Sprintf("unknown endian(%d)", bo)
	}
}

// Use LittleEndian or BigEndian for BitVec initialization
const (
	unknown ByteOrder = iota
	LittleEndian
	BigEndian
)

var hostEndian ByteOrder

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
