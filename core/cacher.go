package core

import (
	"github.com/patrickmn/go-cache"
	"github.com/sillydong/lbsengine/types"
	"time"
)

//缓存器
type Cacher struct {
	client *cache.Cache
}

func (c *Cacher) Init() {
	c.client = cache.New(5*time.Minute, 10*time.Minute)
}

//读缓存
func (c *Cacher) Get(key string, offset, limit int) (types.ScoredDocuments, int) {
	if val, ok := c.client.Get(key); ok {
		docs := val.(types.ScoredDocuments)
		size := len(docs)
		if offset > size {
			return nil, size
		} else if offset+limit > size {
			return docs[offset:], size
		} else {
			return docs[offset : offset+limit], size
		}
	}
	return nil, 0
}

//写缓存
func (c *Cacher) Set(key string, value types.ScoredDocuments) {
	c.client.Set(key, value, 5*time.Minute)
}
