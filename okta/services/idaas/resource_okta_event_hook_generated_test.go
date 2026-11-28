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
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// testAccCheckEventHookDestroy verifies every okta_event_hook in state was deleted
// upstream by reading it back and expecting a 404.
func testAccCheckEventHookDestroy(s *terraform.State) error {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return nil
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resources.OktaIDaaSEventHook {
			continue
		}
		_, httpResp, err := iDaaSAPIClientForTestUtil.OktaSDKClientV9().EventHookAPI.GetEventHook(
			context.Background(),
			rs.Primary.ID,
		).Execute()
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			continue
		}
		if err != nil {
			return fmt.Errorf("error checking destroy of %s %s: %w", resources.OktaIDaaSEventHook, rs.Primary.ID, err)
		}
		return fmt.Errorf("%s %s still exists after destroy", resources.OktaIDaaSEventHook, rs.Primary.ID)
	}
	return nil
}

func TestAccResourceOktaEventHook_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSEventHook, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSEventHook)
	name := acctest.BuildResourceName(mgr.Seed)
	uri := fmt.Sprintf("https://example.com/%s", name)
	filterExpr := "event.target.?[type eq 'User'].size()>0"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             testAccCheckEventHookDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					// Create activates the hook after POST.
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "verification_status", "UNVERIFIED"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.config.uri", uri),
					resource.TestCheckResourceAttr(resourceName, "channel.config.method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "events.type", "EVENT_TYPE"),
					resource.TestCheckResourceAttr(resourceName, "events.items.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "events.items.0", "user.lifecycle.create"),
					resource.TestCheckResourceAttr(resourceName, "events.items.1", "user.lifecycle.delete.initiated"),
					// Single nested block inside a list item — regression for the write-path fix in db9bd7f5.
					resource.TestCheckResourceAttr(resourceName, "events.filter.event_filter_map.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "events.filter.event_filter_map.0.event", "user.lifecycle.create"),
					resource.TestCheckResourceAttr(resourceName, "events.filter.event_filter_map.0.condition.expression", filterExpr),
				),
			},
			{
				// Re-plan after apply must be empty (no perpetual diff).
				Config:   config,
				PlanOnly: true,
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", name+" Updated"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "channel.config.uri", uri),
					resource.TestCheckResourceAttr(resourceName, "events.filter.event_filter_map.0.condition.expression", filterExpr),
				),
			},
			{
				Config:   updatedConfig,
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
