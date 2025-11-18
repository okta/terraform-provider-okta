package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccCollectionAssignment_basic(t *testing.T) {
	t.Skip("Skipping Collection Assignment tests - requires Okta Governance license")
	mgr := newFixtureManager("resources", "okta_collection_assignment", t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := "okta_collection_assignment.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "collection_id"),
					resource.TestCheckResourceAttrSet(resourceName, "principal_id"),
					resource.TestCheckResourceAttr(resourceName, "principal_type", "OKTA_GROUP"),
					resource.TestCheckResourceAttr(resourceName, "actor", "API"),
				),
			},
		},
	})
}

func TestAccCollectionAssignment_withExpiration(t *testing.T) {
	t.Skip("Skipping Collection Assignment tests - requires Okta Governance license")
	mgr := newFixtureManager("resources", "okta_collection_assignment", t.Name())
	config := mgr.GetFixtures("with_expiration.tf", t)
	resourceName := "okta_collection_assignment.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "expiration_time"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "America/Los_Angeles"),
				),
			},
		},
	})
}

func TestAccCollectionAssignment_update(t *testing.T) {
	t.Skip("Skipping Collection Assignment tests - requires Okta Governance license")
	mgr := newFixtureManager("resources", "okta_collection_assignment", t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	configUpdated := mgr.GetFixtures("updated.tf", t)
	resourceName := "okta_collection_assignment.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "expiration_time"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "expiration_time"),
				),
			},
		},
	})
}

func TestAccCollectionAssignmentDataSource_basic(t *testing.T) {
	t.Skip("Skipping Collection Assignment data source tests - requires Okta Governance license")
	mgr := newFixtureManager("data-sources", "okta_collection_assignment", t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	dataSourceName := "data.okta_collection_assignment.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "collection_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "principal_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "principal_type"),
				),
			},
		},
	})
}
