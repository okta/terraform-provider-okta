package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaIdentitySourceSessions_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSIdentitySourceSessions, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_identity_source_sessions.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_identity_source_sessions.test", "status"),
					resource.TestCheckResourceAttrSet("data.okta_identity_source_sessions.test", "import_type"),
				),
			},
		},
	})
}
