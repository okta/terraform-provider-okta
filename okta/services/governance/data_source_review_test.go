package governance_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"testing"
)

func TestAccDataSourceOktaReview_read(t *testing.T) {
	mgr := newFixtureManager("Data-sources", resources.GovernanceReview, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("Data.okta_review.test", "id"),
					resource.TestCheckResourceAttr("Data.okta_review.test", "campaign_id", "icizigd86iM9sOcbN1d6"),
					resource.TestCheckResourceAttr("Data.okta_review.test", "decision", "UNREVIEWED"),
				),
			},
		},
	})
}
