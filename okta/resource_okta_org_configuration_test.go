package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaOrgConfiguration(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", orgConfiguration)
	mgr := newFixtureManager(orgConfiguration, t.Name())
	config := mgr.GetFixtures("standard.tf", t)
	updatedConfig := mgr.GetFixtures("standard_updated.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "company_name", "Hashicorp CI Terraform Provider Okta"),
					resource.TestCheckResourceAttr(resourceName, "website", "https://terraform.io"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "company_name", "Hashicorp CI Terraform Provider Okta Updated"),
					resource.TestCheckResourceAttr(resourceName, "website", "https://terraform.com"),
					resource.TestCheckResourceAttr(resourceName, "phone_number", strconv.Itoa(mgr.Seed)),
				),
			},
			{
				Config: config,
			},
		},
	})
}
