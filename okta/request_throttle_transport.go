package okta

import (
	"context"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type rateLimitThrottle struct {
	endpointsPatterns  []*regexp.Regexp
	noOfRequestsMade   int
	maxRequests        int
	rateLimit          int
	rateLimitResetTime time.Time
	sync.Mutex
}

func newRateLimitThrottle(endpointsRegexes []string, maxRequests int) *rateLimitThrottle {
	endpointsPatterns := make([]*regexp.Regexp, len(endpointsRegexes))
	for i, endpointRegex := range endpointsRegexes {
		endpointsPatterns[i] = regexp.MustCompile(endpointRegex)
	}
	return &rateLimitThrottle{
		endpointsPatterns: endpointsPatterns,
		maxRequests:       maxRequests,
	}
}

func (t *rateLimitThrottle) checkIsEndpoint(path string) bool {
	for _, pattern := range t.endpointsPatterns {
		if len(pattern.FindStringSubmatch(path)) > 0 {
			return true
		}
	}
	return false
}

func (t *rateLimitThrottle) preRequestHook(ctx context.Context, path string) error {
	if !t.checkIsEndpoint(path) {
		return nil
	}
	log.Println("[DEBUG] special preRequestHook request throttle handling")
	t.Lock()
	defer t.Unlock()
	t.noOfRequestsMade++
	if t.rateLimit != 0 && float64(t.noOfRequestsMade) > math.Max(float64(t.rateLimit*t.maxRequests)/100.0, 1) {
		t.noOfRequestsMade = 1
		// add an extra margin to account for the clock skew
		timeToSleep := time.Until(t.rateLimitResetTime.Add(2 * time.Second))
		if timeToSleep > 0 {
			log.Printf(
				"[INFO] Throttling %s requests, sleeping for %s until rate limit reset", path, timeToSleep)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.NewTimer(timeToSleep).C:
				return nil
			}
		}
	}
	return nil
}

func (t *rateLimitThrottle) postRequestHook(resp *http.Response) {
	if !t.checkIsEndpoint(resp.Request.URL.Path) {
		return
	}
	t.Lock()
	defer t.Unlock()
	rateLimit, err := strconv.Atoi(resp.Header.Get("X-Rate-Limit-Limit"))
	if err != nil {
		log.Printf("[WARN] X-Rate-Limit-Limit response header is missing or invalid, skipping postRequestHook: %v", err)
		return
	}
	t.rateLimit = rateLimit
	log.Printf("[DEBUG] %s rate limit: %d", resp.Request.URL.Path, t.rateLimit)
	resetTime, err := strconv.Atoi(resp.Header.Get("X-Rate-Limit-Reset"))
	if err != nil {
		log.Printf("[WARN] X-Rate-Limit-Reset response header is missing or invalid, skipping postRequestHook: %v", err)
		return
	}
	futureResetTime := time.Unix(int64(resetTime), 0)
	if !t.rateLimitResetTime.IsZero() && futureResetTime != t.rateLimitResetTime {
		log.Printf("[DEBUG] %s rate limit reset detected", resp.Request.URL.Path)
		t.noOfRequestsMade = 1
	}
	log.Printf("[DEBUG] future %s rate limit reset time: %v", resp.Request.URL.Path, futureResetTime)
	t.rateLimitResetTime = futureResetTime
}

type requestThrottleTransport struct {
	base               http.RoundTripper
	throttledEndpoints []*rateLimitThrottle
}

func (t *requestThrottleTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, endpoint := range t.throttledEndpoints {
		if err := endpoint.preRequestHook(req.Context(), req.URL.Path); err != nil {
			return nil, err
		}
	}
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	for _, endpoint := range t.throttledEndpoints {
		endpoint.postRequestHook(resp)
	}
	return resp, nil
}

// NewRequestThrottleTransport returns RoundTripper which provides throttling according to maxRequests.
// Every new instance returned has its own local state. Hence for every Okta API client instanced for
// particular Okta Organization the same throttler should be used.
func NewRequestThrottleTransport(base http.RoundTripper, maxRequests int) *requestThrottleTransport {
	apiV1AppsEndpoints := newRateLimitThrottle([]string{
		// the following endpoints share the same rate limit
		`/api/v1/apps$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups/(?P<GroupID>[[:alnum:]]+)$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/users$`,
	}, maxRequests)
	return &requestThrottleTransport{
		base:               base,
		throttledEndpoints: []*rateLimitThrottle{apiV1AppsEndpoints},
	}
}
