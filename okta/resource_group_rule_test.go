package okta

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func sweepGroupRules(client *testClient) error {
	var errorList []error
	// Should never need to deal with pagination
	rules, _, err := client.oktaClient.Group.ListRules(&query.Params{Limit: 300})
	if err != nil {
		return err
	}

	for _, s := range rules {
		if _, err := client.oktaClient.Group.DeactivateRule(s.Id); err != nil {
			errorList = append(errorList, err)
			continue
		}
		if _, err := client.oktaClient.Group.DeleteRule(s.Id, nil); err != nil {
			errorList = append(errorList, err)
		}

	}
	return condenseError(errorList)
}

func TestAccOktaGroupRuleCrud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := buildResourceFQN(groupRule, ri)
	mgr := newFixtureManager("okta_group")
	config := mgr.GetFixtures("okta_group.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_group_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testgroupdifferent"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testgroupdifferent"),
				),
			},
		},
	})
}
