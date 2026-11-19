// Package apprefresh shares complete app collections between distinct resource
// reads. Entries are scoped to one provider configuration, never the process.
package apprefresh

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"
)

type entry struct {
	values   map[string]any
	consumed map[string]bool
	expires  time.Time
}

type Cache struct {
	mu         sync.Mutex
	entries    map[string]*entry
	loading    map[string]chan struct{}
	generation uint64
	writers    int
	now        func() time.Time
}

func (c *Cache) clock() time.Time {
	if c.now != nil {
		return c.now()
	}
	return time.Now()
}

// Get shares a collection only across different IDs. Reading an ID again
// reloads the collection rather than reusing a previous refresh's snapshot.
// Failed or incomplete collections must be returned as errors by load.
func (c *Cache) Get(ctx context.Context, key, id string, load func(context.Context) (map[string]any, error)) (any, bool, error) {
	for {
		if err := ctx.Err(); err != nil {
			return nil, false, err
		}
		c.mu.Lock()
		if e := c.entries[key]; e != nil && c.writers == 0 && c.clock().Before(e.expires) && !e.consumed[id] {
			e.consumed[id] = true
			v, ok := e.values[id]
			c.mu.Unlock()
			return v, ok, nil
		}
		if done := c.loading[key]; done != nil {
			c.mu.Unlock()
			select {
			case <-ctx.Done():
				return nil, false, ctx.Err()
			case <-done:
				continue
			}
		}
		if c.loading == nil {
			c.loading = make(map[string]chan struct{})
		}
		done := make(chan struct{})
		c.loading[key] = done
		generation := c.generation
		c.mu.Unlock()
		values, err := load(ctx)
		c.mu.Lock()
		if err == nil && c.writers == 0 && c.generation == generation {
			if c.entries == nil {
				c.entries = make(map[string]*entry)
			}
			c.entries[key] = &entry{values: values, consumed: map[string]bool{id: true}, expires: c.clock().Add(30 * time.Second)}
		}
		delete(c.loading, key)
		close(done)
		c.mu.Unlock()
		if err != nil {
			return nil, false, err
		}
		v, ok := values[id]
		return v, ok, nil
	}
}

// Invalidate both before and after a write, including failed writes: a response
// failure does not establish that the server made no change.
func (c *Cache) beginWrite() func() {
	c.mu.Lock()
	c.writers++
	c.generation++
	c.entries = nil
	c.mu.Unlock()
	return func() { c.mu.Lock(); c.writers--; c.generation++; c.entries = nil; c.mu.Unlock() }
}

// Written reports whether this provider configuration has attempted a management
// write. Subsequent resource reads use individual endpoints for apply convergence.
func (c *Cache) Written() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.generation != 0
}

type invalidatingTransport struct {
	base  http.RoundTripper
	cache *Cache
}

func InvalidateWrites(base http.RoundTripper, cache *Cache) http.RoundTripper {
	if cache == nil {
		return base
	}
	if base == nil {
		base = http.DefaultTransport
	}
	return &invalidatingTransport{base: base, cache: cache}
}
func (t *invalidatingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.HasPrefix(req.URL.Path, "/api/v1/") && req.Method != http.MethodGet && req.Method != http.MethodHead && req.Method != http.MethodOptions {
		end := t.cache.beginWrite()
		defer end()
	}
	return t.base.RoundTrip(req)
}
