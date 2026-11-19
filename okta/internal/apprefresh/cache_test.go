package apprefresh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBatchConcurrentReads(t *testing.T) {
	var c Cache
	var calls atomic.Int32
	load := func(context.Context) (map[string]any, error) {
		calls.Add(1)
		m := map[string]any{}
		for i := range 100 {
			m[fmt.Sprint(i)] = i
		}
		return m, nil
	}
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			v, ok, err := c.Get(t.Context(), "app", fmt.Sprint(i), load)
			if err != nil || !ok || v != i {
				t.Errorf("read %d: %v %v %v", i, v, ok, err)
			}
		})
	}
	wg.Wait()
	if calls.Load() != 1 {
		t.Fatalf("100 distinct reads made %d collection loads", calls.Load())
	}
	if _, _, err := c.Get(t.Context(), "app", "0", load); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 {
		t.Fatal("repeated resource read reused old snapshot")
	}
}

func TestBatchExpiryAndFailedLoad(t *testing.T) {
	now := time.Now()
	c := Cache{now: func() time.Time { return now }}
	calls := 0
	load := func(context.Context) (map[string]any, error) { calls++; return map[string]any{"a": 1, "b": 2}, nil }
	c.Get(t.Context(), "app", "a", load)
	now = now.Add(31 * time.Second)
	c.Get(t.Context(), "app", "b", load)
	if calls != 2 {
		t.Fatal("expired collection reused")
	}
	bad := func(context.Context) (map[string]any, error) {
		return map[string]any{"a": 1}, errors.New("page two failed")
	}
	if _, _, err := c.Get(t.Context(), "other", "a", bad); err == nil {
		t.Fatal("partial page accepted")
	}
	if _, ok, err := c.Get(t.Context(), "other", "a", load); !ok || err != nil {
		t.Fatal("failed load cached")
	}
}

func TestBatchCanceledWaiter(t *testing.T) {
	var c Cache
	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	load := func(context.Context) (map[string]any, error) { close(started); <-release; return map[string]any{}, nil }
	go func() { defer close(done); c.Get(t.Context(), "app", "a", load) }()
	<-started
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, _, err := c.Get(ctx, "app", "b", load); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	close(release)
	<-done
}

type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestBatchWriteInvalidation(t *testing.T) {
	var c Cache
	calls := 0
	load := func(context.Context) (map[string]any, error) {
		calls++
		return map[string]any{"a": calls, "b": calls}, nil
	}
	c.Get(t.Context(), "app", "a", load)
	transport := InvalidateWrites(roundTripper(func(*http.Request) (*http.Response, error) {
		// A read during a write cannot consume the previous snapshot or seed a new one.
		c.Get(t.Context(), "app", "b", load)
		return nil, errors.New("response lost after write")
	}), &c)
	req, _ := http.NewRequest(http.MethodPut, "https://example.test/api/v1/apps/app", nil)
	transport.RoundTrip(req)
	v, _, err := c.Get(t.Context(), "app", "b", load)
	if err != nil || v != 3 || calls != 3 {
		t.Fatalf("write reused stale snapshot: %v %d %v", v, calls, err)
	}
}

func TestBatchWriteDuringLoad(t *testing.T) {
	var c Cache
	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		c.Get(t.Context(), "app", "a", func(context.Context) (map[string]any, error) {
			close(started)
			<-release
			return map[string]any{"b": "old"}, nil
		})
	}()
	<-started
	end := c.beginWrite()
	end()
	close(release)
	<-done
	v, _, err := c.Get(t.Context(), "app", "b", func(context.Context) (map[string]any, error) { return map[string]any{"b": "new"}, nil })
	if err != nil || v != "new" {
		t.Fatal("load crossed write boundary")
	}
}

func TestBatchTokenAcquisitionKeepsCollections(t *testing.T) {
	var c Cache
	transport := InvalidateWrites(roundTripper(func(*http.Request) (*http.Response, error) { return &http.Response{StatusCode: 200}, nil }), &c)
	req, _ := http.NewRequest(http.MethodPost, "https://example.test/oauth2/v1/token", nil)
	transport.RoundTrip(req)
	if c.Written() {
		t.Fatal("OAuth token acquisition disabled resource batching")
	}
}
