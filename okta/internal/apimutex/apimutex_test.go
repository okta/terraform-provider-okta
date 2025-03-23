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
				// lintignore:R018
				time.Sleep(sleep * time.Millisecond)

				amu.Update(tc.method, tc.endPoint, limit, remaining, reset)
			}(remaining)
		}
		// lintignore:R018
		time.Sleep(300 * time.Millisecond)

		minRemaining := minRemaining(tc.remaining)
		status := amu.Status(tc.method, tc.endPoint)
		if minRemaining != status.remaining {
			t.Fatalf("got %d, should be %d of %+v for the remaining value of %q's api status %+v", status.remaining, minRemaining, tc.remaining, tc.endPoint, status)
		}
	}
}

func TestGetRegex(t *testing.T) {
	tests := []struct {
		method         string
		endPoint       string
		expectedBucket string
	}{
		{http.MethodGet, "/d/n/e", "/"},
		{http.MethodGet, "/.well-known/acme-challenge/ID", "/.well-known/acme-challenge/{token}"},
		{http.MethodGet, "/api/v1/groups", "/api/v1/groups"},
		{http.MethodPost, "/api/v1/groups", "/api/v1/groups"},
		{http.MethodGet, "/api/v1/authorizationServers/0123456789abcdefghij/policies/0123456789abcdefghij", "/api/v1"},
	}
	amu, err := NewAPIMutex(100)
	if err != nil {
		t.Fatalf("api mutex constructor had error %+v", err)
	}
	for _, test := range tests {
		result := amu.get(test.method, test.endPoint)
		if result != amu.status[test.expectedBucket] {
			t.Fatalf("expected endpoint \"%s %s\" to be in status bucket %q", test.method, test.endPoint, test.expectedBucket)
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
