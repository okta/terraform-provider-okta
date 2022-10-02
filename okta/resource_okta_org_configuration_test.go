package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaOrgConfiguration(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", orgConfiguration)
	mgr := newFixtureManager(orgConfiguration)
	config := mgr.GetFixtures("standard.tf", ri, t)
	updatedConfig := mgr.GetFixtures("standard_updated.tf", ri, t)
	companyName := fmt.Sprintf("testAcc-%d Hashicorp CI Terraform Provider Okta", ri)
	companyNameUpdated := fmt.Sprintf("testAcc-%d Hashicorp CI Terraform Provider Okta Updated", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
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
					resource.TestCheckResourceAttr(resourceName, "phone_number", strconv.Itoa(ri)),
				),
			},
			{
				Config: config,
			},
		},
	})
}
