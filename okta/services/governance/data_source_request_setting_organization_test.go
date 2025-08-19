package governance_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"testing"
)

func TestAccDataSourceOktaRequestSettingOrganization_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.GovernanceRequestSettingOrganization, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_request_setting_organization.test", "provisioning_status", "PROVISIONED"),
					resource.TestCheckResourceAttr("data.okta_request_setting_organization.test", "subprocessors_acknowledged", "true"),
				),
			},
		},
	})
}
