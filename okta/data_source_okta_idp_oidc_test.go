package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceIdpOidc_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(idpOidc)
	idpOidcConfig := mgr.GetFixtures("generic_oidc.tf", ri, t)
	config := mgr.GetFixtures("datasource.tf", ri, t)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: idpOidcConfig,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_idp_oidc.test", "id"),
				),
			},
		},
	})
}
