// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package hashgen

import (
	"container/list"
	"log"
	"sync"
)

const defaultCapacity = 1000

type UserRecord struct {

	// unique per hash POST request (job id)
	Uuid string

	// unique per user (when used)
	Salt []byte

	HashBytes []byte
}

// A LRU cache that is safe for concurrent access. All items are also added to
// persistent storage.
type LruCache struct {
	sync.RWMutex

	// Items are evicted when cache is full, can be fetched from persistent storage
	Capacity int

	linkedList *list.List
	catalog    map[string]*UserRecord

	// Persistent storage for crypto hashes by uuid
	dao *UserAccountFile
}

func NewCache(capacity int, filename string) *LruCache {

	useCapacity := defaultCapacity
	if capacity > 0 {
		useCapacity = capacity
	}
	return &LruCache{
		Capacity:   useCapacity,
		linkedList: list.New(),
		catalog:    make(map[string]*UserRecord, useCapacity), // size hint is default capacity
		dao:        New(filename),
	}
}

func (cache *LruCache) Add(uuid string, salt []byte, hashbytes []byte) {

	if cache == nil {
		return
	}
	if len(uuid) == 0 {
		return
	}
	log.Printf("Adding crypto hash bytes to cache for job %q\n", uuid)
	// Takes a write lock while mutating cache
	cache.Lock()
	defer cache.Unlock()
	cache.catalog[uuid] = &UserRecord{uuid, salt, hashbytes}

	// Persist to disk so that eviction does not cause data loss
	// TODO: check for errors
	cache.dao.Append(cache.catalog[uuid])
}

func (cache *LruCache) Get(uuid string) (value *UserRecord, ok bool) {

	if cache.catalog == nil {
		return
	}
	log.Printf("getting record with uuid: %s\n", uuid)
	cache.RLock()
	defer cache.RUnlock()

	if value, ok := cache.catalog[uuid]; ok {
		return value, true
	}
	return
}
