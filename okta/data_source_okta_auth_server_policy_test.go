package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaAuthServerPolicy_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", authServerPolicy, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	createServerWithPolicy := buildTestAuthServerWithPolicy(mgr.Seed)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: createServerWithPolicy,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_auth_server_policy.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_auth_server_policy.test", "description"),
				),
			},
		},
	})
}

func buildTestAuthServerWithPolicy(i int) string {
	return fmt.Sprintf(`
resource "okta_auth_server_policy" "test" {
  status           = "ACTIVE"
  name             = "test"
  description      = "test"
  priority         = 1
  client_whitelist = [
    "ALL_CLIENTS"]
  auth_server_id   = okta_auth_server.test.id
}

resource "okta_auth_server" "test" {
  name        = "testAcc_%d"
  description = "test"
  audiences   = [
    "whatever.rise.zone"]
}`, i)
}
