package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAuthServerScopes(t *testing.T) {
	mgr := newFixtureManager(authServerScopes, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_auth_server_scopes.test", "scopes.#"),
				),
			},
		},
	})
}
