package okta

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccOktaOrgConfiguration(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", orgConfiguration)
	mgr := newFixtureManager(orgConfiguration, t.Name())
	config := mgr.GetFixtures("standard.tf", t)
	updatedConfig := mgr.GetFixtures("standard_updated.tf", t)
	var originalCompanyName string
	companyName := fmt.Sprintf("testAcc-%d Hashicorp CI Terraform Provider Okta", mgr.Seed)
	companyNameUpdated := fmt.Sprintf("testAcc-%d Hashicorp CI Terraform Provider Okta Updated", mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			// We get around the TF testing runtime not having good setup and
			// teardown by using a step to get the name of the current org,
			// saving it, and resetting it back afterwards in a teardown.
			// This saves from renaming our org something like
			// testAcc-123 Hashicorp CI Terraform Provider Okta
			{
				// setup
				Config: `data "okta_groups" "test" { type = "BUILT_IN" }`,
				Check: resource.ComposeTestCheckFunc(
					setupGetOriginalCompanyName(&originalCompanyName),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "company_name", companyName),
					resource.TestCheckResourceAttr(resourceName, "website", "https://terraform.io"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "company_name", companyNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "website", "https://terraform.com"),
					resource.TestCheckResourceAttr(resourceName, "phone_number", strconv.Itoa(mgr.Seed)),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "company_name", companyName),
					resource.TestCheckResourceAttr(resourceName, "website", "https://terraform.io"),
					resource.TestCheckResourceAttr(resourceName, "phone_number", ""),
				),
			},
			{
				// teardown
				Config: `data "okta_groups" "test" { type = "BUILT_IN" }`,
				Check: resource.ComposeTestCheckFunc(
					teardownResetCompanyName(&originalCompanyName),
				),
			},
		},
	})
}

// teardownResetCompanyName Reset the company name back to its original.
func teardownResetCompanyName(name *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if name == nil {
			return nil
		}
		if *name == "" {
			return nil
		}
		setting := sdk.OrgSetting{
			CompanyName: *name + " ?",
		}
		sdkV2ClientForTest().OrgSetting.PartialUpdateOrgSetting(context.Background(), setting)
		return nil
	}
}

// setupGetOriginalCompanyName Get the original org name so it can be rewritten
// back at teardown.
func setupGetOriginalCompanyName(companyName *string) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		client := sdkV2ClientForTest()
		if settings, _, err := client.OrgSetting.GetOrgSettings(context.Background()); err == nil {
			*companyName = settings.CompanyName
		}
		return nil
	}
}
