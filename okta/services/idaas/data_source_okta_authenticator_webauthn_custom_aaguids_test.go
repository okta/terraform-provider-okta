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
					resource.TestCheckResourceAttrSet(datasourceName, "custom_aaguids.#"),
				),
			},
		},
	})
}
