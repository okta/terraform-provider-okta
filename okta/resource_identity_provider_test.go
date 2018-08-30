package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIdentityProvider_create(t *testing.T) {

	// generate a random name for each widget test run, to avoid
	// collisions from multiple concurrent tests.
	// the acctest package includes many helpers such as RandStringFromCharSet
	// See https://godoc.org/github.com/hashicorp/terraform/helper/acctest
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderCreate(rName),
			},
			{
				Config: testAccIdentityProviderUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_identity_provider.test-"+rName, "client_id", "newClientID"),
				),
			},
		},
	})
}

func testOktaIdentityProviderExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[ERROR] Resource Not found: %s", name)
		}

		IdpID, hasID := rs.Primary.Attributes["id"]
		if !hasID {
			return fmt.Errorf("[ERROR] No id found in state for Identity Provider")
		}
		IdpName, hasName := rs.Primary.Attributes["name"]
		if !hasName {
			return fmt.Errorf("[ERROR] No name found in state for Identity Provider")
		}

		err := testPolicyExists(true, IdpID, IdpName)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func testAccIdentityProviderCreate(name string) string {
	return fmt.Sprintf(`
resource "okta_identity_provider" "test-%s" {
  type = "GOOGLE"
  name = "%s"
  protocol_type   = "OIDC"
  protocol_scopes = ["profile", "email", "openid"]
  client_id = "2780nfqgi7gioq39asdg"
  client_secret = "134t98higlhalkgjhakj"
  policy_provisioning_group_assignments = ["banana"]
}`, name, name)
}

func testAccIdentityProviderUpdate(name string) string {
	return fmt.Sprintf(`
resource "okta_identity_provider" "test-%s" {
  type = "GOOGLE"
  name = "%s"
  active = false
  protocol_type   = "OIDC"
  protocol_scopes = ["profile", "email"]
  client_id = "newClientID"
  client_secret = "134t98higlhalkgjhakj"
}`, name, name)
}
