package governance_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccRequestConditionsResource_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceRequestConditions, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceRequestConditions)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRequestConditionsDestroy,
		Steps: []resource.TestStep{
			{
				// Conditions are declared out of priority order; state keeps
				// the config order.
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "resource_id", "0oa10z0gowhJuJBAo1d8"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.name", "test-conditions-c"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.priority", "2"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.name", "test-conditions-a"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.2.name", "test-conditions-b"),
					resource.TestCheckResourceAttr(resourceName, "condition.2.priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.2.status", "INACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "condition.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "condition.1.id"),
					resource.TestCheckResourceAttrSet(resourceName, "condition.2.id"),
					checkRequestConditionsExist(resourceName),
				),
			},
			{
				// Reprioritise, change status, remove one condition and add one.
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "condition.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.name", "test-conditions-c"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.name", "test-conditions-a"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "condition.2.name", "test-conditions-d"),
					resource.TestCheckResourceAttr(resourceName, "condition.2.priority", "2"),
					checkRequestConditionsExist(resourceName),
				),
			},
			{
				Config:   updatedConfig,
				PlanOnly: true,
			},
		},
	})
}

// checkRequestConditionsExist verifies that every condition in state exists in Okta.
func checkRequestConditionsExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
			return nil
		}
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		resourceID := rs.Primary.Attributes["resource_id"]
		var count int
		if _, err := fmt.Sscan(rs.Primary.Attributes["condition.#"], &count); err != nil {
			return err
		}
		for i := 0; i < count; i++ {
			id := rs.Primary.Attributes[fmt.Sprintf("condition.%d.id", i)]
			_, _, err := governanceAPIClientForTestUtil.OktaGovernanceSDKClient().RequestConditionsAPI.
				GetResourceRequestConditionV2(context.Background(), resourceID, id).Execute()
			if err != nil {
				return fmt.Errorf("request condition %s for resource %s does not exist: %v", id, resourceID, err)
			}
		}
		return nil
	}
}

// checkRequestConditionsDestroy verifies that the managed conditions have been destroyed.
func checkRequestConditionsDestroy(s *terraform.State) error {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return nil
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resources.OktaGovernanceRequestConditions {
			continue
		}
		resourceID := rs.Primary.Attributes["resource_id"]
		list, _, err := governanceAPIClientForTestUtil.OktaGovernanceSDKClient().RequestConditionsAPI.
			ListResourceRequestConditionsV2(context.Background(), resourceID).Execute()
		if err != nil {
			return fmt.Errorf("error listing request conditions for resource %s: %v", resourceID, err)
		}
		for _, c := range list.GetData() {
			for _, name := range []string{"test-conditions-a", "test-conditions-b", "test-conditions-c", "test-conditions-d"} {
				if c.GetName() == name {
					return fmt.Errorf("request condition %s for resource %s still exists", name, resourceID)
				}
			}
		}
	}
	return nil
}
