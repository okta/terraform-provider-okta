package governance_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// recordedOrgID is the org ID captured during VCR recording. Used during
// playback so request bodies match the recorded cassette.
const (
	recordedOrgID     = "00o3q4zilpjjJlLLu1d7"
	recordedOrnDomain = "oktapreview"
)

// orgIDForTest returns the org ID to use in test ORNs. During VCR playback,
// returns the recorded org ID so request bodies match the cassette.
func orgIDForTest() string {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return recordedOrgID
	}
	return os.Getenv("TF_VAR_org_id")
}

// ornDomainForTest returns the ORN domain to use in test ORNs.
func ornDomainForTest() string {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return recordedOrnDomain
	}
	return os.Getenv("TF_VAR_orn_domain")
}

// discoverEntitlementBundleORN uses the governance API to find an entitlement
// bundle in the sandbox that we can use for resource_owners testing.
// It also returns the parentResourceOrn (the app ORN).
// During VCR playback, returns hardcoded values from the recorded cassette
// since the API is not reachable.
func discoverEntitlementBundleORN(t *testing.T) (bundleOrn, parentOrn string) {
	t.Helper()

	// During VCR playback, the API is not reachable (fake credentials).
	// Use the values captured during VCR recording.
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return "orn:oktapreview:governance:00o3q4zilpjjJlLLu1d7:entitlement-bundles:enbe1l39ckVceLqn11d6",
			"orn:oktapreview:idp:00o3q4zilpjjJlLLu1d7:apps:saasure:0oa3q4zilxagApEOW1d7"
	}

	cfg := configForTest(t)

	// List entitlement bundles to find one we can use
	resp, _, err := cfg.OktaGovernanceClient.OktaGovernanceSDKClient().EntitlementBundlesAPI.
		ListEntitlementBundles(context.Background()).
		Execute()
	if err != nil {
		t.Fatalf("Failed to list entitlement bundles: %v", err)
	}

	bundles := resp.GetData()
	if len(bundles) == 0 {
		t.Skip("No entitlement bundles found in sandbox — cannot test resource_owners")
	}

	bundle := bundles[0]
	return bundle.GetOrn(), bundle.GetTargetResourceOrn()
}

func TestAccResourceOktaResourceOwners_basicCreateDestroy(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceResourceOwners, t.Name())

	bundleOrn, _ := discoverEntitlementBundleORN(t)
	orgID := orgIDForTest()
	ornDomain := ornDomainForTest()

	config := fmt.Sprintf(`
resource "okta_user" "owner1" {
  first_name = "TestAcc"
  last_name  = "ResourceOwner1"
  login      = "testAcc-resource-owner1-%s@example.com"
  email      = "testAcc-resource-owner1-%s@example.com"
}

resource "okta_resource_owners" "test" {
  resource_orn = "%s"

  principal_orns = [
    "orn:%s:directory:%s:users:${okta_user.owner1.id}",
  ]
}
`, mgr.SeedStr(), mgr.SeedStr(), bundleOrn, ornDomain, orgID)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("okta_resource_owners.test", "id"),
						resource.TestCheckResourceAttr("okta_resource_owners.test", "resource_orn", bundleOrn),
						resource.TestCheckResourceAttrSet("okta_resource_owners.test", "parent_resource_orn"),
						resource.TestCheckResourceAttr("okta_resource_owners.test", "principal_orns.#", "1"),
					),
				},
			},
		})
}

func TestAccResourceOktaResourceOwners_crudAndUpdate(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceResourceOwners, t.Name())

	bundleOrn, _ := discoverEntitlementBundleORN(t)
	orgID := orgIDForTest()
	ornDomain := ornDomainForTest()

	configCreate := fmt.Sprintf(`
resource "okta_user" "owner1" {
  first_name = "TestAcc"
  last_name  = "ResourceOwner1"
  login      = "testAcc-resource-owner1-%s@example.com"
  email      = "testAcc-resource-owner1-%s@example.com"
}

resource "okta_user" "owner2" {
  first_name = "TestAcc"
  last_name  = "ResourceOwner2"
  login      = "testAcc-resource-owner2-%s@example.com"
  email      = "testAcc-resource-owner2-%s@example.com"
}

resource "okta_resource_owners" "test" {
  resource_orn = "%s"

  principal_orns = [
    "orn:%s:directory:%s:users:${okta_user.owner1.id}",
    "orn:%s:directory:%s:users:${okta_user.owner2.id}",
  ]
}
`, mgr.SeedStr(), mgr.SeedStr(), mgr.SeedStr(), mgr.SeedStr(),
		bundleOrn, ornDomain, orgID, ornDomain, orgID)

	configUpdate := fmt.Sprintf(`
resource "okta_user" "owner1" {
  first_name = "TestAcc"
  last_name  = "ResourceOwner1"
  login      = "testAcc-resource-owner1-%s@example.com"
  email      = "testAcc-resource-owner1-%s@example.com"
}

resource "okta_user" "owner2" {
  first_name = "TestAcc"
  last_name  = "ResourceOwner2"
  login      = "testAcc-resource-owner2-%s@example.com"
  email      = "testAcc-resource-owner2-%s@example.com"
}

resource "okta_user" "owner3" {
  first_name = "TestAcc"
  last_name  = "ResourceOwner3"
  login      = "testAcc-resource-owner3-%s@example.com"
  email      = "testAcc-resource-owner3-%s@example.com"
}

resource "okta_resource_owners" "test" {
  resource_orn = "%s"

  principal_orns = [
    "orn:%s:directory:%s:users:${okta_user.owner1.id}",
    "orn:%s:directory:%s:users:${okta_user.owner3.id}",
  ]
}
`, mgr.SeedStr(), mgr.SeedStr(), mgr.SeedStr(), mgr.SeedStr(), mgr.SeedStr(), mgr.SeedStr(),
		bundleOrn, ornDomain, orgID, ornDomain, orgID)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: configCreate,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("okta_resource_owners.test", "id"),
						resource.TestCheckResourceAttr("okta_resource_owners.test", "resource_orn", bundleOrn),
						resource.TestCheckResourceAttr("okta_resource_owners.test", "principal_orns.#", "2"),
					),
				},
				{
					// Update: remove owner2, add owner3
					Config: configUpdate,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("okta_resource_owners.test", "principal_orns.#", "2"),
						resource.TestCheckResourceAttrSet("okta_resource_owners.test", "parent_resource_orn"),
					),
				},
				{
					ResourceName:      "okta_resource_owners.test",
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateIdFunc: func(s *terraform.State) (string, error) {
						rs := s.RootModule().Resources["okta_resource_owners.test"]
						pOrn := rs.Primary.Attributes["parent_resource_orn"]
						rOrn := rs.Primary.Attributes["resource_orn"]
						return fmt.Sprintf("%s/%s", pOrn, rOrn), nil
					},
				},
			},
		})

}
