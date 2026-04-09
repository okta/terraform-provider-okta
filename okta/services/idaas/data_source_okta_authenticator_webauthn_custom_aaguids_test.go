package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAuthenticatorWebauthnCustomAAGUIDs_read(t *testing.T) {
	datasourceName := fmt.Sprintf("data.%s.sample", resources.OktaIDaaSAuthenticatorWebauthnCustomAAGUIDs)
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAuthenticatorWebauthnCustomAAGUIDs, t.Name())
	config := mgr.GetFixtures("okta_authenticator_webauthn_custom_aaguids.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "authenticator_id"),
					resource.TestCheckResourceAttr(datasourceName, "custom_aaguids.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "custom_aaguids.0.aaguid", "cb69481e-8ff7-4039-93ec-0a2729a154a8"),
					resource.TestCheckResourceAttr(datasourceName, "custom_aaguids.0.name", "YubiKey 5 Series"),
				),
			},
		},
	})
}
