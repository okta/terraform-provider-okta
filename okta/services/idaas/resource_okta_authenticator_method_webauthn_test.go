package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAuthenticatorMethodWebauthn_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.sample", resources.OktaIDaaSAuthenticatorMethodWebauthn)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthenticatorMethodWebauthn, t.Name())
	config := mgr.GetFixtures("okta_authenticator_method_webauthn.tf", t)
	updatedConfig := mgr.GetFixtures("okta_authenticator_method_webauthn_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "authenticator_id"),
					resource.TestCheckResourceAttr(resourceName, "user_verification", "PREFERRED"),
					resource.TestCheckResourceAttr(resourceName, "attachment", "ANY"),
					resource.TestCheckResourceAttr(resourceName, "aaguid_group.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "aaguid_group.0.name", "TestYubiKeys"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "user_verification", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "aaguid_group.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
