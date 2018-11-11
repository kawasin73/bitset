# Bi-endianess Bit Vector

bitset は、ビッグエンディアン、リトルエンディアンの両方に対応したビットベクトルのGo言語で実装されたライブラリです。

マシンのエンディアンに関わらず初期化時に指定されたエンディアンでビットベクトルを構築することで、ファイルに保存したバイト列を異なるエンディアンのマシン上で利用することができます。

内部では uint64 で操作を行い、同じ振る舞いをするようにマシンのエンディアンを検出して最適化された処理に切り替えます。

このパッケージの個々の操作は、[github.com/willf/bitset](https://github.com/willf/bitset) を参考にして作られています。

## 特徴

- Zero Copy : `unsafe` パッケージを利用して渡された `[]byte` を `[]uint64` に型変換するためメモリコピーが発生しない。
- Bie-Endian : バイト列を指定されたエンディアンで操作するために、マシンのエンディアンを検出してそれぞれに最適化された処理に切り替える。
- Compatibility : 異なるエンディアンのマシンにファイルを転送して利用する時に変換をする必要がない。

## 使用例

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

## インストール方法

```bash
go get github.com/kawasin73/bitset
```

## 制限

- 初期化時に渡すバイト列の長さは、8 の倍数であることが要求されます。8 で割り切れない長さのバイト列で初期化された場合は `bitset.ErrInvalidLength` エラーを発生させます。
- マシンのエンディアン・バイト列を操作するエンディアンは、それぞれビッグエンディアンとリトルエンディアンのみに対応しています。ミドルエンディアンなど他のエンディアンには対応していません。
- サイズの自動拡張は行いません。バイト列に保存できるサイズを超えた場合は、新しいバイト列を確保して `bitset.BitVec` を作成しなおしてください。

## ライセンス

MIT
