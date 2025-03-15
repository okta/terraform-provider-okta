package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaBrands_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSBrands, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.#"),
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.0.id"),
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.0.name"),
					resource.TestCheckResourceAttrSet("data.okta_brands.test", "brands.0.links"),
				),
			},
		},
	})
}
