package governance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccCollection_basic(t *testing.T) {
	t.Skip("Skipping Collection tests - requires Okta Governance license")
	mgr := newFixtureManager("resources", resources.OktaGovernanceCollection, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceCollection)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Collection"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test collection for managing entitlements"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
		},
	})
}

func TestAccCollection_update(t *testing.T) {
	t.Skip("Skipping Collection tests - requires Okta Governance license")
	mgr := newFixtureManager("resources", resources.OktaGovernanceCollection, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	configUpdated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaGovernanceCollection)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test Collection"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test collection for managing entitlements"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Updated Collection"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccCollectionDataSource_basic(t *testing.T) {
	t.Skip("Skipping Collection data source tests - requires Okta Governance license")
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceCollection, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	dataSourceName := fmt.Sprintf("data.%s.test", resources.OktaGovernanceCollection)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created"),
				),
			},
		},
	})
}
