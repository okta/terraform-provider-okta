package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	//"github.com/hashicorp/terraform/terraform"
)

// Witiz1932@teleworm.us is a fake email address created at fakemailgenerator.com
// view inbox: http://www.fakemailgenerator.com/inbox/teleworm.us/witiz1932/

func TestAccOktaGroups(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaGroups_create(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
			},
		},
	})
}

func testOktaGroups_create(rInt int) string {
	return fmt.Sprintf(`
resource "okta_group" "test-%d" {
  name = "testgroup-%d"
  description = "testing, testing, 1 2..."
}
`, rInt, rInt)
}
