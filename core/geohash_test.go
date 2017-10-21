package core

import (
	"github.com/mmcloughlin/geohash"
	"testing"
)

func TestEncode(t *testing.T) {
	latitude := 31.286305
	longitude := 121.448463
	hash := geohash.EncodeWithPrecision(latitude, longitude, 5)
	t.Log(hash)
}

func TestEncodeInt(t *testing.T) {
	latitude := 31.286305
	longitude := 121.448463
	hash := geohash.EncodeIntWithPrecision(latitude, longitude, 6)
	t.Log(hash)
}

func TestNeighbour(t *testing.T) {
	hashs := geohash.Neighbors("wtw3g")
	t.Logf("%+v", hashs)
}

func BenchmarkNeighbour(b *testing.B) {
	for i := 0; i < b.N; i++ {
		geohash.Neighbors("wtw3g")
	}
}

func TestLoop(t *testing.T) {
	latitude := 31.286305
	longitude := 121.448463
	neightbours := LoopNeighbours(latitude, longitude, 5, 3)
	t.Logf("%+v", len(neightbours))
}

func BenchmarkLoop(b *testing.B) {
	latitude := 31.286305
	longitude := 121.448463
	for i := 0; i < b.N; i++ {
		LoopNeighbours(latitude, longitude, 5, 1)
	}
}
