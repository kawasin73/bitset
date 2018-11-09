# Bi-endian BitSet

bitset is Bit Vector (Array) library supporting both Big Endian and Little Endian for Golang.

[日本語の README](./README.ja.md) はこちら

## Features

TODO

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

TODO:

## LICENSE

MIT
