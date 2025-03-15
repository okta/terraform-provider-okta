package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAuthServerPolicy_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAuthServerPolicy, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	createServerWithPolicy := buildTestAuthServerWithPolicy(mgr.Seed)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
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
