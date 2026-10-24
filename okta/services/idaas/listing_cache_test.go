package idaas

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/sdk"
)

// fakeListing serves a fixed set of IDs and counts how often it is listed.
type fakeListing struct {
	ids   []string
	err   error
	lists atomic.Int32
}

func (f *fakeListing) list(context.Context) (map[string]string, error) {
	f.lists.Add(1)
	if f.err != nil {
		return nil, f.err
	}
	byID := map[string]string{}
	for _, id := range f.ids {
		byID[id] = "listed " + id
	}
	return byID, nil
}

// expectLookup fails the test unless looking up id gives ok, and, when ok,
// the listed object.
func expectLookup(t *testing.T, c *listingCache[string], f *fakeListing, scope, id string, ok bool) {
	t.Helper()
	obj, gotOK, err := c.lookup(context.Background(), scope, id, f.list)
	if err != nil {
		t.Fatalf("lookup(%q, %q) returned error %v", scope, id, err)
	}
	if gotOK != ok || (ok && obj != "listed "+id) {
		t.Fatalf("lookup(%q, %q) = %q, %v; want ok %v", scope, id, obj, gotOK, ok)
	}
}

func expectLists(t *testing.T, f *fakeListing, want int32) {
	t.Helper()
	if got := f.lists.Load(); got != want {
		t.Errorf("listed %d times, want %d", got, want)
	}
}

func TestListingCacheOneIDNeverLists(t *testing.T) {
	f := &fakeListing{ids: []string{"a"}}
	c := newListingCache[string](2)
	expectLookup(t, c, f, "s", "a", false)
	expectLookup(t, c, f, "s", "a", false)
	expectLists(t, f, 0)
}

func TestListingCacheSecondIDListsOnce(t *testing.T) {
	f := &fakeListing{ids: []string{"a", "b", "c"}}
	c := newListingCache[string](2)
	expectLookup(t, c, f, "s", "a", false)
	expectLookup(t, c, f, "s", "b", true)
	expectLookup(t, c, f, "s", "c", true)
	expectLookup(t, c, f, "s", "a", true)
	expectLists(t, f, 1)
}

func TestListingCacheListsAtListAfter(t *testing.T) {
	f := &fakeListing{ids: []string{"a", "b", "c", "d"}}
	c := newListingCache[string](3)
	expectLookup(t, c, f, "s", "a", false)
	expectLookup(t, c, f, "s", "b", false)
	expectLookup(t, c, f, "s", "a", false)
	expectLists(t, f, 0)
	expectLookup(t, c, f, "s", "c", true)
	expectLookup(t, c, f, "s", "d", true)
	expectLookup(t, c, f, "s", "a", true)
	expectLists(t, f, 1)
}

func TestListingCacheScopesAreIndependent(t *testing.T) {
	f := &fakeListing{ids: []string{"a", "b"}}
	c := newListingCache[string](2)
	expectLookup(t, c, f, "s1", "a", false)
	expectLookup(t, c, f, "s2", "b", false)
	expectLists(t, f, 0)
}

func TestListingCacheMissReadsDirectly(t *testing.T) {
	f := &fakeListing{ids: []string{"a", "b"}}
	c := newListingCache[string](2)
	expectLookup(t, c, f, "s", "a", false)
	expectLookup(t, c, f, "s", "b", true)
	expectLookup(t, c, f, "s", "created-after-listing", false)
	expectLists(t, f, 1)
}

func TestListingCacheFailedListingReadsDirectly(t *testing.T) {
	f := &fakeListing{ids: []string{"a", "b", "c"}, err: errors.New("429")}
	c := newListingCache[string](2)
	expectLookup(t, c, f, "s", "a", false)
	expectLookup(t, c, f, "s", "b", false)
	expectLookup(t, c, f, "s", "c", false)
	expectLists(t, f, 1)
}

func TestListingCacheForgottenIDReadsDirectly(t *testing.T) {
	f := &fakeListing{ids: []string{"a", "b"}}
	c := newListingCache[string](2)
	expectLookup(t, c, f, "s", "a", false)
	expectLookup(t, c, f, "s", "b", true)
	c.forget("s", "b")
	expectLookup(t, c, f, "s", "b", false)
	expectLookup(t, c, f, "s", "a", true)
	expectLists(t, f, 1)
}

func TestListingCacheForgetDuringListing(t *testing.T) {
	listing, release := make(chan struct{}), make(chan struct{})
	slow := func(ctx context.Context) (map[string]string, error) {
		close(listing)
		<-release
		return map[string]string{"a": "stale a", "b": "listed b"}, nil
	}
	c := newListingCache[string](2)
	if _, ok, _ := c.lookup(context.Background(), "s", "a", slow); ok {
		t.Fatal("first lookup was served from a listing")
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _, _ = c.lookup(context.Background(), "s", "b", slow)
	}()
	<-listing
	c.forget("s", "a")
	close(release)
	<-done
	if obj, ok, _ := c.lookup(context.Background(), "s", "a", slow); ok {
		t.Fatalf("an ID written during the listing was served from it: %q", obj)
	}
}

