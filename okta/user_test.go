package okta

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestHttpClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestUserSetAllGroups(t *testing.T) {
	s := dataSourceUser().Schema
	d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
		"id": "foo",
	})
	ctx := context.Background()

	firstPageOfGroups, err := json.Marshal([]*okta.Group{
		{Id: "foo"},
		{Id: "bar"},
	})

	if err != nil {
		t.Fatal(err)
	}

	secondPageOfGroups, err := json.Marshal([]*okta.Group{
		{Id: "baz"},
		{Id: "qux"},
	})

	if err != nil {
		t.Fatal(err)
	}

	h := NewTestHttpClient(func(req *http.Request) *http.Response {
		q := req.URL.Query()

		if q.Has("after") {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(secondPageOfGroups)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(firstPageOfGroups)),
			Header: http.Header{
				"Link": []string{"<https://foo.okta.com?limit=2&after=0>; rel=\"next\""},
			},
		}
	})

	oktaCtx, c, err := okta.NewClient(
		ctx,
		okta.WithOrgUrl("https://foo.okta.com"),
		okta.WithToken("f0oT0k3n"),
		okta.WithHttpClientPtr(h),
	)

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = setAllGroups(oktaCtx, d, c)

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	groups := convertInterfaceToStringSetNullable(d.Get("group_memberships"))

	if len(groups) != 4 {
		t.Fatalf("expected 4 groups; got %d", len(groups))
	}

}
