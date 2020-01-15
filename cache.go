package screws

import (
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

//ICache 缓存器接口
type ICache interface {
	Set(key string, value []byte, flags uint32, expiration int32) error
	Get(key string) ([]byte, error)
	Expiration(key string) (int32, error)
	Delete(key string) error
}

//cache 缓存器
type cache struct {
	MC *memcache.Client
}

//NewCache 初始化缓存器
func NewCache(hosts ...string) ICache {
	mc := memcache.New(hosts...)
	if err := mc.Ping(); err != nil {
		log.Fatalln(err)
	}
	return &cache{
		MC: mc,
	}
}

//MSet 缓存设置
func (c *cache) Set(key string, value []byte, flags uint32, expiration int32) error {
	if err := c.MC.Set(&memcache.Item{Key: key, Value: value, Flags: flags, Expiration: expiration}); err != nil {
		return err
	}
	return nil
}

//GetValue 缓存查询
func (c *cache) Get(key string) ([]byte, error) {
	item, err := c.MC.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

//GetValue 有效期查询
func (c *cache) Expiration(key string) (int32, error) {
	item, err := c.MC.Get(key)
	if err != nil {
		return -1, err
	}
	return item.Expiration, nil
}

//MGet 缓存删除
func (c *cache) Delete(key string) error {
	return c.MC.Delete(key)
}
