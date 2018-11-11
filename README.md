# Bi-endianess Bit Vector

bitset is Bit Vector (Array) library supporting both Little Endian and Big Endian for Golang.

bitset write bit vector to byte array with specified endianness (Little Endian or Big Endian) regardless of host endianness.
It enables to transfer a file holding bit vector to different endianness machine without any conversion process.

bitset calculate each bit in `uint64`.
bitset switches optimized bit vector operation (`Set`、`Unset`、`Get` etc...) by each host endianness.

[日本語の README](./README.ja.md) はこちら

## Features

- **Zero Copy** : Cast `[]byte` provided by user to `[]uint64` without any memory copy using `unsafe` package.
- **Bi-Endianness** : Switches bit vector operation by host endianness (Little Endian or Big Endian)
- **Compatibility** : Ready to transfer a file holding bit vector to different endianness machine without any conversion process

## Example

```go
func main() {
	// in memory usage
	buf := make([]byte, 2*8)
	b, _ := bitset.New(buf, bitset.LittleEndian)
	for _, v := range []uint{0, 1, 3, 6, 10, 64, 127, 128} {
		b.Set(v) // when v == 128 returns false because overflow
	}
	fmt.Println(buf) // [75 4 0 0 0 0 0 0 1 0 0 0 0 0 0 128]

	b.Unset(127)
	fmt.Println(b.Get(127)) // false

	for v, ok := b.FindFirstOne(0); ok; v, ok = b.FindFirstOne(v + 1) {
		fmt.Println(v) // 0 1 3 6 10 64
	}

	// File + mmap usage
	const pageSize = 4 * 1024
	f, _ := os.OpenFile("example.bit", os.O_RDWR|os.O_CREATE, 0666)
	f.Truncate(pageSize)
	buf, _ = syscall.Mmap(int(f.Fd()), 0, pageSize,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED)
	defer func() {
		f.Sync()
		syscall.Munmap(buf)
		f.Close()
	}()

	b, _ = bitset.New(buf, bitset.BigEndian)
	for v, ok := b.FindFirstOne(0); ok; v, ok = b.FindFirstOne(v + 1) {
		fmt.Println(v) // 0 1 3 6 10 64  if executed twice
	}
	for _, v := range []uint{0, 1, 3, 6, 10, 64, 127, 128} {
		b.Set(v) // when v == 128 returns false because overflow
	}
	fmt.Println(buf) // [0 0 0 0 0 0 4 75 128 0 0 0 0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 ....
}
```

## Installation

```bash
go get github.com/kawasin73/bitset
```

## Notices

- Length of the buffer (`[]byte`) provided by user MUST be a multiple of 8. (or `New()` returns error `bitset.ErrInvalidLength`)
- bitset supports only `Little Endian` and `Big Endian`, not `middle endian` or other endianness.
- bitset never auto expand provided buffer. If you need to expand bit vector then re-create `bitset.BitVec` with expanded buffer by user.

## LICENSE

MIT
