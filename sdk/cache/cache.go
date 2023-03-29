package cache

import (
	"io"
	"net/http"
)

type Cache interface {
	Get(key string) *http.Response
	Set(key string, value *http.Response)
	GetString(key string) string
	SetString(key, value string)
	Delete(key string)
	Clear()
	Has(key string) bool
}

func CreateCacheKey(req *http.Request) string {
	s := req.URL.Scheme + "://" + req.URL.Host + req.URL.RequestURI()
	return s
}

func CopyResponse(resp *http.Response) *http.Response {
	c := *resp
	c.Body = io.NopCloser(resp.Body)

	return &c
}
