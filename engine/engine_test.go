package engine

import (
	"fmt"
	"github.com/huichen/murmur"
	"testing"
	"unsafe"
)

func TestParse(t *testing.T) {
	t.Logf("%+v", []byte(fmt.Sprintf("%d", 1234567890)))
	t.Logf("%+v", parse(1234567890))
}

func BenchmarkStringShard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := uint64(i)
		murmur.Murmur3([]byte(fmt.Sprintf("%d", x)))
	}
}

func BenchmarkUintShard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := uint64(i)
		murmur.Murmur3(parse(x))
	}
}

func parse(d uint64) []byte {
	s := fmt.Sprintf("%d", d)
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
