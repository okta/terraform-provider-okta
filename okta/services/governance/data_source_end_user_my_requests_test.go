package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestDataSourceMyRequests(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceEndUserMyRequests, t.Name())
	config := mgr.GetFixtures("data-source.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_end_user_my_requests.example", "entry_id", "cen123456789abcdefgh"),
					resource.TestCheckResourceAttr("data.okta_end_user_my_requests.example", "id", "req123abcd456ghijklm"),
					resource.TestCheckResourceAttr("data.okta_end_user_my_requests.example", "status", "PENDING"),
				),
			},
		},
	})
}
