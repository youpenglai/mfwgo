package consul

import (
	"log"
	"sync"
	"time"
)

type ConsulCache struct {
	kv map[string]*ConsulCacheItem
	l  sync.Mutex
}

var g_consulCache = ConsulCache{
	kv: make(map[string]*ConsulCacheItem),
}

func (c *ConsulCache) Get(key string) *ServiceInfo {
	c.l.Lock()
	defer c.l.Unlock()

	v, ok := c.kv[key]
	if !ok {
		return nil
	}
	return v.Get()
}

func (c *ConsulCache) Set(key string, val []*ServiceInfo) error {
	c.l.Lock()
	defer c.l.Unlock()

	item, exist := c.kv[key]
	if !exist {
		item = newConsulCacheItem(key)
		c.kv[key] = item
	}

	return item.Set(val)
}

func (c *ConsulCache) IsExist(key string) bool {
	c.l.Lock()
	defer c.l.Unlock()

	_, ok := c.kv[key]
	return ok
}

func (c *ConsulCache) Delete(key string) error {
	c.l.Lock()
	defer c.l.Unlock()

	delete(c.kv, key)
	return nil
}

type ConsulCacheItem struct {
	n    int
	l    sync.Mutex
	data []*ServiceInfo
}

func newConsulCacheItem(serviceName string) *ConsulCacheItem {
	item := new(ConsulCacheItem)
	item.data = make([]*ServiceInfo, 0, 10)

	go func() {
		s := NewConsulService()
		for {
			time.Sleep(time.Minute * 30)
			if err := s.GetServicesToCache(serviceName); err != nil {
				log.Println(err.Error())
				return
			}
		}
	}()
	return item
}

func (c *ConsulCacheItem) Get() *ServiceInfo {
	c.l.Lock()
	defer c.l.Unlock()

	if 0 == len(c.data) {
		return nil
	}

	if c.n >= len(c.data) {
		c.n = 0
	}

	count := c.n
	c.n++
	return c.data[count]
}

func (c *ConsulCacheItem) Set(val []*ServiceInfo) error {
	c.l.Lock()
	defer c.l.Unlock()

	c.data = val
	return nil
}

/*
func (c *ConsulCacheItem) IsExist(ID string) bool {
	c.l.Lock()
	defer c.l.Unlock()

	for index, _ := range c.data {
		if c.data[index].ID == ID {
			return true
		}
	}

	return false
}

func (c *ConsulCacheItem) Delete(ID string) error {
	c.l.Lock()
	defer c.l.Unlock()

	for index, _ := range c.data {
		if c.data[index].ID == ID {
			length := len(c.data)
			if index < length-1 {
				c.data[index] = c.data[length-1]
			}
			c.data = c.data[:length-1]
			return nil
		}
	}

	return errors.New("Not Found Service ID: " + ID)
}
*/