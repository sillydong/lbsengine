package core

import (
	"github.com/patrickmn/go-cache"
	"time"
	"github.com/sillydong/lbsengine/types"
)

//缓存器
type Cacher struct{
	client *cache.Cache
}

//读缓存
func(c *Cacher)Get(key string) types.IndexedDocument{
	if val,ok := c.client.Get(key);ok{
		return val.(types.IndexedDocument)
	}
	return types.IndexedDocument{}
}

//写缓存
func(c *Cacher)Set(key string,value types.IndexedDocument){
	c.client.Set(key,value,5*time.Minute)
}
