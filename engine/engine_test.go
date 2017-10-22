package engine

import (
	"fmt"
	"github.com/huichen/murmur"
	"github.com/sillydong/goczd/godata"
	"github.com/sillydong/lbsengine/types"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

var e *Engine

func init() {
	rand.Seed(time.Now().Unix())
	e = &Engine{}
	e.Init(nil)
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lat, lng := RandomPoint()
		e.Add(&types.IndexedDocument{
			DocId:     uint64(i),
			Latitude:  lat,
			Longitude: lng,
			Fields: map[string]string{
				"a": "b",
			},
		})
	}
}

func TestSearch(t *testing.T) {
	lat, lng := RandomPoint()
	fmt.Println(lat, lng)
	result := e.Search(&types.SearchRequest{
		Latitude:  lat,
		Longitude: lng,
		Offset:    0,
		Limit:     100,
	})
	fmt.Printf("%+v\n", result)
	//x,_ := json.Marshal(result)
	//fmt.Println(string(x))
}

func BenchmarkSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lat, lng := RandomPoint()
		e.Search(&types.SearchRequest{
			Latitude:  lat,
			Longitude: lng,
			Offset:    0,
			Limit:     10,
			//SearchOption:&types.SearchOptions{
			//	Refresh:false,
			//	Circles:2,
			//},
		})
		//fmt.Println(resp.Count)
	}
}

func TestPoint(t *testing.T) {
	lat, lng := RandomPoint()
	fmt.Printf("%+v - %+v", lat, lng)
}

func RandomPoint() (lat, lng float64) {
	lat = 40.8137674 + godata.Round(rand.Float64()*0.01, 7)
	lng = -73.8525142 + godata.Round(rand.Float64()*0.01, 7)
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
