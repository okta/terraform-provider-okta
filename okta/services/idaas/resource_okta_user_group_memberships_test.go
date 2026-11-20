package idaas_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jarcoal/httpmock"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/stretchr/testify/require"
)

func TestAccResourceOktaUserGroupMemberships_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUserGroupMemberships, t.Name())
	start := mgr.GetFixtures("basic.tf", t)
	update := mgr.GetFixtures("basic_update.tf", t)
	remove := mgr.GetFixtures("basic_removal.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: start,
			},
			{
				Config: update,
			},
			{
				Config: remove,
			},
		},
	})
}

// TestUserGroupMembershipsReadPaginationDoesNotDropGroups is a regression test
// for a pointer aliasing bug in checkIfUserHasGroups.
//
// The slice each page was decoded into used to be declared once, outside the
// pagination loop. encoding/json decodes into a slice that already holds
// non-nil pointers by reusing those pointers and overwriting the structs in
// place, so decoding a page silently rewrote groups already appended to the
// accumulator from an earlier page. Groups the user really is a member of
// disappeared from the collected set, and the resource then reported the
// membership as missing.
//
// The corruption only begins on the THIRD page: page two decodes into a nil
// slice, which allocates fresh. So reproducing it needs three pages, with the
// expected group positioned inside the prefix that page three overwrites
// (index < len(page three)). Two pages, or a group late in page two, both pass
// even with the bug present.
func TestUserGroupMembershipsReadPaginationDoesNotDropGroups(t *testing.T) {
	t.Setenv("OKTA_ORG_NAME", "test")
	t.Setenv("OKTA_BASE_URL", "example.com")
	t.Setenv("OKTA_API_TOKEN", "token")

	const (
		userID  = "00uTESTUSER00000000"
		groupID = "00gEXPECTED00000000"
		cursor1 = "00gCURSOR1000000000"
		cursor2 = "00gCURSOR2000000000"
	)

	groupsURL := fmt.Sprintf("https://test.example.com/api/v1/users/%s/groups", userID)

	groupsBody := func(ids ...string) string {
		out := "["
		for i, id := range ids {
			if i > 0 {
				out += ","
			}
			out += fmt.Sprintf(`{"id":%q,"type":"OKTA_GROUP","profile":{"name":%q}}`, id, "group-"+id)
		}
		return out + "]"
	}

	// groupID sits at index 0 of page two, and page three is non-empty, so the
	// aliasing bug overwrites groupID with page three's first group.
	pages := map[string]struct {
		body string
		next string
	}{
		"":      {body: groupsBody("00gPAGE1A0000000000", "00gPAGE1B0000000000"), next: cursor1},
		cursor1: {body: groupsBody(groupID, "00gPAGE2B0000000000"), next: cursor2},
		cursor2: {body: groupsBody("00gPAGE3A0000000000"), next: ""},
	}

	d := schema.TestResourceDataRaw(t, idaas.ProviderResources()[resources.OktaIDaaSUserGroupMemberships].Schema, map[string]interface{}{
		"user_id": userID,
		"groups":  []interface{}{groupID},
	})
	d.SetId(userID)

	c := config.NewConfig(d)
	require.NoError(t, c.LoadAPIClient())

	defer httpmock.DeactivateAndReset()
	httpmock.ActivateNonDefault(c.OktaIDaaSClient.OktaSDKClientV2().GetConfig().HttpClient)
	httpmock.RegisterResponder("GET", groupsURL, func(req *http.Request) (*http.Response, error) {
		after := req.URL.Query().Get("after")
		page, ok := pages[after]
		if !ok {
			return nil, fmt.Errorf("unexpected after cursor %q", after)
		}
		resp := httpmock.NewStringResponse(http.StatusOK, page.body)
		resp.Header.Set("Content-Type", "application/json")
		resp.Request = req
		resp.Header.Add("Link", fmt.Sprintf(`<%s?limit=200>; rel="self"`, groupsURL))
		if page.next != "" {
			resp.Header.Add("Link", fmt.Sprintf(`<%s?after=%s&limit=200>; rel="next"`, groupsURL, page.next))
		}
		return resp, nil
	})

	diags := idaas.ProviderResources()[resources.OktaIDaaSUserGroupMemberships].ReadContext(context.Background(), d, c)
	require.False(t, diags.HasError(), "read returned errors: %+v", diags)

	// Guard the test itself: all three pages must have been walked, otherwise
	// the scenario no longer covers the aliasing bug.
	require.Equal(t, 3, httpmock.GetTotalCallCount(), "expected all three pages to be fetched")

	// Read blanks the ID when it believes the user is missing an expected group.
	require.Equal(t, userID, d.Id(),
		"expected group %s was reported missing; it is returned on page two of the user's groups, so a page-three decode dropped it", groupID)
}
