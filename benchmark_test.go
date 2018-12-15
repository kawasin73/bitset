package bitset

import (
	"math/rand"
	"testing"
	"time"
	"encoding/binary"
)

const randomSize = 1024 * 8

var (
	randomSet   = make([]uint, randomSize)
	randomSet64 = make([]uint64, randomSize)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < randomSize; i++ {
		randomSet[i] = uint(rand.Intn(randomSize))
		randomSet64[i] = uint64(rand.Intn(randomSize))
	}
}

func BenchmarkBitVec_Set_LittleEndian(b *testing.B) {
	benchmarkBitVecSet(b, binary.LittleEndian)
}

func BenchmarkBitVec_Set_BigEndian(b *testing.B) {
	benchmarkBitVecSet(b, binary.BigEndian)
}

func benchmarkBitVecSet(b *testing.B, order binary.ByteOrder) {
	buf := make([]byte, randomSize/8)
	bv, err := New(buf, order)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		bv.Set(randomSet[i%randomSize])
	}
}

func BenchmarkBitVec_Unset_LittleEndian(b *testing.B) {
	benchmarkBitVecUnset(b, binary.LittleEndian)
}

func BenchmarkBitVec_Unset_BigEndian(b *testing.B) {
	benchmarkBitVecUnset(b, binary.BigEndian)
}

func benchmarkBitVecUnset(b *testing.B, order binary.ByteOrder) {
	buf := make([]byte, randomSize/8)
	bv, err := New(buf, order)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		bv.Unset(randomSet[i%randomSize])
	}
}

func BenchmarkBitVec_Get_LittleEndian(b *testing.B) {
	benchmarkBitVecGet(b, binary.LittleEndian)
}

func BenchmarkBitVec_Get_BigEndian(b *testing.B) {
	benchmarkBitVecGet(b, binary.BigEndian)
}

func benchmarkBitVecGet(b *testing.B, order binary.ByteOrder) {
	buf := make([]byte, randomSize/8)
	bv, err := New(buf, order)
	if err != nil {
		b.Fatal(err)
	}
	for _, v := range randomSet[:randomSize/2] {
		bv.Set(v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bv.Get(randomSet[i%randomSize])
	}
}

func BenchmarkSwap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		swapUint64(randomSet64[i%randomSize])
	}
}
