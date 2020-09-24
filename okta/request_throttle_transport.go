package okta

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	apiV1AppsEndpointCalls              = 0
	apiV1AppsEndpointRateLimitLimit     = 0
	apiV1AppsEndpointRateLimitResetTime time.Time
	apiV1AppsEndpointMux                sync.Mutex
	apiV1AppsPattern                    = regexp.MustCompile(
		`/api/v1/apps$`)
	apiV1AppsGroupsPattern = regexp.MustCompile(
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups$`)
	apiV1AppsUsersPattern = regexp.MustCompile(
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/users$`)
	apiV1AppsGroupPattern = regexp.MustCompile(
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups/(?P<GroupID>[[:alnum:]]+)$`)
)

func isApiV1AppsEndpoint(path string) bool {
	return len(apiV1AppsPattern.FindStringSubmatch(path))+
		len(apiV1AppsGroupsPattern.FindStringSubmatch(path))+
		len(apiV1AppsUsersPattern.FindStringSubmatch(path))+
		len(apiV1AppsGroupPattern.FindStringSubmatch(path)) != 0
}

type RequestThrottleTransport struct {
	base        http.RoundTripper
	maxRequests int
}

func (throttleTransport *RequestThrottleTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	appsEndpoint := isApiV1AppsEndpoint(req.URL.Path)
	if appsEndpoint {
		log.Println("[DEBUG] special isApiV1AppsEndpoint request throttle handling")
		apiV1AppsEndpointMux.Lock()
		apiV1AppsEndpointCalls += 1
		// TODO hardcoded 10 for testing, change it to calculating on the fly using
		// apiV1AppsEndpointRateLimitLimit * c.maxRequests / 100
		if apiV1AppsEndpointCalls >= 10 {
			apiV1AppsEndpointCalls = 0
			// add an extra margin to account for the clock skew
			timeToSleep := apiV1AppsEndpointRateLimitResetTime.Add(2 * time.Second).Sub(time.Now())
			if timeToSleep > 0 {
				log.Printf(
					"[INFO] Throttling /api/v1/apps requests, sleeping until rate limit reset for %s",
					timeToSleep,
				)
				time.Sleep(timeToSleep)
			}
		}
		apiV1AppsEndpointMux.Unlock()
	}
	resp, err := throttleTransport.base.RoundTrip(req)
	if err == nil && appsEndpoint {
		apiV1AppsEndpointMux.Lock()
		apiV1AppsEndpointRateLimitLimit, err = strconv.Atoi(resp.Header.Get("X-Rate-Limit-Limit"))
		if err != nil {
			// TODO
		}
		log.Printf("[DEBUG] /api/v1/apps rate limit limit: %v", apiV1AppsEndpointRateLimitLimit)
		resetTime, err := strconv.Atoi(resp.Header.Get("X-Rate-Limit-Reset"))
		if err != nil {
			// TODO
		}
		futureResetTime := time.Unix(int64(resetTime), 0)
		log.Printf("[DEBUG] future /api/v1/apps rate limit reset time: %v", futureResetTime)
		if futureResetTime != apiV1AppsEndpointRateLimitResetTime {
			apiV1AppsEndpointCalls = 1
		}
		apiV1AppsEndpointRateLimitResetTime = futureResetTime
		apiV1AppsEndpointMux.Unlock()
	}
	return resp, err
}

func NewRequestThrottleTransport(base http.RoundTripper, maxRequests int) *RequestThrottleTransport {
	return &RequestThrottleTransport{base: base, maxRequests: maxRequests}
}
