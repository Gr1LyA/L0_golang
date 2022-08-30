package database

import "sync"

//type tableDB struct {
//	uidOrder string
//	jsonData string
//}

type mapMutex struct {
	mx sync.Mutex
	m  map[string]string
}

func NewCounters() *mapMutex {
	return &mapMutex{
		m: make(map[string]string),
	}
}

func (c *mapMutex) Load(key string) (string, bool) {
	c.mx.Lock()
	val, ok := c.m[key]
	c.mx.Unlock()
	return val, ok
}

func (c *mapMutex) Store(key string, value string) {
	c.mx.Lock()
	c.m[key] = value
	c.mx.Unlock()
}
