package zcbit

import (
	"math/rand"
	"testing"
	"time"
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
	benchmarkBitVecSet(b, LittleEndian)
}

func BenchmarkBitVec_Set_BigEndian(b *testing.B) {
	benchmarkBitVecSet(b, BigEndian)
}

func benchmarkBitVecSet(b *testing.B, endian Endianness) {
	buf := make([]byte, randomSize/8)
	bv, err := New(buf, endian)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		bv.Set(randomSet[i%randomSize])
	}
}

func BenchmarkBitVec_Clear_LittleEndian(b *testing.B) {
	benchmarkBitVecClear(b, LittleEndian)
}

func BenchmarkBitVec_Clear_BigEndian(b *testing.B) {
	benchmarkBitVecClear(b, BigEndian)
}

func benchmarkBitVecClear(b *testing.B, endian Endianness) {
	buf := make([]byte, randomSize/8)
	bv, err := New(buf, endian)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		bv.Clear(randomSet[i%randomSize])
	}
}

func BenchmarkBitVec_Test_LittleEndian(b *testing.B) {
	benchmarkBitVecTest(b, LittleEndian)
}

func BenchmarkBitVec_Test_BigEndian(b *testing.B) {
	benchmarkBitVecTest(b, BigEndian)
}

func benchmarkBitVecTest(b *testing.B, endian Endianness) {
	buf := make([]byte, randomSize/8)
	bv, err := New(buf, endian)
	if err != nil {
		b.Fatal(err)
	}
	for _, v := range randomSet[:randomSize/2] {
		bv.Set(v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bv.Test(randomSet[i%randomSize])
	}
}

func BenchmarkSwap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		swapUint64(randomSet64[i%randomSize])
	}
}
