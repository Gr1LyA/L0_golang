package model

import "sync"

type MapMutex struct {
	mx sync.RWMutex
	m  map[string]string
}

func NewRWMap() *MapMutex {
	return &MapMutex{
		m: make(map[string]string),
	}
}

func (c *MapMutex) Load(key string) (string, bool) {
	c.mx.RLock()
	val, ok := c.m[key]
	c.mx.RUnlock()
	return val, ok
}

func (c *MapMutex) Store(key string, value string) {
	c.mx.Lock()
	c.m[key] = value
	c.mx.Unlock()
}