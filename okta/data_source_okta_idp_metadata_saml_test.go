package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceIdpMetadataSaml_read(t *testing.T) {
	mgr := newFixtureManager(idpMetadataSaml, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	resourceName := "data.okta_idp_metadata_saml.test"

	oktaResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "signing_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "encryption_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "http_post_binding"),
					resource.TestCheckResourceAttrSet(resourceName, "metadata"),
					resource.TestCheckResourceAttrSet(resourceName, "entity_id"),
				),
			},
		},
	})
}
