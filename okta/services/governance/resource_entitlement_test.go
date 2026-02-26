package governance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccEntitlement_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceEntitlement, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceEntitlement)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Entitlement Bundle"),
					resource.TestCheckResourceAttr(resourceName, "external_value", "Entitlement Bundle"),
					resource.TestCheckResourceAttr(resourceName, "parent.type", "APPLICATION"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "values.0.external_value", "value_1"),
				),
			},
		},
	})
}
