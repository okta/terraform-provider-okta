package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
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
				// use a dynamic configuration with the random name from above
				Config: testAccIdentityProviderCreate(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("identity_provider.test-"+rName, "name", rName),
				),
			},
		},
	})
}

// testAccExampleResource returns an configuration for an Example Widget with the provided name
func testAccIdentityProviderCreate(name string) string {
	return fmt.Sprintf(`
resource "okta_identity_provider" "foo" {
  type = "GOOGLE"
  name = "%s"
}`, name)
}
