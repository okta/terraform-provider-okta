package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaCampaign_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceCampaign, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_campaign.test", "id"),
					resource.TestCheckResourceAttr("data.okta_campaign.test", "name", "Monthly access review of sales team"),
					resource.TestCheckResourceAttr("data.okta_campaign.test", "resource_settings.type", "GROUP"),
					resource.TestCheckResourceAttr("data.okta_campaign.test", "principal_scope_settings.type", "USERS"),
				),
			},
		},
	})
}
