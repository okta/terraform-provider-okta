package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaDataSourceIdpMetadataSaml_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager("okta_idp_metadata_saml")
	config := mgr.GetFixtures("datasource.tf", ri, t)
	resourceName := "data.okta_idp_metadata_saml.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
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
