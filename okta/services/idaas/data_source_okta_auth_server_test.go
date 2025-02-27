package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccDataSourceOktaAuthServer_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAuthServer, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	authServer := buildTestAuthServer(mgr.Seed)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: authServer,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_auth_server.test", "id"),
					resource.TestCheckResourceAttr("data.okta_auth_server.test", "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr("data.okta_auth_server.test", "status", idaas.StatusActive),
					resource.TestCheckResourceAttrSet("data.okta_auth_server.test", "issuer"),
				),
			},
		},
	})
}

func buildTestAuthServer(i int) string {
	return fmt.Sprintf(`
resource "okta_auth_server" "test" {
  audiences   = ["whatever.rise.zone"]
  description = "test"
  name        = "testAcc_%d"
}`, i)
}
