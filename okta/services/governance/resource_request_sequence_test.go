package governance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccRequestSequenceResource_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceRequestSequence, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceRequestSequence)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy: func(state *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				ImportState:        true,
				ResourceName:       "okta_request_sequence.test",
				ImportStateId:      "0oa15zse2by0vGzyP1d8/6af8f21dc78cc307be8e6b8b",
				ImportStatePersist: true,
				Config:             config,
				PlanOnly:           true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf testing seq"),
					resource.TestCheckResourceAttr(resourceName, "description", "Sequence for testing TF"),
				),
			},
		},
	})
}
