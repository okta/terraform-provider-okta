package apimutex

import (
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestHasCapacity(t *testing.T) {
	amu, err := NewAPIMutex(50)
	if err != nil {
		t.Fatalf("api mutex constructor had error %+v", err)
	}

	endPoint := "/api/v1/users"
	reset := (time.Now().Unix() + int64(60))
	// endpoint, limit, remaining, reset
	amu.Update(http.MethodGet, endPoint, 90, 46, reset)
	if !amu.HasCapacity(http.MethodGet, endPoint) {
		t.Fatalf("ami mutex should have capacity, 50%% threshold, 90 limit, 46 remaining")
	}

	amu.Update(http.MethodGet, endPoint, 90, 45, reset)
	if !amu.HasCapacity(http.MethodGet, endPoint) {
		t.Fatalf("ami mutex should have capacity, 50%% threshold, 90 limit, 45 remaining")
	}

	amu.Update(http.MethodGet, endPoint, 90, 44, reset)
	if amu.HasCapacity(http.MethodGet, endPoint) {
		t.Fatalf("ami mutex shouldn't have capacity, 50%% threshold, 90 limit, 44 remaining")
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		method    string
		endPoint  string
		remaining []int
	}{
		{method: http.MethodGet, endPoint: "/api/v1/apps", remaining: []int{1, 2, 3}},
		{method: http.MethodGet, endPoint: "/api/v1/apps/foo", remaining: []int{4, 5, 6}},
		{method: http.MethodGet, endPoint: "/api/v1/certificateAuthorities", remaining: []int{7, 8, 9}},
		{method: http.MethodGet, endPoint: "/oauth2/v1/clients", remaining: []int{10, 11, 12}},
		{method: http.MethodGet, endPoint: "/api/v1/devices", remaining: []int{13, 14, 15}},
		{method: http.MethodGet, endPoint: "/api/v1/events", remaining: []int{16, 17, 18}},
		{method: http.MethodGet, endPoint: "/api/v1/groups", remaining: []int{19, 20, 21}},
		{method: http.MethodGet, endPoint: "/api/v1/groups/foo", remaining: []int{22, 23, 24}},
		{method: http.MethodGet, endPoint: "/api/v1/logs", remaining: []int{25, 26, 27}},
		{method: http.MethodGet, endPoint: "/api/v1/users", remaining: []int{28, 29, 30}},
		{method: http.MethodGet, endPoint: "/api/v1/users/foo", remaining: []int{31, 32, 33}},
		{method: http.MethodPost, endPoint: "/api/v1/users/foo", remaining: []int{34, 35, 36}},
		{method: http.MethodGet, endPoint: "/api/v1/domains", remaining: []int{37, 38, 39}},
	}

	limit := 500
	reset := (time.Now().Unix() + int64(60))
	for _, tc := range tests {
		// here, we are testing that regardless of threading parallelism the api
		// mutex will have the highest remaining value set for any given class
		// of endpoint
		amu, err := NewAPIMutex(100)
		if err != nil {
			t.Fatalf("api mutex constructor had error %+v", err)
		}
		for _, remaining := range tc.remaining {
			go func(remaining int) {
				sleep := time.Duration(rand.Intn(100))
				time.Sleep(sleep * time.Millisecond)

				amu.Update(tc.method, tc.endPoint, limit, remaining, reset)
			}(remaining)
		}
		time.Sleep(300 * time.Millisecond)

		minRemaining := minRemaining(tc.remaining)
		status := amu.Status(tc.method, tc.endPoint)
		if minRemaining != status.remaining {
			t.Fatalf("got %d, should be %d of %+v for the remaining value of %q's api status %+v", status.remaining, minRemaining, tc.remaining, tc.endPoint, status)
		}
	}
}

func TestGet(t *testing.T) {
	amu, err := NewAPIMutex(100)
	if err != nil {
		t.Fatalf("api mutex constructor had error %+v", err)
	}
	if len(amu.status) != 13 {
		t.Fatalf("amu status map should sized 13 but was sized %d", len(amu.status))
	}
	keys := []string{
		"users",
		"user-id",
		"apps",
		"app-id",
		"groups",
		"group-id",
		"other",
	}
	for _, key := range keys {
		if _, found := amu.status[key]; !found {
			t.Fatalf("amu should have status for key %q", key)
		}
	}
}

func TestNormalizeKey(t *testing.T) {
	// Attempts to cover the rules listed at:
	// https://developer.okta.com/docs/reference/rl-global-mgmt/
	tests := []struct {
		method   string
		endPoint string
		expected string
	}{
		// Create or list applications:
		// /api/v1/apps except /api/v1/apps/${id}
		{method: http.MethodGet, endPoint: "/api/v1/apps", expected: "apps"},
		{method: http.MethodPost, endPoint: "/api/v1/apps", expected: "apps"},

		// Get, update, or delete an application by ID:
		// /api/v1/apps/${id} only
		{method: http.MethodGet, endPoint: "/api/v1/apps/XXX", expected: "app-id"},
		{method: http.MethodPut, endPoint: "/api/v1/apps/XXX", expected: "app-id"},
		{method: http.MethodPost, endPoint: "/api/v1/apps/XXX", expected: "app-id"},
		{method: http.MethodDelete, endPoint: "/api/v1/apps/XXX", expected: "app-id"},

		// Create or list groups:
		// /api/v1/groups except /api/v1/groups/${id}
		{method: http.MethodGet, endPoint: "/api/v1/groups", expected: "groups"},
		{method: http.MethodPost, endPoint: "/api/v1/groups", expected: "groups"},

		// Get, update, or delete a group by ID:
		// /api/v1/groups/${id} only
		{method: http.MethodGet, endPoint: "/api/v1/groups/XXX", expected: "group-id"},
		{method: http.MethodPost, endPoint: "/api/v1/groups/XXX", expected: "group-id"},
		{method: http.MethodPut, endPoint: "/api/v1/groups/XXX", expected: "group-id"},
		{method: http.MethodDelete, endPoint: "/api/v1/groups/XXX", expected: "group-id"},

		// Create or list users:
		// Only GET or POST to /api/v1/users
		{method: http.MethodGet, endPoint: "/api/v1/users", expected: "users"},
		{method: http.MethodPost, endPoint: "/api/v1/users", expected: "users"},

		// Update or delete a user by ID:
		// Only POST, PUT or DELETE to /api/v1/users/${id}
		{method: http.MethodPost, endPoint: "/api/v1/users/foo", expected: "user-id"},
		{method: http.MethodPut, endPoint: "/api/v1/users/foo", expected: "user-id"},
		{method: http.MethodDelete, endPoint: "/api/v1/users/foo", expected: "user-id"},

		// Get System Log data:
		// /api/v1/logs
		{method: http.MethodGet, endPoint: "/api/v1/logs", expected: "logs"},

		// Get event data:
		// /api/v1/events
		{method: http.MethodGet, endPoint: "/api/v1/events", expected: "events"},

		// OAuth2 client configuration requests:
		// /oauth2/v1/clients
		{method: http.MethodGet, endPoint: "/oauth2/v1/clients", expected: "clients"},

		// Get a user by ID or sign in
		// Only GET to /api/v1/users/${idOrLogin}
		{method: http.MethodGet, endPoint: "/api/v1/users/foo", expected: "user-id-get"},

		// /api/v1/certificateAuthorities
		{method: http.MethodGet, endPoint: "/api/v1/certificateAuthorities", expected: "cas-id"},

		// /api/v1/devices
		{method: http.MethodGet, endPoint: "/api/v1/devices", expected: "devices"},

		// Most other API actions:
		// /api/v1
		{method: http.MethodGet, endPoint: "/api/v1/apps/XXX/users", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/authorizationServers", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/authorizationServers/XXX", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/behaviors", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/behaviors/XXX", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/domains", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/domains/XXX", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/groups/rules/XXX", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/groups/XXX/roles/XXX", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/idps", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/idps/XXX", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/internal", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/internal/XXX", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/mappings", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/mappings/foo", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/meta", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/meta/foo", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/org", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/org/foo", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/policies", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/policies/foo", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/templates", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/templates/foo", expected: "other"},
		{method: http.MethodGet, endPoint: "/api/v1/trustedOrigins/foo", expected: "other"},
	}

	amu, err := NewAPIMutex(100)
	if err != nil {
		t.Fatalf("api mutex constructor had error %+v", err)
	}
	for _, tc := range tests {
		// test that private normalizedKey function is operating correctly
		key := amu.normalizeKey(tc.method, tc.endPoint)
		if key != tc.expected {
			t.Fatalf("got %q, expected %q for method: %q, endPoint: %q", key, tc.expected, tc.method, tc.endPoint)
		}
	}
}

func minRemaining(remaining []int) int {
	var result int
	first := true
	for _, r := range remaining {
		if first {
			result = r
			first = false
			continue
		}
		if r < result {
			result = r
		}
	}

	return result
}
