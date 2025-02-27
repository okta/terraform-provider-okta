package idaas_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaOrgConfiguration_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSOrgConfiguration)
	mgr := newFixtureManager("resources", resources.OktaIDaaSOrgConfiguration, t.Name())
	config := mgr.GetFixtures("standard.tf", t)
	updatedConfig := mgr.GetFixtures("standard_updated.tf", t)
	var originalCompanyName string
	companyName := fmt.Sprintf("testAcc-%d Hashicorp CI Terraform Provider Okta", mgr.Seed)
	companyNameUpdated := fmt.Sprintf("testAcc-%d Hashicorp CI Terraform Provider Okta Updated", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
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
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
		client.OrgSetting.PartialUpdateOrgSetting(context.Background(), setting)
		return nil
	}
}

// setupGetOriginalCompanyName Get the original org name so it can be rewritten
// back at teardown.
func setupGetOriginalCompanyName(companyName *string) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
		if settings, _, err := client.OrgSetting.GetOrgSettings(context.Background()); err == nil {
			*companyName = settings.CompanyName
		}
		return nil
	}
}
