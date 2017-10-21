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

//读缓存
func (c *Cacher) Get(key string, offset, limit int) (types.ScoredDocuments, int) {
	if val, ok := c.client.Get(key); ok {
		docs := val.(types.ScoredDocuments)
		return docs[offset : offset+limit], len(docs)
	}
	return nil, 0
}

//写缓存
func (c *Cacher) Set(key string, value types.ScoredDocuments) {
	c.client.Set(key, value, 5*time.Minute)
}
