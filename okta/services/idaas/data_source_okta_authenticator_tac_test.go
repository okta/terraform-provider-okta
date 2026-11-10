package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAuthenticatorTac_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAuthenticatorTac, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	resourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSAuthenticatorTac)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", "data.okta_authenticator.tac", "id"),
					resource.TestCheckResourceAttr(resourceName, "type", "tac"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "last_updated"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_tac.type"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_tac.configuration.length"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_tac.configuration.default_ttl"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_tac.configuration.min_ttl"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_tac.configuration.max_ttl"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_tac.configuration.multi_use_allowed"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_tac.configuration.complexity.letters"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_tac.configuration.complexity.numbers"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_tac.configuration.complexity.special_characters"),
				),
			},
		},
	})
}
