package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaEndUserMyCatalogsEntry_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceEndUsersMyCatalogsEntry, t.Name())
	config := mgr.GetFixtures("basic.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_end_user_my_catalogs_entry.test", "id", "cen1043mnuxScMKl91d7"),
					resource.TestCheckResourceAttr("data.okta_end_user_my_catalogs_entry.test", "name", "Workplace by Facebook"),
					resource.TestCheckResourceAttr("data.okta_end_user_my_catalogs_entry.test", "requestable", "false"),
					resource.TestCheckResourceAttr("data.okta_end_user_my_catalogs_entry.test", "label", "Application"),
				),
			},
		},
	})
}
