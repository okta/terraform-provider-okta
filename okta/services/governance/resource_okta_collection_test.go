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

func TestAccCollection_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceCollection, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceCollection)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		CheckDestroy:             testAccCheckCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc-collection-%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "Collection managed by Terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc-collection-%d-updated", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "description", "Collection managed by Terraform"),
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

// testAccCheckCollectionDestroy expects every okta_collection in state to read back as 404.
func testAccCheckCollectionDestroy(s *terraform.State) error {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return nil
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resources.OktaGovernanceCollection {
			continue
		}
		_, httpResp, err := governanceAPIClientForTestUtil.OktaGovernanceSDKClient().CollectionsV2API.GetCollectionV2(context.Background(), rs.Primary.ID).Execute()
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			continue
		}
		if err != nil {
			return fmt.Errorf("error checking destroy of %s %s: %w", rs.Type, rs.Primary.ID, err)
		}
		return fmt.Errorf("%s %s still exists", rs.Type, rs.Primary.ID)
	}
	return nil
}
