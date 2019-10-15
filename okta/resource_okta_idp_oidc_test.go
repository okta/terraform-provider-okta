package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func TestAccOktaIdpOidc_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(idpResource)
	config := mgr.GetFixtures("generic_oidc.tf", ri, t)
	updatedConfig := mgr.GetFixtures("generic_oidc_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", idpResource)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(idpResource, createDoesIdpExist(&sdk.OIDCIdentityProvider{})),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "acs_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "acs_type", "INSTANCE"),
					resource.TestCheckResourceAttr(resourceName, "authorization_url", "https://idp.example.com/authorize"),
					resource.TestCheckResourceAttr(resourceName, "authorization_binding", "HTTP-REDIRECT"),
					resource.TestCheckResourceAttr(resourceName, "token_url", "https://idp.example.com/token"),
					resource.TestCheckResourceAttr(resourceName, "token_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "user_info_url", "https://idp.example.com/userinfo"),
					resource.TestCheckResourceAttr(resourceName, "user_info_binding", "HTTP-REDIRECT"),
					resource.TestCheckResourceAttr(resourceName, "jwks_url", "https://idp.example.com/keys"),
					resource.TestCheckResourceAttr(resourceName, "jwks_binding", "HTTP-REDIRECT"),
					resource.TestCheckResourceAttr(resourceName, "client_id", "efg456"),
					resource.TestCheckResourceAttr(resourceName, "client_secret", "efg456"),
					resource.TestCheckResourceAttr(resourceName, "issuer_url", "https://id.example.com"),
					resource.TestCheckResourceAttr(resourceName, "username_template", "idpuser.email"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "acs_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "acs_type", "INSTANCE"),
					resource.TestCheckResourceAttr(resourceName, "authorization_url", "https://idp.example.com/authorize2"),
					resource.TestCheckResourceAttr(resourceName, "authorization_binding", "HTTP-REDIRECT"),
					resource.TestCheckResourceAttr(resourceName, "token_url", "https://idp.example.com/token2"),
					resource.TestCheckResourceAttr(resourceName, "token_binding", "HTTP-POST"),
					resource.TestCheckResourceAttr(resourceName, "user_info_url", "https://idp.example.com/userinfo2"),
					resource.TestCheckResourceAttr(resourceName, "user_info_binding", "HTTP-REDIRECT"),
					resource.TestCheckResourceAttr(resourceName, "jwks_url", "https://idp.example.com/keys2"),
					resource.TestCheckResourceAttr(resourceName, "jwks_binding", "HTTP-REDIRECT"),
					resource.TestCheckResourceAttr(resourceName, "client_id", "efg456"),
					resource.TestCheckResourceAttr(resourceName, "client_secret", "efg456"),
					resource.TestCheckResourceAttr(resourceName, "issuer_url", "https://id.example.com"),
					resource.TestCheckResourceAttr(resourceName, "username_template", "idpuser.email"),
				),
			},
		},
	})
}
