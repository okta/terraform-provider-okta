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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// Creates an assignment without an expiry, then adds expiration_time/time_zone. The update is
// sent as JSON-Patch ADD ops, which only works if the unset fields were stored as null.
func TestAccCollectionsAssignment_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceCollectionsAssignment, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceCollectionsAssignment)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		// The collection's destroy check covers this too: deleting the collection removes its assignments.
		CheckDestroy: testAccCheckCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "collection_id", "okta_collection.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "principal.external_id", "00uto13wvy2Y5RCv91d7"),
					resource.TestCheckResourceAttr(resourceName, "principal.type", "OKTA_USER"),
					resource.TestCheckNoResourceAttr(resourceName, "expiration_time"),
					resource.TestCheckNoResourceAttr(resourceName, "time_zone"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_time", "2027-12-31T23:59:59Z"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "Asia/Kolkata"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: importIDFromParent(resourceName, "collection_id"),
				ImportStateVerify: true,
			},
		},
	})
}

// importIDFromParent builds the "{<parentAttr>}/{id}" import ID used by collection child resources.
func importIDFromParent(resourceName, parentAttr string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes[parentAttr], rs.Primary.ID), nil
	}
}
