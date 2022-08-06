package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaAppSignOnPolicy_crud(t *testing.T) {
	mgr := newFixtureManager(appSignOnPolicy, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	renamedConfig := mgr.GetFixtures("basic_renamed.tf", t)
	resourceName := fmt.Sprintf("%v.test", appSignOnPolicy)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createPolicyCheckDestroy(appSignOnPolicy),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("Test_App", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "The app signon policy used by our test app."),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("Test_App", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "The updated app signon policy used by our test app."),
				),
			},
			{
				Config: renamedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("Test_App_Renamed", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "The app signon policy used by our test app."),
				),
			},
		},
	})

}
