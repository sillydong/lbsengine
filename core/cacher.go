package core

import (
	"github.com/patrickmn/go-cache"
	"time"
)

//缓存器
type Cacher struct{
	client *cache.Cache
}

//读缓存
func(c *Cacher)Get(key string) interface{}{
	if val,ok := c.client.Get(key);ok{
		return val
	}
	return nil
}

//写缓存
func(c *Cacher)Set(key string,value interface{}){
	c.client.Set(key,value,5*time.Minute)
}
