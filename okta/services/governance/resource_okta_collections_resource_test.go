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
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// Adds an app to a collection with one entitlement value, then swaps it for another — the
// update must read the new entitlement ids back from the PUT response.
func TestAccCollectionsResource_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceCollectionsResource, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceCollectionsResource)
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
					resource.TestCheckResourceAttrPair(resourceName, "id", resourceName, "resource_id"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.0.id", "esp1mrruxopcczTEr1d7"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.0.values.0.id", "ent1mrruxps9YxXFb1d7"),
					resource.TestCheckResourceAttr(resourceName, "resource_configuration.type", "APPLICATION"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.0.id", "esp1mpfbbbqD2NAwG1d7"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.0.values.0.id", "ent1mpfbbczwILObX1d7"),
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
