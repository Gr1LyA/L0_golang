package database

import "sync"

type mapMutex struct {
	mx sync.RWMutex
	m  map[string]string
}

func NewCounters() *mapMutex {
	return &mapMutex{
		m: make(map[string]string),
	}
}

func (c *mapMutex) Load(key string) (string, bool) {
	c.mx.RLock()
	val, ok := c.m[key]
	c.mx.RUnlock()
	return val, ok
}

func (c *mapMutex) Store(key string, value string) {
	c.mx.Lock()
	c.m[key] = value
	c.mx.Unlock()
}
