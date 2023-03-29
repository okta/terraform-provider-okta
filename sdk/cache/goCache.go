package cache

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"
	"time"

	patrickmnGoCache "github.com/patrickmn/go-cache"
)

type GoCache struct {
	ttl         time.Duration
	tti         time.Duration
	rootLibrary *patrickmnGoCache.Cache
}

func NewGoCache(ttl, tti int32) GoCache {
	c := patrickmnGoCache.New(time.Duration(ttl)*time.Second, time.Duration(tti)*time.Second)

	gc := GoCache{
		ttl:         time.Duration(ttl) * time.Second,
		tti:         time.Duration(tti) * time.Second,
		rootLibrary: c,
	}

	return gc
}

func (c GoCache) Get(key string) *http.Response {
	item, found := c.rootLibrary.Get(key)
	if found {
		r := bufio.NewReader(bytes.NewReader(item.([]byte)))
		resp, _ := http.ReadResponse(r, nil)
		return resp
	}

	return nil
}

func (c GoCache) Set(key string, value *http.Response) {
	cacheableResponse, _ := httputil.DumpResponse(value, true)

	c.rootLibrary.Set(key, cacheableResponse, c.ttl)
}

func (c GoCache) GetString(key string) string {
	item, found := c.rootLibrary.Get(key)
	if found {
		return item.(string)
	}

	return ""
}

func (c GoCache) SetString(key, value string) {
	c.rootLibrary.Set(key, value, c.ttl)
}

func (c GoCache) Delete(key string) {
	c.rootLibrary.Delete(key)
}

func (c GoCache) Clear() {
	c.rootLibrary.Flush()
}

func (c GoCache) Has(key string) bool {
	_, found := c.rootLibrary.Get(key)
	return found
}
