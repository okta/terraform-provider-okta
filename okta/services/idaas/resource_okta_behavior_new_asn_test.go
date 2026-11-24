// Copyright 2025 - Present Okta, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package idaas_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccBehaviorNewAsn_basic(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_behavior_new_asn", t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", "okta_behavior_new_asn")
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkBehaviorNewAsnDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "NEW_ASN"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "settings.max_events_used_for_evaluation", "20"),
					resource.TestCheckResourceAttr(resourceName, "settings.min_events_needed_for_evaluation", "0"),
				),
			},
			{
				// Re-plan after apply must be empty (idempotency / no perpetual diff).
				Config:   config,
				PlanOnly: true,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "settings.max_events_used_for_evaluation", "50"),
					resource.TestCheckResourceAttr(resourceName, "settings.min_events_needed_for_evaluation", "5"),
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

func checkBehaviorNewAsnDestroy(s *terraform.State) error {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return nil
	}
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV9()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "okta_behavior_new_asn" {
			continue
		}
		_, resp, err := client.BehaviorAPI.GetBehaviorDetectionRule(context.Background(), rs.Primary.ID).Execute()
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			continue
		}
		if err != nil {
			return fmt.Errorf("error checking destroy of okta_behavior_new_asn %s: %w", rs.Primary.ID, err)
		}
		return fmt.Errorf("behavior_new_asn %s still exists", rs.Primary.ID)
	}
	return nil
}
