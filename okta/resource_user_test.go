package okta

import (
	"fmt"
  "strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
  "github.com/hashicorp/terraform/terraform"
)

func TestAccOktaUserNew(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := "okta_user.test_acc_" + rName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testOktaUserConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", rName),
					resource.TestCheckResourceAttr(resourceName, "login", "test-acc-"+rName+"@testing.com"),
				),
			},
      {
        Config: testOktaUserConfig_updated(rName),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
          resource.TestCheckResourceAttr(resourceName, "last_name", rName),
          resource.TestCheckResourceAttr(resourceName, "login", "test-acc-"+rName+"@testing.com"),
        ),
      },
		},
	})
}

func testOktaUserConfig(r string) string {
	return fmt.Sprintf(`
resource "okta_user" "test_acc_%s" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "%s"
  login       = "test-acc-%s@testing.com"
  status      = "STAGED"
}
`, r, r, r)
}

func testOktaUserConfig_updated(r string) string {
  return fmt.Sprintf(`
resource "okta_user" "test_acc_%s" {
  admin_roles      = ["ORG_ADMIN"]
  first_name       = "TestAcc"
  last_name        = "%s"
  login            = "test-acc-%s@testing.com"
  honorific_prefix = "Dr."
  honorific_suffix = "Jr."
}
`, r, r, r)
}

func testAccCheckUserDestroy(s *terraform.State) error {
  client := testAccProvider.Meta().(*Config).oktaClient

  for _, r := range s.RootModule().Resources {
    if _, resp, err := client.User.GetUser(r.Primary.ID, nil); err != nil {
      if strings.Contains(resp.Response.Status, "404") {
        continue
      }
      return fmt.Errorf("[ERROR] Error Getting User in Okta: %v", err)
    }
    return fmt.Errorf("User still exists")
  }

  return nil
}
