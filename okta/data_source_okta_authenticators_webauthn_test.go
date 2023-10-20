package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaAuthenticatorsWebauthn_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_authenticators_webauthn", t.Name())

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("datasource.tf", t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_authenticators_webauthn.test", "auth_devices.#"),
					resource.TestCheckResourceAttr("data.okta_authenticators_webauthn.test", "auth_devices.#", "106"),
					resource.TestCheckResourceAttrSet("data.okta_authenticators_webauthn.test", "auth_devices.0.aaguid"),
					resource.TestCheckResourceAttrSet("data.okta_authenticators_webauthn.test", "auth_devices.0.model_name"),
				),
			},
		},
	})
}
