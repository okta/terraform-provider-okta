package okta4ai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaApiServer_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaApiServer, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaApiServer)
	resourceURL := fmt.Sprintf("https://testacc-api-server-%d.example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "orn"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "resource_url", resourceURL),
					resource.TestCheckResourceAttr(resourceName, "metadata.display_name", fmt.Sprintf("testAcc-api-server-%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "metadata.description", "Terraform acceptance test API server"),
					resource.TestCheckResourceAttr(resourceName, "metadata.scopes_supported.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.scopes_supported.0", "read:items"),
				),
			},
			{
				// Changing metadata updates the API server in place.
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "resource_url", resourceURL),
					resource.TestCheckResourceAttr(resourceName, "metadata.display_name", fmt.Sprintf("testAcc-api-server-updated-%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "metadata.description", "Terraform acceptance test API server, updated"),
					resource.TestCheckResourceAttr(resourceName, "metadata.scopes_supported.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.scopes_supported.1", "write:items"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
