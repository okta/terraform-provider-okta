package okta

import (
	"testing"
)

func TestIsApiV1AppsEndpoint(t *testing.T) {
	if isApiV1AppsEndpoint("/api/v1/apps") != true {
		t.Error()
	}
	if isApiV1AppsEndpoint("/api/v1/apps/123/groups") != true {
		t.Error()
	}
	if isApiV1AppsEndpoint("/api/v1/apps/123/users") != true {
		t.Error()
	}
	if isApiV1AppsEndpoint("/api/v1/apps/123/groups/456") != true {
		t.Error()
	}
	if isApiV1AppsEndpoint("/api/v1/apps/123") != false {
		t.Error()
	}
}
