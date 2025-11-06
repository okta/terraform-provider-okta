package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaRequestSequence_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceRequestSequence, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_request_sequence.test", "id"),
					resource.TestCheckResourceAttr("data.okta_request_sequence.test", "name", "Business Justification + Requester's Manager Approval/Justification"),
					resource.TestCheckResourceAttr("data.okta_request_sequence.test", "compatible_resource_types.0", "APP"),
					resource.TestCheckResourceAttr("data.okta_request_sequence.test", "compatible_resource_types.1", "GROUP"),
				),
			},
		},
	})
}
