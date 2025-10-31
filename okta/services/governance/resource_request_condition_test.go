package governance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccRequestConditionResource_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceRequestCondition, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceRequestCondition)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-condition"),
					resource.TestCheckResourceAttr(resourceName, "requester_settings.type", "EVERYONE"),
				),
			},
			{
				Config: mgr.ConfigReplace(updatedConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-condition"),
					resource.TestCheckResourceAttr(resourceName, "requester_settings.type", "GROUPS"),
				),
			},
		},
	})
}

func TestAccRequestConditionResource_Issue2510(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceRequestCondition, t.Name())
	config := mgr.GetFixtures("basic_issue2510.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceRequestCondition)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "issue-2510"),
					resource.TestCheckResourceAttr(resourceName, "requester_settings.type", "GROUPS"),
				),
			},
		},
	})
}
