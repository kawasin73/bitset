package bitset

import (
	"testing"
	"encoding/binary"
)

func TestHostOrder(t *testing.T) {
	t.Log("Host Endian :", hostEndian)
	if hostEndian != binary.LittleEndian && hostEndian != binary.BigEndian {
		t.Errorf("unknown endian %v", hostEndian)
	}
}
