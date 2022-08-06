package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceAppMetadataSaml_read(t *testing.T) {
	mgr := newFixtureManager(appMetadataSaml, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	resourceName := "data.okta_app_metadata_saml.test"

	oktaResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "http_post_binding"),
					resource.TestCheckResourceAttrSet(resourceName, "metadata"),
					resource.TestCheckResourceAttrSet(resourceName, "entity_id"),
				),
			},
		},
	})
}
