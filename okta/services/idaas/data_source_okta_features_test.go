package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaFeature_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSFeature, t.Name())
	config := mgr.GetFixtures("data-source.tf", t)
	appCreate := buildTestApp(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: appCreate,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_features.test", "substring"),
					resource.TestCheckResourceAttr("data.okta_features.test", "substring", "MFA"),
				),
			},
		},
	})
}
