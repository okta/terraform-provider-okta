package apimutex

import (
	"math/rand"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		endPoint  string
		remaining []int
	}{
		{endPoint: "/api/v1/users", remaining: []int{1, 2, 3}},
		{endPoint: "/api/v1/users/foo", remaining: []int{6, 5, 4}},
		{endPoint: "/api/v1/apps", remaining: []int{7, 9, 8}},
		{endPoint: "/api/v1/apps/foo", remaining: []int{10, 12, 11}},
		{endPoint: "/api/v1/domains", remaining: []int{15, 13, 14}},
		{endPoint: "/api/v1/domains/foo", remaining: []int{16, 18, 17}},
	}

	limit := 500
	reset := (time.Now().Unix() + int64(60))
	for _, tc := range tests {
		// here, we are testing that regardless of threading parallelism the api
		// mutex will have the highest remaining value set for any given class
		// of endpoint
		amu := NewApiMutex()
		for _, remaining := range tc.remaining {
			go func(remaining int) {
				sleep := time.Duration(rand.Intn(100))
				time.Sleep(sleep * time.Millisecond)

				amu.Update(tc.endPoint, limit, remaining, reset)
			}(remaining)
		}
		time.Sleep(300 * time.Millisecond)

		maxRemaining := maxRemaining(tc.remaining)
		status := amu.Status(tc.endPoint)
		if maxRemaining != status.remaining {
			t.Fatalf("got %d, should be %d for the remaining value of %q's api status %+v", status.remaining, maxRemaining, tc.endPoint, status)
		}
	}
}

func TestGet(t *testing.T) {
	amu := NewApiMutex()
	if len(amu.status) != 0 {
		t.Fatalf("amu status map should be empty, but was sized %d", len(amu.status))
	}
	uris := []string{
		"/api/v1/users",
		"/api/v1/users",
		"/api/v1/apps",
		"/api/v1/apps",
		"/api/v1/domains",
		"/api/v1/domains",
	}
	for _, uri := range uris {
		_ = amu.get(uri)
	}
	if len(amu.status) != 3 {
		t.Fatalf("amu status map should sized 3 but was sized %d", len(amu.status))
	}
	keys := []string{
		"users",
		"apps",
		"other",
	}
	for _, key := range keys {
		if _, found := amu.status[key]; !found {
			t.Fatalf("amu should have status for key %q", key)
		}
	}
}

func TestNormalizeKey(t *testing.T) {
	tests := []struct {
		endPoint string
		expected string
	}{
		{endPoint: "/api/v1/users", expected: "users"},
		{endPoint: "/api/v1/users/foo", expected: "users"},
		{endPoint: "/api/v1/apps", expected: "apps"},
		{endPoint: "/api/v1/apps/foo", expected: "apps"},
		{endPoint: "/api/v1/authorizationServers", expected: "other"},
		{endPoint: "/api/v1/authorizationServers/foo", expected: "other"},
		{endPoint: "/api/v1/behaviors", expected: "other"},
		{endPoint: "/api/v1/behaviors/foo", expected: "other"},
		{endPoint: "/api/v1/domains", expected: "other"},
		{endPoint: "/api/v1/domains/foo", expected: "other"},
		{endPoint: "/api/v1/idps", expected: "other"},
		{endPoint: "/api/v1/idps/foo", expected: "other"},
		{endPoint: "/api/v1/internal", expected: "other"},
		{endPoint: "/api/v1/internal/foo", expected: "other"},
		{endPoint: "/api/v1/mappings", expected: "other"},
		{endPoint: "/api/v1/mappings/foo", expected: "other"},
		{endPoint: "/api/v1/meta", expected: "other"},
		{endPoint: "/api/v1/meta/foo", expected: "other"},
		{endPoint: "/api/v1/org", expected: "other"},
		{endPoint: "/api/v1/org/foo", expected: "other"},
		{endPoint: "/api/v1/policies", expected: "other"},
		{endPoint: "/api/v1/policies/foo", expected: "other"},
		{endPoint: "/api/v1/templates", expected: "other"},
		{endPoint: "/api/v1/templates/foo", expected: "other"},
	}

	amu := NewApiMutex()
	for _, tc := range tests {
		// test that private normalizedKey function is operating correctly
		key := amu.normalizeKey(tc.endPoint)
		if key != tc.expected {
			t.Fatalf("got %q, expected %q for endPoint %q", key, tc.expected, tc.endPoint)
		}
	}
}

func maxRemaining(remaining []int) int {
	var result int
	for _, r := range remaining {
		if r > result {
			result = r
		}
	}

	return result
}
