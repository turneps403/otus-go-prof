package hw04lrucache

import "sync"

// https://stackoverflow.com/questions/48861029/what-is-the-benefit-of-using-rwmutex-instead-of-mutex/48861083
// lock/unlock takes 100ns for every ops.
var (
	mapMutex = sync.RWMutex{}
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return Cache(&lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	})
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	wasInCache := false
	mapMutex.RLock()
	li, ok := lru.items[key]
	mapMutex.RUnlock()
	if ok {
		// remove if exists
		lru.queue.Remove(li)
		mapMutex.Lock()
		delete(lru.items, key)
		mapMutex.Unlock()
		wasInCache = true
	}
	// add to head
	el := cacheItem{key: key, value: value}
	li = lru.queue.PushFront(el)
	mapMutex.Lock()
	lru.items[key] = li
	mapMutex.Unlock()
	if lru.queue.Len() > lru.capacity {
		// remove the oldest
		li := lru.queue.Back()
		lru.queue.Remove(li)
		mapMutex.Lock()
		delete(lru.items, li.Val.(cacheItem).key)
		mapMutex.Unlock()
	}
	return wasInCache
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	mapMutex.RLock()
	li, ok := lru.items[key]
	mapMutex.RUnlock()
	if ok {
		ci := li.Val.(cacheItem)
		lru.Set(ci.key, ci.value)
		return ci.value, true
	}
	return nil, false
}

func (lru *lruCache) Clear() {
	lru.queue = NewList()
	mapMutex.Lock()
	lru.items = make(map[Key]*ListItem, lru.capacity)
	mapMutex.Unlock()
}
