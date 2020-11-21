package okta

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestRateLimitThrottleCheckIsEndpoint(t *testing.T) {
	apiV1AppsEndpoints := newRateLimitThrottle([]string{
		// the following endpoints share the same rate limit
		`/api/v1/apps$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups/(?P<GroupID>[[:alnum:]]+)$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/users$`,
	}, 5)
	if apiV1AppsEndpoints.checkIsEndpoint("/api/v1/apps") != true {
		t.Error()
	}
	if apiV1AppsEndpoints.checkIsEndpoint("/api/v1/apps/123/groups") != true {
		t.Error()
	}
	if apiV1AppsEndpoints.checkIsEndpoint("/api/v1/apps/123/users") != true {
		t.Error()
	}
	if apiV1AppsEndpoints.checkIsEndpoint("/api/v1/apps/123/groups/456") != true {
		t.Error()
	}
	if apiV1AppsEndpoints.checkIsEndpoint("/api/v1/apps/123") != false {
		t.Error()
	}
}

func TestRateLimitThrottlePostRequestHook(t *testing.T) {
	apiV1AppsEndpoints := newRateLimitThrottle([]string{
		// the following endpoints share the same rate limit
		`/api/v1/apps$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups/(?P<GroupID>[[:alnum:]]+)$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/users$`,
	}, 5)

	request := http.Request{
		URL: &url.URL{
			Path: "/api/v1/apps",
		},
	}
	rateLimit := 25
	rateLimitResetTime := time.Now()
	responseHeaders := http.Header{}
	responseHeaders.Add("X-Rate-Limit-Limit", fmt.Sprintf("%v", rateLimit))
	responseHeaders.Add("X-Rate-Limit-Reset", fmt.Sprintf("%v", rateLimitResetTime.Unix()))
	response := http.Response{
		Request: &request,
		Header:  responseHeaders,
	}
	apiV1AppsEndpoints.postRequestHook(&response)
	if apiV1AppsEndpoints.rateLimit != rateLimit {
		t.Errorf(
			"Rate limit header got parsed incorrectly. Expected %v, got: %v.",
			rateLimit,
			apiV1AppsEndpoints.rateLimit,
		)
	}
	if apiV1AppsEndpoints.rateLimitResetTime.Equal(rateLimitResetTime) {
		t.Errorf(
			"Rate limit reset time header got parsed incorrectly. Expected %v, got: %v.",
			rateLimitResetTime,
			apiV1AppsEndpoints.rateLimitResetTime,
		)
	}
}
