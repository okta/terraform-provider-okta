package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// TestAccResourceOktaAuthenticatorTac_crud imports the org's existing TAC authenticator rather
// than creating one: an org can hold only one (create returns 409 Conflict) and Delete is a
// no-op. updated.tf holds the org's original settings so the last step restores them.
func TestAccResourceOktaAuthenticatorTac_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthenticatorTac, t.Name())
	lookup := mgr.GetFixtures("lookup.tf", t)
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthenticatorTac)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: lookup,
			},
			{
				Config:             config,
				ResourceName:       resourceName,
				ImportState:        true,
				ImportStateIdFunc:  importStateIdForAuthenticatorTac,
				ImportStatePersist: true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", "data.okta_authenticator.tac", "id"),
					resource.TestCheckResourceAttr(resourceName, "key", "tac"),
					resource.TestCheckResourceAttr(resourceName, "name", "Temporary Access Code"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.type", "TAC"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.length", "12"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.min_ttl", "15"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.default_ttl", "45"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.max_ttl", "120"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.multi_use_allowed", "false"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.complexity.letters", "true"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.complexity.numbers", "true"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.complexity.special_characters", "false"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test TAC-3"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.length", "16"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.min_ttl", "20"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.default_ttl", "30"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.max_ttl", "60"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.multi_use_allowed", "true"),
					resource.TestCheckResourceAttr(resourceName, "provider_tac.configuration.complexity.special_characters", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// type and description are Optional-only and unset in config; import maps the API values.
				ImportStateVerifyIgnore: []string{"type", "description"},
			},
		},
	})
}

func importStateIdForAuthenticatorTac(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["data.okta_authenticator.tac"]
	if !ok {
		return "", fmt.Errorf("data.okta_authenticator.tac not found in state")
	}
	return rs.Primary.ID, nil
}
