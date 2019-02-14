package okta

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func sweepGroups(client *testClient) error {
	var errorList []error
	// Should never need to deal with pagination, limit is 10,000 by default
	groups, _, err := client.oktaClient.Group.ListGroups(&query.Params{Q: testResourcePrefix})
	if err != nil {
		return err
	}

	for _, s := range groups {
		if _, err := client.oktaClient.Group.DeleteGroup(s.Id); err != nil {
			errorList = append(errorList, err)
		}

	}
	return condenseError(errorList)
}

func TestAccOktaGroupsCreate(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := buildResourceFQN(oktaGroup, ri)
	mgr := newFixtureManager("okta_group")
	config := mgr.GetFixtures("okta_group.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_group_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testgroupdifferent")),
			},
		},
	})
}
