package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAuthServerPolicy_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(authServerPolicy)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	createServerWithPolicy := buildTestAuthServerWithPolicy(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
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
