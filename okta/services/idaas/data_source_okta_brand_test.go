package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaBrand_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSBrand, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config:  config,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_brand.example", "id"),
					resource.TestCheckResourceAttrSet("data.okta_brand.example", "name"),
					resource.TestCheckResourceAttrSet("data.okta_brand.example", "links"),
					resource.TestCheckResourceAttrSet("data.okta_brand.default", "id"),
					resource.TestCheckResourceAttrSet("data.okta_brand.default", "links"),
				),
			},
		},
	})
}
