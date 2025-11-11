package governance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccRequestV2Resource_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceRequestV2, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceRequestV2)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "SUBMITTED"),
					resource.TestCheckResourceAttr(resourceName, "requested.access_scope_type", "APPLICATION"),
					resource.TestCheckResourceAttr(resourceName, "requested.resource_type", "APPLICATION"),
					resource.TestCheckResourceAttr(resourceName, "requested.type", "CATALOG_ENTRY"),
					resource.TestCheckResourceAttr(resourceName, "requested_for.type", "OKTA_USER"),
				),
			},
		},
	})
}
