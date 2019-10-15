package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func TestAccOktaIdpSaml_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(idpSaml)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", idpSaml)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(idpSaml, createDoesIdpExist(&sdk.SAMLIdentityProvider{})),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "acs_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "acs_type", "INSTANCE"),
					resource.TestCheckResourceAttrSet(resourceName, "audience"),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "sso_destination", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "sso_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "username_template", "idpuser.email"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "https://idp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_scope", "REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_scope", "ANY"),
					resource.TestCheckResourceAttrSet(resourceName, "kid"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "acs_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "acs_type", "INSTANCE"),
					resource.TestCheckResourceAttrSet(resourceName, "audience"),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "https://idp.example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "sso_destination", "https://idp.example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "sso_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "username_template", "idpuser.email"),
					resource.TestCheckResourceAttr(resourceName, "issuer", "https://idp.example.com/issuer"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr(resourceName, "request_signature_scope", "REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "response_signature_scope", "RESPONSE"),
					resource.TestCheckResourceAttrSet(resourceName, "kid"),
				),
			},
		},
	})
}
