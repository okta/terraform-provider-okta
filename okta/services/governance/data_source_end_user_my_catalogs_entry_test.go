package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaEndUserMyCatalogsEntry_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceCampaign, t.Name())
	config := mgr.GetFixtures("basic.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_end_user_my_requests_entry.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_end_user_my_requests_entry.test", "name"),
					resource.TestCheckResourceAttrSet("data.okta_end_user_my_requests_entry.test", "requestable"),
					resource.TestCheckResourceAttrSet("data.okta_end_user_my_requests_entry.test", "label"),
					// Parent might be empty for root entries
					resource.TestCheckResourceAttr("data.okta_end_user_my_requests_entry.test", "counts.resource_counts.applications", "1"),
				),
			},
		},
	})
}
