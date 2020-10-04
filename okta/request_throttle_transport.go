package okta

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type rateLimitThrottle struct {
	endpointsPatterns  []*regexp.Regexp
	noRequestsMade     int
	maxRequests        int
	rateLimit          int
	rateLimitResetTime time.Time
	sync.Mutex
}

func newRateLimitThrottle(endpointsRegexes []string, maxRequests int) rateLimitThrottle {
	endpointsPatterns := make([]*regexp.Regexp, len(endpointsRegexes))
	for i, endpointRegex := range endpointsRegexes {
		endpointsPatterns[i] = regexp.MustCompile(endpointRegex)
	}
	return rateLimitThrottle{
		endpointsPatterns: endpointsPatterns,
		maxRequests:       maxRequests,
	}
}

func (throttle *rateLimitThrottle) checkIsEndpoint(path string) bool {
	for _, pattern := range throttle.endpointsPatterns {
		if len(pattern.FindStringSubmatch(path)) != 0 {
			return true
		}
	}
	return false
}

func (throttle *rateLimitThrottle) preRequestHook(path string) {
	if !throttle.checkIsEndpoint(path) {
		return
	}
	log.Println("[DEBUG] special preRequestHook request throttle handling")
	throttle.Lock()
	defer throttle.Unlock()
	throttle.noRequestsMade += 1
	if throttle.rateLimit != 0 && throttle.noRequestsMade >= (throttle.rateLimit*throttle.maxRequests/100) {
		throttle.noRequestsMade = 1
		// add an extra margin to account for the clock skew
		timeToSleep := throttle.rateLimitResetTime.Add(2 * time.Second).Sub(time.Now())
		if timeToSleep > 0 {
			log.Printf(
				"[INFO] Throttling %s requests, sleeping until rate limit reset for %s",
				path,
				timeToSleep,
			)
			time.Sleep(timeToSleep)
		}
	}
}

func (throttle *rateLimitThrottle) postRequestHook(resp *http.Response) {
	if !throttle.checkIsEndpoint(resp.Request.URL.Path) {
		return
	}
	throttle.Lock()
	defer throttle.Unlock()
	var err error
	throttle.rateLimit, err = strconv.Atoi(resp.Header.Get("X-Rate-Limit-Limit"))
	if err != nil {
		// TODO
	}
	log.Printf("[DEBUG] %s rate limit limit: %v", resp.Request.URL.Path, throttle.rateLimit)
	resetTime, err := strconv.Atoi(resp.Header.Get("X-Rate-Limit-Reset"))
	if err != nil {
		// TODO
	}
	futureResetTime := time.Unix(int64(resetTime), 0)
	log.Printf("[DEBUG] future %s rate limit reset time: %v", resp.Request.URL.Path, futureResetTime)
	if futureResetTime != throttle.rateLimitResetTime {
		log.Printf("[DEBUG] %s rate limit reset detected", resp.Request.URL.Path)
		throttle.noRequestsMade = 1
	}
	throttle.rateLimitResetTime = futureResetTime
}

type requestThrottleTransport struct {
	base               http.RoundTripper
	throttledEndpoints []*rateLimitThrottle
}

func (transport *requestThrottleTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	for i, _ := range transport.throttledEndpoints {
		transport.throttledEndpoints[i].preRequestHook(req.URL.Path)
	}
	resp, err := transport.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	for i, _ := range transport.throttledEndpoints {
		transport.throttledEndpoints[i].postRequestHook(resp)
	}
	return resp, err
}

// NewRequestThrottleTransport returns RoundTripper which provides throttling according to maxRequests.
// Every new instance returned has its own local state. Hence for every Okta API client instanced for
// particular Okta Organization the same throttler should be used.
func NewRequestThrottleTransport(base http.RoundTripper, maxRequests int) *requestThrottleTransport {
	apiV1AppsEndpoint := newRateLimitThrottle([]string{
		// the following endpoints share the same rate limit
		`/api/v1/apps$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups/(?P<GroupID>[[:alnum:]]+)$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/users$`,
	}, maxRequests)
	return &requestThrottleTransport{
		base:               base,
		throttledEndpoints: []*rateLimitThrottle{&apiV1AppsEndpoint},
	}
}