func TestListingCacheConcurrentLookupsShareOneListing(t *testing.T) {
	const n = 200
	ids := make([]string, n)
	for i := range ids {
		ids[i] = fmt.Sprintf("id%d", i)
	}
	f := &fakeListing{ids: ids}
	c := newListingCache[string](2)

	var wg sync.WaitGroup
	var direct atomic.Int32
	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			obj, ok, err := c.lookup(context.Background(), "s", id, f.list)
			switch {
			case err != nil:
				t.Errorf("lookup(%q) returned error %v", id, err)
			case !ok:
				direct.Add(1)
			case obj != "listed "+id:
				t.Errorf("lookup(%q) = %q", id, obj)
			}
		}(id)
	}
	wg.Wait()
	expectLists(t, f, 1)
	if got := direct.Load(); got != 1 {
		t.Errorf("%d lookups read directly, want only the first", got)
	}
}

func TestListingCacheWaiterHonoursContext(t *testing.T) {
	listing, release := make(chan struct{}), make(chan struct{})
	slow := func(context.Context) (map[string]string, error) {
		close(listing)
		<-release
		return nil, nil
	}
	c := newListingCache[string](2)
	_, _, _ = c.lookup(context.Background(), "s", "a", slow)
	go func() { _, _, _ = c.lookup(context.Background(), "s", "b", slow) }()
	<-listing

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := c.lookup(ctx, "s", "c", slow); !errors.Is(err, context.Canceled) {
		t.Fatalf("lookup with a cancelled context returned %v, want context.Canceled", err)
	}
	close(release)
}

func TestListingCachesArePerConfig(t *testing.T) {
	a, b := &config.Config{}, &config.Config{}
	if listingCachesFor(a) != listingCachesFor(a) {
		t.Error("one config returned two caches")
	}
	if listingCachesFor(a) == listingCachesFor(b) {
		t.Error("two configs shared a cache")
	}
}

// seedGroupListing makes groups the completed listing of meta's org scope,
// taken at listedAt.
func seedGroupListing(meta *config.Config, listedAt time.Time, groups ...*sdk.Group) {
	listed := make(chan struct{})
	close(listed)
	byID := map[string]*listedGroup{}
	for _, g := range groups {
		byID[g.Id] = &listedGroup{Group: g, listedAt: listedAt}
	}
	c := listingCachesFor(meta).groups
	c.scopes[""] = &listingScope[*listedGroup]{seen: map[string]struct{}{}, listed: listed, byID: byID, dirty: map[string]struct{}{}}
}

func TestRecordableMembershipUpdatedWaitsOutItsSecond(t *testing.T) {
	changed := time.Date(2026, 10, 23, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name     string
		listedAt time.Time
		want     string
	}{
		{"listed in the same second", changed.Add(999 * time.Millisecond), ""},
		{"listed a second later", changed.Add(time.Second), "2026-10-23T12:00:00Z"},
		{"listing time unknown", time.Time{}, ""},
	} {
		g := &listedGroup{Group: &sdk.Group{LastMembershipUpdated: &changed}, listedAt: tc.listedAt}
		if got := recordableMembershipUpdated(g); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
	if got := recordableMembershipUpdated(&listedGroup{Group: &sdk.Group{}, listedAt: changed}); got != "" {
		t.Errorf("no lastMembershipUpdated: got %q, want empty", got)
	}
}

func TestGroupMembershipsReadSkipsUnchangedMembership(t *testing.T) {
	changed := time.Date(2026, 10, 23, 12, 0, 0, 0, time.UTC)
	// No Okta client: reading the members would panic.
	meta := &config.Config{}
	seedGroupListing(meta, changed.Add(time.Hour), &sdk.Group{Id: "00g1", LastMembershipUpdated: &changed})

	d := schema.TestResourceDataRaw(t, resourceGroupMemberships().Schema, map[string]interface{}{
		"group_id":        "00g1",
		"users":           []interface{}{"00u1"},
		"track_all_users": true,
	})
	if err := d.Set("last_membership_updated", changed.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	if diags := resourceGroupMembershipsRead(context.Background(), d, meta); diags.HasError() {
		t.Fatalf("read returned %v", diags)
	}
	if got := d.Get("users").(*schema.Set).List(); len(got) != 1 || got[0] != "00u1" {
		t.Errorf("users = %v, want the stored [00u1]", got)
	}
}
