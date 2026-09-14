package idaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/okta/terraform-provider-okta/sdk"
)

// pagedUserGroupsServer serves GET /api/v1/users/{id}/groups from the given
// pages, emitting the same Link headers Okta uses for cursor pagination.
func pagedUserGroupsServer(t *testing.T, pages [][]string) *httptest.Server {
	t.Helper()

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := 0
		if after := r.URL.Query().Get("after"); after != "" {
			for i, ids := range pages {
				if ids[len(ids)-1] == after {
					page = i + 1
					break
				}
			}
		}
		if page >= len(pages) {
			t.Errorf("unexpected page request: %s", r.URL.String())
			http.Error(w, "unexpected page", http.StatusBadRequest)
			return
		}

		groups := make([]*sdk.Group, 0, len(pages[page]))
		for _, id := range pages[page] {
			groups = append(groups, &sdk.Group{Id: id, Profile: &sdk.GroupProfile{Name: "group " + id}})
		}

		w.Header().Add("Link", fmt.Sprintf(`<%s%s?limit=200>; rel="self"`, server.URL, r.URL.Path))
		if page < len(pages)-1 {
			last := pages[page][len(pages[page])-1]
			w.Header().Add("Link", fmt.Sprintf(`<%s%s?after=%s&limit=200>; rel="next"`, server.URL, r.URL.Path, last))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(groups)
	}))

	return server
}

func groupIDs(prefix string, from, to int) []string {
	ids := make([]string, 0, to-from+1)
	for i := from; i <= to; i++ {
		ids = append(ids, fmt.Sprintf("%s%04d", prefix, i))
	}
	return ids
}

// TestCheckIfUserHasGroups_Pagination reproduces a user whose group list spans
// three pages (over 400 groups). Before the fix, decoding page three reused the
// slice that page two had been decoded into, overwriting the *sdk.Group values
// already collected from page two, so groups near the start of page two were
// never seen by the membership check.
func TestCheckIfUserHasGroups_Pagination(t *testing.T) {
	pages := [][]string{
		groupIDs("00g", 1, 200),
		groupIDs("00g", 201, 400),
		groupIDs("00g", 401, 473),
	}

	server := pagedUserGroupsServer(t, pages)
	defer server.Close()

	ctx, client, err := sdk.NewClient(
		context.Background(),
		sdk.WithOrgUrl(server.URL),
		sdk.WithToken("00fake-token"),
		sdk.WithTestingDisableHttpsCheck(true),
		sdk.WithCache(false),
		sdk.WithRateLimitMaxRetries(0),
	)
	if err != nil {
		t.Fatalf("unable to create client: %v", err)
	}

	tests := []struct {
		name   string
		groups []string
		want   bool
	}{
		{name: "first page", groups: []string{"00g0001"}, want: true},
		{name: "start of second page", groups: []string{"00g0201"}, want: true},
		{name: "end of second page", groups: []string{"00g0400"}, want: true},
		{name: "last page", groups: []string{"00g0473"}, want: true},
		{name: "one group per page", groups: []string{"00g0100", "00g0210", "00g0450"}, want: true},
		{name: "missing group", groups: []string{"00g9999"}, want: false},
		{name: "one missing among present", groups: []string{"00g0210", "00g9999"}, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := checkIfUserHasGroups(ctx, client, "00u0001", tc.groups)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("checkIfUserHasGroups(%v) = %v, want %v", tc.groups, got, tc.want)
			}
		})
	}
}
