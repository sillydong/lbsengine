package dbtest

import (
	"github.com/go-redis/redis"
	"github.com/sillydong/goczd/godata"
	"math/rand"
	"strconv"
	"testing"
)

var rclient *redis.Client

func init() {
	rclient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   4,
	})
}

func BenchmarkRAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lat, lng := RandomPoint()
		rclient.GeoAdd("geo", &redis.GeoLocation{
			Name:      strconv.Itoa(i),
			Longitude: lng,
			Latitude:  lat,
		})
	}
}

func BenchmarkRSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lat, lng := RandomPoint()
		rclient.GeoRadius("geo", lng, lat, &redis.GeoRadiusQuery{
			Unit:     "m",
			WithDist: true,
			Sort:     "ASC",
			Count:    10,
		})
	}
}

func RandomPoint() (lat, lng float64) {
	lat = 40.8137674 + godata.Round(rand.Float64()*0.01, 7)
	lng = -73.8525142 + godata.Round(rand.Float64()*0.01, 7)
	return
}
