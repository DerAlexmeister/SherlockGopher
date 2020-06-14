package sherlockanalyser

import "sync"

type AnalyserCacheInterface interface {
	Request(addr string) bool
	Register(addr string)
}

type AnalyserCache struct {
	crawledLinks		[]string
	mutex				sync.Mutex
}


func (cache *AnalyserCache) Request(addr string) bool {
	return cache.contains(cache.crawledLinks, addr)
}

func (cache *AnalyserCache) Register(addr string) {
	cache.mutex.Lock()
	cache.crawledLinks = append(cache.crawledLinks, addr)
	cache.mutex.Unlock()
}

func (cache *AnalyserCache) contains(s []string, e string) bool {
	cache.mutex.Lock()
	ret := false
	for _, a := range s {
		if a == e {
			ret = true
		}
	}
	cache.mutex.Unlock()
	return ret
}

func NewAnalyserCache() AnalyserCache {
	return AnalyserCache{
		crawledLinks: make([]string, 0),
		mutex:        sync.Mutex{},
	}
}
