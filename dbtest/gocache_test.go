package dbtest

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"testing"
	"time"
)

var gocacheclient *cache.Cache

func init() {
	gocacheclient = cache.New(5*time.Minute, 10*time.Minute)
}

func BenchmarkASet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%v", i)
		value := fmt.Sprintf("value%v", i)
		gocacheclient.Set(key, value, 5*time.Minute)
	}
}

func BenchmarkAGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%v", i)
		gocacheclient.Get(key)
	}
}
