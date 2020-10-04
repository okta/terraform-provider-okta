package okta

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func isAPIV1AppsEndpoint(path string) bool {
	endpointsPatterns := []*regexp.Regexp{
		// apps pattern
		regexp.MustCompile(`/api/v1/apps$`),
		// groups pattern
		regexp.MustCompile(`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups$`),
		// group pattern
		regexp.MustCompile(`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups/(?P<GroupID>[[:alnum:]]+)$`),
		// users pattern
		regexp.MustCompile(`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/users$`),
	}
	// Check if endpoint match to one of the patterns.
	var foundLen int
	for _, p := range endpointsPatterns {
		foundLen += len(p.FindStringSubmatch(path))
	}
	return foundLen != 0
}

type requestThrottleTransport struct {
	base                                http.RoundTripper
	percentageOfLimitRate               int
	apiV1AppsEndpointCalls              int
	apiV1AppsEndpointRateLimit          int
	apiV1AppsEndpointRateLimitResetTime time.Time
	sync.Mutex
}

func (tt *requestThrottleTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	// Nothing to throttle, just bypass.
	if !isAPIV1AppsEndpoint(req.URL.Path) {
		return tt.base.RoundTrip(req)
	}

	// Perform throttling (wait with making a request).
	log.Println("[DEBUG] special isApiV1AppsEndpoint request throttle handling")
	tt.Lock()
	tt.apiV1AppsEndpointCalls++
	// TODO hardcoded 10 for testing, change it to calculating on the fly using
	// apiV1AppsEndpointRateLimit * c.maxRequests / 100
	if tt.apiV1AppsEndpointCalls >= 10 {
		tt.apiV1AppsEndpointCalls = 0
		// add an extra margin to account for the clock skew
		timeToSleep := tt.apiV1AppsEndpointRateLimitResetTime.Add(2 * time.Second).Sub(time.Now())
		if timeToSleep > 0 {
			log.Printf(
				"[INFO] Throttling /api/v1/apps requests, sleeping until rate limit reset for %s",
				timeToSleep,
			)
			time.Sleep(timeToSleep)
		}
	}
	tt.Unlock()

	// Make a request after throttling phase.
	resp, err := tt.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	// Above throttled request ended successfully, update information about limits in throttler.
	tt.Lock()
	tt.apiV1AppsEndpointRateLimit, err = strconv.Atoi(resp.Header.Get("X-Rate-Limit-Limit"))
	if err != nil {
		// TODO
	}
	log.Printf("[DEBUG] /api/v1/apps rate limit limit: %v", tt.apiV1AppsEndpointRateLimit)
	resetTime, err := strconv.Atoi(resp.Header.Get("X-Rate-Limit-Reset"))
	if err != nil {
		// TODO
	}
	futureResetTime := time.Unix(int64(resetTime), 0)
	log.Printf("[DEBUG] future /api/v1/apps rate limit reset time: %v", futureResetTime)
	if futureResetTime != tt.apiV1AppsEndpointRateLimitResetTime {
		tt.apiV1AppsEndpointCalls = 1
	}
	tt.apiV1AppsEndpointRateLimitResetTime = futureResetTime
	tt.Unlock()

	return resp, nil
}

// NewRequestThrottleTransport returns RoundTripper which provides throttling according to maxRequests.
// Every time new instance is returned which does not share any state with any of returned previously.
// Hence for every Okta API client instanced for particular Okta Organization the same throttler should be used.
func NewRequestThrottleTransport(base http.RoundTripper, percentageOfLimitRate int) *requestThrottleTransport {
	return &requestThrottleTransport{base: base, percentageOfLimitRate: percentageOfLimitRate}
}
