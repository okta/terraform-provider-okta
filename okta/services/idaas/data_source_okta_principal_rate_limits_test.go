package idaas_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"testing"
)

func TestAccDataSourceOktaPrincipalRateLimits_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSPrincipalRateLimits, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_principal_rate_limits.test", "default_percentage", "50"),
					resource.TestCheckResourceAttr("data.okta_principal_rate_limits.test", "default_concurrency_percentage", "50"),
				),
			},
		},
	})
}
