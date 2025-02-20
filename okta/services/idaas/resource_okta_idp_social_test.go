package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaIdpSocial_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSIdpSocial, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	disabledConf := mgr.GetFixtures("auto_provision_disabled.tf", t)
	fbName := fmt.Sprintf("%s.facebook", resources.OktaIDaaSIdpSocial)
	microName := fmt.Sprintf("%s.microsoft", resources.OktaIDaaSIdpSocial)
	googleName := fmt.Sprintf("%s.google", resources.OktaIDaaSIdpSocial)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSIdpSocial, createDoesIdpExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fbName, "type", "FACEBOOK"),
					resource.TestCheckResourceAttr(fbName, "protocol_type", "OAUTH2"),
					resource.TestCheckResourceAttr(fbName, "name", fmt.Sprintf("testAcc_facebook_%d", mgr.Seed)),
					resource.TestCheckResourceAttr(fbName, "client_id", "abcd123"),
					resource.TestCheckResourceAttr(fbName, "client_secret", "abcd123"),
					resource.TestCheckResourceAttr(fbName, "username_template", "idpuser.email"),

					resource.TestCheckResourceAttr(microName, "type", "MICROSOFT"),
					resource.TestCheckResourceAttr(microName, "protocol_type", "OIDC"),
					resource.TestCheckResourceAttr(microName, "name", fmt.Sprintf("testAcc_microsoft_%d", mgr.Seed)),
					resource.TestCheckResourceAttr(microName, "client_id", "abcd123"),
					resource.TestCheckResourceAttr(microName, "client_secret", "abcd123"),
					resource.TestCheckResourceAttr(microName, "username_template", "idpuser.userPrincipalName"),
					resource.TestCheckResourceAttr(microName, "groups_action", "ASSIGN"),
					resource.TestCheckResourceAttr(microName, "groups_assignment.#", "1"),

					resource.TestCheckResourceAttr(googleName, "type", "GOOGLE"),
					resource.TestCheckResourceAttr(googleName, "protocol_type", "OIDC"),
					resource.TestCheckResourceAttr(googleName, "name", fmt.Sprintf("testAcc_google_%d", mgr.Seed)),
					resource.TestCheckResourceAttr(googleName, "client_id", "abcd123"),
					resource.TestCheckResourceAttr(googleName, "client_secret", "abcd123"),
					resource.TestCheckResourceAttr(googleName, "username_template", "idpuser.email"),
				),
			},
			{
				Config: disabledConf,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(googleName, "type", "GOOGLE"),
					resource.TestCheckResourceAttr(googleName, "protocol_type", "OIDC"),
					resource.TestCheckResourceAttr(googleName, "name", fmt.Sprintf("testAcc_google_%d", mgr.Seed)),
					resource.TestCheckResourceAttr(googleName, "client_id", "abcd123"),
					resource.TestCheckResourceAttr(googleName, "client_secret", "abcd123"),
					resource.TestCheckResourceAttr(googleName, "username_template", "idpuser.email"),
					resource.TestCheckResourceAttr(googleName, "provisioning_action", "DISABLED"),
				),
			},
		},
	})
}
