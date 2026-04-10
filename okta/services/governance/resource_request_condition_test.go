package governance_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		CheckDestroy:             checkRequestConditionDestroy,
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
		CheckDestroy:             checkRequestConditionDestroy,
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

func TestAccRequestConditionResource_Status(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceRequestCondition, t.Name())
	configActive := mgr.GetFixtures("status_active.tf", t)
	configInactive := mgr.GetFixtures("status_inactive.tf", t)
	configReactivated := mgr.GetFixtures("status_active.tf", t)
	resourceName := fmt.Sprintf("%s.test_status", resources.OktaGovernanceRequestCondition)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRequestConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config: configActive,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-condition-status"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: configInactive,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-condition-status"),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
				),
			},
			{
				Config: configReactivated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-condition-status"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

// TestAccRequestConditionResource_Priority verifies that creating two conditions
// with explicit priorities does not cause "inconsistent result after apply" errors
// when the API silently reassigns priorities. It also swaps priorities in a second
// step to exercise the Update path.
func TestAccRequestConditionResource_Priority(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceRequestCondition, t.Name())
	config := mgr.GetFixtures("priority.tf", t)
	updatedConfig := mgr.GetFixtures("priority_updated.tf", t)
	firstName := fmt.Sprintf("%s.priority_first", resources.OktaGovernanceRequestCondition)
	secondName := fmt.Sprintf("%s.priority_second", resources.OktaGovernanceRequestCondition)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRequestConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(firstName, "name", "priority-first"),
					resource.TestCheckResourceAttr(firstName, "priority", "1"),
					resource.TestCheckResourceAttr(secondName, "name", "priority-second"),
					resource.TestCheckResourceAttr(secondName, "priority", "2"),
				),
			},
			{
				Config:             mgr.ConfigReplace(updatedConfig),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(firstName, "name", "priority-first"),
					resource.TestCheckResourceAttr(firstName, "priority", "2"),
					resource.TestCheckResourceAttr(secondName, "name", "priority-second"),
					resource.TestCheckResourceAttr(secondName, "priority", "1"),
				),
			},
		},
	})
}

// TestAccRequestConditionResource_PriorityNonSequential verifies that creating
// a condition with a non-sequential priority (e.g., 100 when only 2 conditions
// exist) does not cause "inconsistent result after apply" errors. This hits the
// Create code path directly — the API is likely to normalize the priority down
// to the actual position, which is the most common trigger for the bug.
func TestAccRequestConditionResource_PriorityNonSequential(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceRequestCondition, t.Name())
	config := mgr.GetFixtures("priority_non_sequential.tf", t)
	firstName := fmt.Sprintf("%s.priority_first", resources.OktaGovernanceRequestCondition)
	secondName := fmt.Sprintf("%s.priority_second", resources.OktaGovernanceRequestCondition)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkRequestConditionDestroy,
		Steps: []resource.TestStep{
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(firstName, "name", "priority-first"),
					resource.TestCheckResourceAttr(firstName, "priority", "1"),
					resource.TestCheckResourceAttr(secondName, "name", "priority-second"),
					resource.TestCheckResourceAttr(secondName, "priority", "100"),
				),
			},
		},
	})
}

// checkRequestConditionDestroy verifies that request conditions have been destroyed
func checkRequestConditionDestroy(s *terraform.State) error {
	// Skip destroy check in VCR playback mode
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return nil
	}

	// Use the shared governance client
	client := governanceAPIClientForTestUtil

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resources.OktaGovernanceRequestCondition {
			continue
		}

		resourceID := rs.Primary.Attributes["resource_id"]
		conditionID := rs.Primary.ID

		// Try to get the request condition
		_, resp, err := client.OktaGovernanceSDKClient().RequestConditionsAPI.GetResourceRequestConditionV2(
			context.Background(),
			resourceID,
			conditionID,
		).Execute()

		// If we get a 404, the resource is successfully deleted
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			continue
		}

		// If there's an error other than 404, return it
		if err != nil {
			return fmt.Errorf("error checking if request condition %s was destroyed: %v", conditionID, err)
		}

		// If we got here, the resource still exists
		return fmt.Errorf("request condition %s for resource %s still exists", conditionID, resourceID)
	}

	return nil
}
