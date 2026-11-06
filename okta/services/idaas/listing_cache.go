package idaas

import (
	"context"
	"net/http"
	"sync"
	"time"

	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

// listingCache lets the refreshes of many objects that one API call can list
// share that call. Reading each object with its own GET spends one request of
// a per-minute rate limit apiece, so a configuration with hundreds of groups or
// push mappings stalls every refresh on those limits.
type listingCache[T any] struct {
	// listAfter is how many distinct IDs of a scope are read directly before
	// the scope is listed instead.
	listAfter int

	mu     sync.Mutex
	scopes map[string]*listingScope[T]
}

type listingScope[T any] struct {
	// seen holds the IDs requested before a listing started.
	seen map[string]struct{}
	// listed is non-nil once a listing has started and is closed when it ends.
	listed chan struct{}
	// byID is nil until a listing succeeds.
	byID map[string]T
	// dirty holds the IDs written since the listing may have been taken.
	dirty map[string]struct{}
}

type listFunc[T any] func(ctx context.Context) (map[string]T, error)

func newListingCache[T any](listAfter int) *listingCache[T] {
	return &listingCache[T]{listAfter: listAfter, scopes: map[string]*listingScope[T]{}}
}

func (c *listingCache[T]) scope(name string) *listingScope[T] {
	s := c.scopes[name]
	if s == nil {
		s = &listingScope[T]{seen: map[string]struct{}{}, dirty: map[string]struct{}{}}
		c.scopes[name] = s
	}
	return s
}

// lookup returns the object with the given ID from one listing of its scope.
// The scope is listed once listAfter distinct IDs in it have been requested,
// so a configuration with fewer objects per scope makes no extra request. ok
// is false when the caller should read the object itself: before the listing,
// after a failed one, for an ID it lacks, and for an ID written since.
func (c *listingCache[T]) lookup(ctx context.Context, scope, id string, list listFunc[T]) (obj T, ok bool, err error) {
	c.mu.Lock()
	s := c.scope(scope)
	start := false
	if s.listed == nil {
		s.seen[id] = struct{}{}
		if len(s.seen) >= c.listAfter {
			s.listed = make(chan struct{})
			start = true
		}
	}
	listed := s.listed
	c.mu.Unlock()

	if start {
		byID, err := list(ctx)
		if err == nil {
			c.mu.Lock()
			s.byID = byID
			c.mu.Unlock()
		}
		close(listed)
	}

	if listed == nil {
		return obj, false, nil
	}
	select {
	case <-listed:
	case <-ctx.Done():
		return obj, false, ctx.Err()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, written := s.dirty[id]; written {
		return obj, false, nil
	}
	obj, ok = s.byID[id]
	return obj, ok, nil
}

// forget makes later lookups of an ID the caller has just written read it
// directly, including when a listing still in flight finishes afterwards.
func (c *listingCache[T]) forget(scope, id string) {
	c.mu.Lock()
	c.scope(scope).dirty[id] = struct{}{}
	c.mu.Unlock()
}

// providerListingCaches holds one cache per object type for a configured
// provider, so aliased providers never share listings.
type providerListingCaches struct {
	groupPushMappings *listingCache[v6okta.GroupPushMapping] // scoped by app ID
	groups            *listingCache[*listedGroup]            // one scope: the org
}

// listedGroup is a group as the org listing returned it, with the server time
// the listing began.
type listedGroup struct {
	*sdk.Group
	listedAt time.Time
}

// recordableMembershipUpdated returns the group's lastMembershipUpdated to
// keep in state, or "" while a later change could still share its second:
// Okta truncates it to whole seconds, so a change in the same second as the
// last one would not move it.
func recordableMembershipUpdated(g *listedGroup) string {
	if g.LastMembershipUpdated == nil {
		return ""
	}
	changed := g.LastMembershipUpdated.UTC()
	if g.listedAt.Before(changed.Add(time.Second)) {
		return ""
	}
	return changed.Format(time.RFC3339Nano)
}

var listingCaches sync.Map // map[*config.Config]*providerListingCaches

func listingCachesFor(meta interface{}) *providerListingCaches {
	c, _ := meta.(*config.Config)
	if v, ok := listingCaches.Load(c); ok {
		return v.(*providerListingCaches)
	}
	v, _ := listingCaches.LoadOrStore(c, &providerListingCaches{
		// An app's mappings are a small listing, worth it from the second.
		groupPushMappings: newListingCache[v6okta.GroupPushMapping](2),
		// Every OKTA_GROUP in the org is a large one, worth it only for a
		// configuration of many groups.
		groups: newListingCache[*listedGroup](10),
	})
	return v.(*providerListingCaches)
}

// listGroupPushMappings lists every push mapping of an app at the API's page
// maximum.
func listGroupPushMappings(c *config.Config, appID string) listFunc[v6okta.GroupPushMapping] {
	return func(ctx context.Context) (map[string]v6okta.GroupPushMapping, error) {
		byID := map[string]v6okta.GroupPushMapping{}
		var cursor string
		for {
			mappings, resp, err := c.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.ListGroupPushMappings(ctx, appID).After(cursor).Limit(1000).Execute()
			if err != nil {
				return nil, err
			}
			for _, m := range mappings {
				byID[m.GetId()] = m
			}
			if !resp.HasNextPage() || len(mappings) == 0 {
				return byID, nil
			}
			cursor = mappings[len(mappings)-1].GetId()
		}
	}
}

// listOktaGroups lists every OKTA_GROUP, the only type okta_group creates, at
// the API's page maximum. It filters rather than searches: search results come
// from an eventually consistent index, and a refresh must see a change made
// just before it.
func listOktaGroups(meta interface{}) listFunc[*listedGroup] {
	return func(ctx context.Context) (map[string]*listedGroup, error) {
		client := getOktaClientFromMetadata(meta)
		groups, resp, err := client.Group.ListGroups(ctx, &query.Params{Filter: `type eq "OKTA_GROUP"`, Limit: 10000})
		if err != nil {
			return nil, err
		}
		// A missing or unparsable Date leaves the zero time, so no membership
		// timestamp is recorded from this listing.
		listedAt, _ := http.ParseTime(resp.Header.Get("Date"))
		byID := map[string]*listedGroup{}
		for {
			for _, g := range groups {
				byID[g.Id] = &listedGroup{Group: g, listedAt: listedAt}
			}
			if !resp.HasNextPage() {
				return byID, nil
			}
			groups = nil
			if resp, err = resp.Next(ctx, &groups); err != nil {
				return nil, err
			}
		}
	}
}
