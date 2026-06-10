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

func TestAccDelegateAppointments_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceDelegateAppointments, t.Name())
	config := mgr.ConfigReplace(mgr.GetFixtures("basic.tf", t))
	updatedConfig := mgr.ConfigReplace(mgr.GetFixtures("updated.tf", t))
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceDelegateAppointments)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkDelegateAppointmentsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "principal_id"),
					resource.TestCheckResourceAttr(resourceName, "principal_type", "OKTA_USER"),
					resource.TestCheckResourceAttr(resourceName, "appointments.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "appointments.0.delegate_id"),
					resource.TestCheckResourceAttr(resourceName, "appointments.0.note", "Covering while on PTO"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "appointments.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "appointments.0.delegate_id"),
					resource.TestCheckResourceAttr(resourceName, "appointments.0.note", "Switched delegate"),
				),
			},
		},
	})
}

func TestAccDelegateAppointments_import(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceDelegateAppointments, t.Name())
	config := mgr.ConfigReplace(mgr.GetFixtures("basic.tf", t))
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceDelegateAppointments)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "principal_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// checkDelegateAppointmentsDestroy verifies that delegate appointments have been removed
func checkDelegateAppointmentsDestroy(s *terraform.State) error {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return nil
	}

	client := governanceAPIClientForTestUtil

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resources.OktaGovernanceDelegateAppointments {
			continue
		}

		principalId := rs.Primary.Attributes["principal_id"]
		filter := fmt.Sprintf(`delegatorId eq "%s"`, principalId)

		listResp, apiResp, err := client.OktaGovernanceSDKClient().DelegatesAPI.ListDelegateAppointments(
			context.Background(),
		).Filter(filter).Execute()

		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			continue
		}

		if err != nil {
			return fmt.Errorf("error checking if delegate appointments for principal %s were destroyed: %v", principalId, err)
		}

		if len(listResp.Data) > 0 {
			return fmt.Errorf("delegate appointments for principal %s still exist (%d remaining)", principalId, len(listResp.Data))
		}
	}

	return nil
}
