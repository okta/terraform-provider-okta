package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIdp(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(idpResource)
	config := mgr.GetFixtures("generic_oidc.tf", ri, t)
	updatedConfig := mgr.GetFixtures("generic_oidc_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", idpResource)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

func TestAccidpSocial(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(idpSocial)
	config := mgr.GetFixtures("basic.tf", ri, t)
	fbName := fmt.Sprintf("%s.facebook", idpSocial)
	microName := fmt.Sprintf("%s.microsoft", idpSocial)
	googleName := fmt.Sprintf("%s.google", idpSocial)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fbName, "type", "FACEBOOK"),
					resource.TestCheckResourceAttr(fbName, "protocol_type", "OAUTH2"),
					resource.TestCheckResourceAttr(fbName, "name", fmt.Sprintf("facebook_%d", ri)),
					resource.TestCheckResourceAttr(fbName, "client_id", "abcd123"),
					resource.TestCheckResourceAttr(fbName, "client_secret", "abcd123"),
					resource.TestCheckResourceAttr(fbName, "username_template", "idpuser.email"),

					resource.TestCheckResourceAttr(microName, "type", "MICROSOFT"),
					resource.TestCheckResourceAttr(microName, "protocol_type", "OIDC"),
					resource.TestCheckResourceAttr(microName, "name", fmt.Sprintf("microsoft_%d", ri)),
					resource.TestCheckResourceAttr(microName, "client_id", "abcd123"),
					resource.TestCheckResourceAttr(microName, "client_secret", "abcd123"),
					resource.TestCheckResourceAttr(microName, "username_template", "idpuser.userPrincipalName"),

					resource.TestCheckResourceAttr(googleName, "type", "GOOGLE"),
					resource.TestCheckResourceAttr(googleName, "protocol_type", "OAUTH2"),
					resource.TestCheckResourceAttr(googleName, "name", fmt.Sprintf("google_%d", ri)),
					resource.TestCheckResourceAttr(googleName, "client_id", "abcd123"),
					resource.TestCheckResourceAttr(googleName, "client_secret", "abcd123"),
					resource.TestCheckResourceAttr(googleName, "username_template", "idpuser.email"),
				),
			},
		},
	})
}
