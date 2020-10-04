package okta

import (
	"testing"
)

func TestIsApiV1AppsEndpoint(t *testing.T) {
	apiV1AppsEndpoint := newRateLimitThrottle([]string{
		// the following endpoints share the same rate limit
		`/api/v1/apps$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/groups/(?P<GroupID>[[:alnum:]]+)$`,
		`/api/v1/apps/(?P<AppID>[[:alnum:]]+)/users$`,
	}, 5)
	if apiV1AppsEndpoint.checkIsEndpoint("/api/v1/apps") != true {
		t.Error()
	}
	if apiV1AppsEndpoint.checkIsEndpoint("/api/v1/apps/123/groups") != true {
		t.Error()
	}
	if apiV1AppsEndpoint.checkIsEndpoint("/api/v1/apps/123/users") != true {
		t.Error()
	}
	if apiV1AppsEndpoint.checkIsEndpoint("/api/v1/apps/123/groups/456") != true {
		t.Error()
	}
	if apiV1AppsEndpoint.checkIsEndpoint("/api/v1/apps/123") != false {
		t.Error()
	}
}
