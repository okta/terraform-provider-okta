package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaAppSignOnPolicy_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSignOnPolicy)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	renamedConfig := mgr.GetFixtures("basic_renamed.tf", ri, t)
	resourceName := fmt.Sprintf("%v.test", appSignOnPolicy)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createPolicyCheckDestroy(appSignOnPolicy),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("Test_App", ri)),
					resource.TestCheckResourceAttr(resourceName, "description", "The app signon policy used by our test app."),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("Test_App", ri)),
					resource.TestCheckResourceAttr(resourceName, "description", "The updated app signon policy used by our test app."),
				),
			},
			{
				Config: renamedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensurePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceNameWithPrefix("Test_App_Renamed", ri)),
					resource.TestCheckResourceAttr(resourceName, "description", "The app signon policy used by our test app."),
				),
			},
		},
	})
}
