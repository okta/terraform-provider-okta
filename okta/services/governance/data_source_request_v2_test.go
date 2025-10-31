package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaRequestV2_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceRequestV2, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_request_v2.test", "requested.access_scope_type", "APPLICATION"),
					resource.TestCheckResourceAttr("data.okta_request_v2.test", "requested.type", "CATALOG_ENTRY"),
					resource.TestCheckResourceAttr("data.okta_request_v2.test", "requested_by.type", "OKTA_USER"),
				),
			},
		},
	})
}
