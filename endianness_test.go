package bibit

import "testing"

func TestHostOrder(t *testing.T) {
	t.Log("Host Endian :", hostEndian)
	if hostEndian != LittleEndian && hostEndian != BigEndian {
		t.Errorf("unknown endian %v", hostEndian)
	}
}
