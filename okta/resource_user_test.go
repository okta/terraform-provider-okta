package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOktaUserNew(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := "okta_user.test_acc_" + rName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUsersDestroy,
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
}
`, r, r, r)
}

func testOktaUserConfig_updated(r string) string {
  return fmt.Sprintf(`
resource "okta_user" "test_acc_%s" {
  admin_roles      = ["APP_ADMIN"]
  first_name       = "TestAcc"
  last_name        = "%s"
  login            = "test-acc-%s@testing.com"
  honorific_prefix = "Dr."
  honorific_suffix = "Jr."
}
`, r, r, r)
}
