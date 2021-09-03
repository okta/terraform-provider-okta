package okta

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func getGroupPagesJson(t *testing.T) ([]byte, []byte) {
	firstPageOfGroups, err := json.Marshal([]*okta.Group{
		{Id: "foo"},
		{Id: "bar"},
	})
	if err != nil {
		t.Fatalf("could not serialize first set of groups: %s", err)
	}

	secondPageOfGroups, err := json.Marshal([]*okta.Group{
		{Id: "baz"},
		{Id: "qux"},
	})
	if err != nil {
		t.Fatalf("could not serialize second set of groups: %s", err)
	}

	return firstPageOfGroups, secondPageOfGroups
}

type userGroupFunc func(ctx context.Context, d *schema.ResourceData, c *okta.Client) error

func testUserGroupFetchesAllPages(t *testing.T, fn userGroupFunc) {
	s := dataSourceUser().Schema
	d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
		"Id": "foo",
	})

	firstPageOfGroups, secondPageOfGroups := getGroupPagesJson(t)

	ctx, c, err := newTestOktaClientWithResponse(func(req *http.Request) *http.Response {
		q := req.URL.Query()

		if q.Has("after") {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(secondPageOfGroups)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(firstPageOfGroups)),
			Header: http.Header{
				"Link":         []string{"<https://foo.okta.com?limit=2&after=0>; rel=\"next\""},
				"Content-Type": []string{"application/json"},
			},
		}
	})
	if err != nil {
		t.Fatalf("could not create an okta client instance: %s", err)
	}

	err = fn(ctx, d, c)

	if err != nil {
		t.Fatalf("fetching groups failed: %s", err)
	}

	groups := convertInterfaceToStringSetNullable(d.Get("group_memberships"))

	if len(groups) != 4 {
		t.Fatalf("expected 4 groups; got %d", len(groups))
	}

	expected := []string{"bar", "baz", "foo", "qux"}
	sort.Strings(groups)
	allMatch := true

	for i, group := range expected {
		if group != groups[i] {
			allMatch = false
		}
	}

	if !allMatch {
		t.Fatalf("expected %s; got %s", expected, groups)
	}
}

func TestUserSetAllGroups(t *testing.T) {
	testUserGroupFetchesAllPages(t, setAllGroups)
}

func TestUserSetGroups(t *testing.T) {
	testUserGroupFetchesAllPages(t, setGroups)
}
