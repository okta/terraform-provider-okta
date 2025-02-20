package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

// Test failed due to 500 error. The actual script run normally
func TestAccResourceOktaTrustedServers_crud(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_trusted_server", t.Name())

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.AccMergeProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("basic.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_trusted_server.example", "trusted.#", "2"),
				),
			},
			{
				Config: mgr.GetFixtures("update.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_trusted_server.example", "trusted.#", "2"),
				),
			},
			{
				Config: mgr.GetFixtures("update_2.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_trusted_server.example", "trusted.#", "2"),
				),
			},
		},
	})
}
