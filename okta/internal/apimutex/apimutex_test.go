package apimutex

import (
	"math/rand"
	"testing"
	"time"
)

func TestHasCapacity(t *testing.T) {
	amu, err := NewApiMutex(50)
	if err != nil {
		t.Fatalf("api mutex constructor had error %+v", err)
	}

	endPoint := "/api/v1/users"
	reset := (time.Now().Unix() + int64(60))
	// endpoint, limit, remaining, reset
	amu.Update(endPoint, 90, 46, reset)
	if !amu.HasCapacity(endPoint) {
		t.Fatalf("ami mutex should have capacity, 50%% threshold, 90 limit, 46 remaining")
	}

	amu.Update(endPoint, 90, 45, reset)
	if !amu.HasCapacity(endPoint) {
		t.Fatalf("ami mutex should have capacity, 50%% threshold, 90 limit, 45 remaining")
	}

	amu.Update(endPoint, 90, 44, reset)
	if amu.HasCapacity(endPoint) {
		t.Fatalf("ami mutex shouldn't have capacity, 50%% threshold, 90 limit, 44 remaining")
	}
}

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
		amu, err := NewApiMutex(100)
		if err != nil {
			t.Fatalf("api mutex constructor had error %+v", err)
		}
		for _, remaining := range tc.remaining {
			go func(remaining int) {
				sleep := time.Duration(rand.Intn(100))
				time.Sleep(sleep * time.Millisecond)

				amu.Update(tc.endPoint, limit, remaining, reset)
			}(remaining)
		}
		time.Sleep(300 * time.Millisecond)

		minRemaining := minRemaining(tc.remaining)
		status := amu.Status(tc.endPoint)
		if minRemaining != status.remaining {
			t.Fatalf("got %d, should be %d of %+v for the remaining value of %q's api status %+v", status.remaining, minRemaining, tc.remaining, tc.endPoint, status)
		}
	}
}

func TestGet(t *testing.T) {
	amu, err := NewApiMutex(100)
	if err != nil {
		t.Fatalf("api mutex constructor had error %+v", err)
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

	amu, err := NewApiMutex(100)
	if err != nil {
		t.Fatalf("api mutex constructor had error %+v", err)
	}
	for _, tc := range tests {
		// test that private normalizedKey function is operating correctly
		key := amu.normalizeKey(tc.endPoint)
		if key != tc.expected {
			t.Fatalf("got %q, expected %q for endPoint %q", key, tc.expected, tc.endPoint)
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
