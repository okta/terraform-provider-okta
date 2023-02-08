package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaDomain_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(domain)
	config := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config:  config,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_domain.by-id", "domain", "www.example.com"),
					resource.TestCheckResourceAttr("data.okta_domain.by-name", "domain", "www.example.com"),
				),
			},
		},
	})
}
