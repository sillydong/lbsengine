package engine

import (
	"fmt"
	"github.com/huichen/murmur"
	"testing"
	"unsafe"
	"math/rand"
	"github.com/sillydong/lbsengine/types"
)

var e *Engine
func init() {
	e = &Engine{}
	e.Init(nil)
}

func TestSearch(t *testing.T) {
	result := e.Search(&types.SearchRequest{
		Latitude: 40.7137674,
		Longitude: -73.9525142,
		Offset:0,
		Limit:100,
	})
	fmt.Printf("%+v\n",result)
}

func BenchmarkSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.Search(&types.SearchRequest{
			Latitude:  40.7137674,
			Longitude: -73.9525142,
			Offset:    0,
			Limit:     100,
		})
	}
}


func randomPoints(n int) [][2]float64 {
	var points [][2]float64
	for i := 0; i < n; i++ {
		lat, lon := RandomPoint()
		points = append(points, [2]float64{lat, lon})
	}
	return points
}

func RandomPoint() (lat, lng float64) {
	lat = -90 + 180*rand.Float64()
	lng = -180 + 360*rand.Float64()
	return
}

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
