package main

import (
	"github.com/spaolacci/murmur3"
	"sync"
)

type murmurHashGenerator struct {
	setMu       sync.RWMutex
	storeMu     sync.Mutex
	urlMapCache map[string]string
	urlSet      map[string]struct{}
}

func (mmh *murmurHashGenerator) Generate(input string) string {
	// load from cache
	mmh.setMu.RLock()
	shorten, ok := mmh.urlMapCache[input]
	mmh.setMu.RUnlock()
	if ok {
		return shorten
	}
	// generate short code
	sum32 := murmur3.Sum64([]byte(input))
	// convert to string
	//shortCode := strconv.FormatInt(int64(sum32), 36)
	shortCode := FormatInt62(sum32)
	// check if short code is already present ,using db unique index or bloom filter
	shortCode = handleCollisions(input, shortCode)

	return shortCode
}

func (mmh *murmurHashGenerator) Store(input, shorten string) {
	// 缓存短链，避免重复生成
	// todo：添加过期时间，LRU策略
	mmh.setMu.Lock()
	mmh.urlMapCache[input] = shorten
	mmh.setMu.Unlock()

	// 用于判断是否Hash冲突
	mmh.storeMu.Lock()
	mmh.urlSet[shorten] = struct{}{}
	mmh.storeMu.Unlock()
}

func (mmh *murmurHashGenerator) handleCollisions(input, shortCode string) string {
	for i := 0; i < maxRetries; i++ {

		mmh.setMu.RLock()
		_, ok := urlSet[shortCode]
		mmh.setMu.RUnlock()
		if !ok {
			return shortCode
		}
		shortCode = mmh.Generate(input + string(rune(i)))
	}
	return shortCode
}
