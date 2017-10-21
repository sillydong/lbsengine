package dbtest

import (
	"github.com/muesli/cache2go"
	"testing"
	"fmt"
	"time"
)

var cachegoclient *cache2go.CacheTable

func init() {
	cachegoclient = cache2go.Cache("mycache")
}

func BenchmarkBSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%v",i)
		value := fmt.Sprintf("value%v",i)
		cachegoclient.Add(key,5*time.Minute,value)
	}
}

func BenchmarkBGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%v",i)
		cachegoclient.Value(key)
	}
}
