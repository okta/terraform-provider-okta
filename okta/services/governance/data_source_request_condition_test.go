package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaRequestCondition_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceRequestCondition, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_request_condition.test", "id"),
					resource.TestCheckResourceAttr("data.okta_request_condition.test", "name", "request-condition-test"),
					resource.TestCheckResourceAttr("data.okta_request_condition.test", "access_scope_settings.type", "RESOURCE_DEFAULT"),
					resource.TestCheckResourceAttr("data.okta_request_condition.test", "requester_settings.type", "EVERYONE"),
				),
			},
		},
	})
}
