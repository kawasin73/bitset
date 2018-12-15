package bitset

import (
	"unsafe"
	"encoding/binary"
)

var hostEndian binary.ByteOrder

func init() {
	var x uint64 = 0x0011223344556677
	b := *(*[8]byte)(unsafe.Pointer(&x))
	if b[0] == 0x77 && b[7] == 0x00 {
		hostEndian = binary.LittleEndian
	} else if b[0] == 0x00 && b[7] == 0x77 {
		hostEndian = binary.BigEndian
	} else {
		// zcbit does not support Middle Endian etc...
		hostEndian = nil
	}
}

// HostEndian returns the endianness of machine.
func HostEndian() binary.ByteOrder {
	return hostEndian
}
