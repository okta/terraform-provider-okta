package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaPrincipalAccess_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.GovernancePrincipalAccess, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_principal_access.test", "target_principal.type", "OKTA_USER"),
					resource.TestCheckResourceAttr("data.okta_principal_access.test", "parent.type", "APPLICATION"),
					resource.TestCheckResourceAttr("data.okta_principal_access.test", "additional.0.grant_type", "ENTITLEMENT-BUNDLE"),
				),
			},
		},
	})
}
