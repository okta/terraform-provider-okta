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

func TestAccDataSourceOktaCollectionsResource_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceCollectionsResource, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := fmt.Sprintf("data.%s.test", resources.OktaGovernanceCollectionsResource)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		CheckDestroy:             testAccCheckCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", "okta_collections_resource.test", "resource_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "resource_orn", "okta_collections_resource.test", "resource_orn"),
					resource.TestCheckResourceAttr(dataSourceName, "resource_configuration.type", "APPLICATION"),
					resource.TestCheckResourceAttr(dataSourceName, "resource_configuration.counts.entitlements", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "entitlements.#", "1"),
					// labels, related_apps and resource_profile are returned only with ?include=..., which
					// the data source does not send yet, so they are not asserted here.
				),
			},
		},
	})
}
