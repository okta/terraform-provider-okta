package governance_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccGrant_custom(t *testing.T) {
	t.Skip("Skipping Grant tests - requires Okta Governance license")
	mgr := newFixtureManager("resources", "okta_grant", t.Name())
	config := mgr.GetFixtures("custom.tf", t)
	resourceName := "okta_grant.test"

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
					resource.TestCheckResourceAttr(resourceName, "grant_type", "CUSTOM"),
					resource.TestCheckResourceAttrSet(resourceName, "target_principal_id"),
					resource.TestCheckResourceAttr(resourceName, "target_principal_type", "OKTA_GROUP"),
					resource.TestCheckResourceAttrSet(resourceName, "target_resource_orn"),
					resource.TestCheckResourceAttr(resourceName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "actor", "API"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "1"),
				),
			},
		},
	})
}

func TestAccGrant_invalidCombination_customWithBundle(t *testing.T) {
	t.Skip("Skipping Grant tests - requires Okta Governance license")
	mgr := newFixtureManager("resources", "okta_grant", t.Name())
	// hypothetical fixture combining custom + bundle id
	config := `resource "okta_grant" "test" {
	   grant_type = "CUSTOM"
	   target_principal_id = "00uDummy"
	   target_principal_type = "OKTA_USER"
	   target_resource_orn = "orn:okta:idp:00oOrg:apps:dummy:0oaApp"
	   entitlement_bundle_id = "ebDummy"
	 }`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{{
			Config:      config,
			ExpectError: regexp.MustCompile("entitlement_bundle_id must not be set"),
		}},
	})
}

func TestAccGrant_invalidCombination_bundleWithEntitlements(t *testing.T) {
	t.Skip("Skipping Grant tests - requires Okta Governance license")
	config := `resource "okta_grant" "test" {
	   grant_type = "ENTITLEMENT-BUNDLE"
	   target_principal_id = "00uDummy"
	   target_principal_type = "OKTA_USER"
	   target_resource_orn = "orn:okta:idp:00oOrg:apps:dummy:0oaApp"
	   entitlement_bundle_id = "ebDummy"
	   entitlements { id = "ent1" }
	 }`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{{
			Config:      config,
			ExpectError: regexp.MustCompile("entitlements must not be set"),
		}},
	})
}

func TestAccGrant_invalidTimeZoneWithoutExpiration(t *testing.T) {
	t.Skip("Skipping Grant tests - requires Okta Governance license")
	config := `resource "okta_grant" "test" {
	   grant_type = "CUSTOM"
	   target_principal_id = "00uDummy"
	   target_principal_type = "OKTA_USER"
	   target_resource_orn = "orn:okta:idp:00oOrg:apps:dummy:0oaApp"
	   time_zone = "UTC"
	   entitlements { id = "ent1" }
	 }`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{{
			Config:      config,
			ExpectError: regexp.MustCompile("time_zone requires expiration_date"),
		}},
	})
}

func TestAccGrant_pastExpiration(t *testing.T) {
	t.Skip("Skipping Grant tests - requires Okta Governance license")
	past := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	config := fmt.Sprintf(`resource "okta_grant" "test" {
	   grant_type = "CUSTOM"
	   target_principal_id = "00uDummy"
	   target_principal_type = "OKTA_USER"
	   target_resource_orn = "orn:okta:idp:00oOrg:apps:dummy:0oaApp"
	   expiration_date = "%s"
	   entitlements { id = "ent1" }
	 }`, past)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{{
			Config:      config,
			ExpectError: regexp.MustCompile("expiration_date must be a future timestamp"),
		}},
	})
}

func TestAccGrant_withExpiration(t *testing.T) {
	t.Skip("Skipping Grant tests - requires Okta Governance license")
	mgr := newFixtureManager("resources", "okta_grant", t.Name())
	config := mgr.GetFixtures("with_expiration.tf", t)
	resourceName := "okta_grant.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "expiration_date"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "America/New_York"),
				),
			},
		},
	})
}

func TestAccGrant_update(t *testing.T) {
	t.Skip("Skipping Grant tests - requires Okta Governance license")
	mgr := newFixtureManager("resources", "okta_grant", t.Name())
	config := mgr.GetFixtures("custom.tf", t)
	configUpdated := mgr.GetFixtures("updated_expiration.tf", t)
	resourceName := "okta_grant.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "expiration_date"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "expiration_date"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC"),
				),
			},
		},
	})
}

func TestAccGrantDataSource_basic(t *testing.T) {
	// existing content below

	t.Skip("Skipping Grant data source tests - requires Okta Governance license")
	mgr := newFixtureManager("data-sources", "okta_grant", t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	dataSourceName := "data.okta_grant.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "grant_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "target_principal_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "target_principal_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
				),
			},
		},
	})
}
