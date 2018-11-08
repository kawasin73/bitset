package zcbit

import (
	"testing"
)

var (
	endians = []ByteOrder{LittleEndian, BigEndian}
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		size   int
		endian ByteOrder
		err    error
	}{
		"should create length : 0":    {size: 0, endian: LittleEndian, err: nil},
		"should create length : 1":    {size: 8, endian: LittleEndian, err: nil},
		"should create length : 1024": {size: 8 * 1024, endian: LittleEndian, err: nil},
		"should create length : 12":   {size: 100, endian: LittleEndian, err: nil},
		"should be big endian":        {size: 100, endian: BigEndian, err: nil},
		"unsupported endian unknown":  {size: 1024, endian: unknown, err: ErrInvalidByteOrder},
		"unsupported endian 100":      {size: 1024, endian: 100, err: ErrInvalidByteOrder},
	}

	for name, v := range tests {
		b, err := New(make([]byte, v.size), v.endian)
		if err != v.err {
			t.Errorf("%v b : %v, err : %v, expected err : %v", name, b, err, v.err)
		}
	}
}

func TestBitVec_Set_Clear(t *testing.T) {
	t.Run("set and clear in LittleEndian", func(t *testing.T) {
		buf := make([]byte, 8*10)
		b, err := New(buf, LittleEndian)
		if err != nil {
			t.Fatalf("failed to create bit vec %v", err)
		}
		arr := []uint{0, 1, 3, 6, 10, 64, 127}
		for _, v := range arr {
			if !b.Set(v) {
				t.Errorf("failed to set %v", v)
			}
		}
		buf2 := []byte{
			75, 4, 0, 0, 0, 0, 0, 0,
			1, 0, 0, 0, 0, 0, 0, 128,
			0, 0, 0, 0, 0, 0, 0, 0,
		}
		for i, v := range buf2 {
			if buf[i] != v {
				t.Errorf("set not match i : %v, v : %v, expected %v", i, buf[i], v)
			}
		}
		for _, v := range arr[1:] {
			if !b.Clear(v) {
				t.Errorf("failed to clear %v", v)
			}
		}
		buf3 := []byte{
			1, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		}
		for i, v := range buf3 {
			if buf[i] != v {
				t.Errorf("clear not match i : %v, v : %v, expected %v", i, buf[i], v)
			}
		}
	})

	t.Run("set and clear in BigEndian", func(t *testing.T) {
		buf := make([]byte, 8*10)
		b, err := New(buf, BigEndian)
		if err != nil {
			t.Fatalf("failed to create bit vec %v", err)
		}
		arr := []uint{0, 1, 3, 6, 10, 64, 127}
		for _, v := range arr {
			if !b.Set(v) {
				t.Errorf("failed to set %v", v)
			}
		}
		buf2 := []byte{
			0, 0, 0, 0, 0, 0, 4, 75,
			128, 0, 0, 0, 0, 0, 0, 1,
			0, 0, 0, 0, 0, 0, 0, 0,
		}
		for i, v := range buf2 {
			if buf[i] != v {
				t.Errorf("set not match i : %v, v : %v, expected %v", i, buf[i], v)
			}
		}
		for _, v := range arr[1:] {
			if !b.Clear(v) {
				t.Errorf("failed to clear %v", v)
			}
		}
		buf3 := []byte{
			0, 0, 0, 0, 0, 0, 0, 1,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		}
		for i, v := range buf3 {
			if buf[i] != v {
				t.Errorf("clear not match i : %v, v : %v, expected %v", i, buf[i], v)
			}
		}
	})

	for _, endian := range endians {
		t.Run("vector edge "+endian.String(), func(t *testing.T) {
			buf := make([]byte, 16)
			b, err := New(buf, endian)
			if err != nil {
				t.Fatalf("failed to create bit vec %v", err)
			}
			if !b.Set(16*8 - 1) {
				t.Errorf("failed to set last bit")
			}
			if b.Set(16 * 8) {
				t.Errorf("invalid to success to set next to the last bit")
			}
			if !b.Clear(16*8 - 1) {
				t.Errorf("failed to clear last bit")
			}
			if b.Clear(16 * 8) {
				t.Errorf("invalid to success to clear next to the last bit")
			}
		})
	}

}

func TestBitVec_Test(t *testing.T) {
	for _, endian := range endians {
		t.Run(endian.String(), func(t *testing.T) {
			buf := make([]byte, 8*3)
			b, err := New(buf, LittleEndian)
			if err != nil {
				t.Fatalf("failed to create bit vec %v", err)
			}
			arr := []uint{0, 1, 3, 6, 10, 64, 127}
			for _, v := range arr {
				if !b.Set(v) {
					t.Errorf("failed to set %v", v)
				}
			}
			for _, v := range arr {
				if !b.Test(v) {
					t.Errorf("failed to test %v not found", v)
				}
			}
			for _, v := range []uint{2, 4, 5, 7, 8, 100, 8*3*8 - 1, 8 * 3 * 8} {
				if b.Test(v) {
					t.Errorf("unexpectedly successed to test %v", v)
				}
			}
		})
	}
}
