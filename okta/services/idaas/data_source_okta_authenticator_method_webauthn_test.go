package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAuthenticatorMethodWebauthn_read(t *testing.T) {
	datasourceName := fmt.Sprintf("data.%s.sample", resources.OktaIDaaSAuthenticatorMethodWebauthn)
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAuthenticatorMethodWebauthn, t.Name())
	config := mgr.GetFixtures("okta_authenticator_method_webauthn.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "authenticator_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "status"),
				),
			},
		},
	})
}
