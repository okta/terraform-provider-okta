package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaRequestConditions_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceRequestConditions, t.Name())
	config := mgr.GetFixtures("datasource_list.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_request_conditions.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_request_conditions.test", "resource_id"),
					resource.TestCheckResourceAttrSet("data.okta_request_conditions.test", "conditions.#"),
				),
			},
		},
	})
}
