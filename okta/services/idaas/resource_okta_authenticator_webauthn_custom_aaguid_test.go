package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAuthenticatorWebauthnCustomAAGUID_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.sample", resources.OktaIDaaSAuthenticatorWebauthnCustomAAGUID)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthenticatorWebauthnCustomAAGUID, t.Name())
	config := mgr.GetFixtures("okta_authenticator_webauthn_custom_aaguid.tf", t)
	updatedConfig := mgr.GetFixtures("okta_authenticator_webauthn_custom_aaguid_updated.tf", t)

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
					resource.TestCheckResourceAttr(resourceName, "aaguid", "cb69481e-8ff7-4039-93ec-0a2729a154a8"),
					resource.TestCheckResourceAttr(resourceName, "name", "YubiKey 5 Series"),
					resource.TestCheckResourceAttr(resourceName, "authenticator_characteristics.fips_compliant", "true"),
					resource.TestCheckResourceAttr(resourceName, "authenticator_characteristics.hardware_protected", "true"),
					resource.TestCheckResourceAttr(resourceName, "authenticator_characteristics.platform_attached", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "YubiKey 5 Series (Updated)"),
					resource.TestCheckResourceAttr(resourceName, "authenticator_characteristics.platform_attached", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["authenticator_id"], rs.Primary.Attributes["aaguid"]), nil
				},
			},
		},
	})
}
