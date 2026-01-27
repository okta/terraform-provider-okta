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
				ImportStateId:      "0oaoum6j3cElINe1z1d7/68cbc2b263c689fc3336bfac",
				ImportStatePersist: true,
				Config:             config,
				PlanOnly:           true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Sequence"),
					resource.TestCheckResourceAttr(resourceName, "description", "Sequence for testing TF"),
				),
			},
		},
	})
}
