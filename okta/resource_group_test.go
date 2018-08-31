package okta

import (
	"fmt"
	"testing"
	"strconv"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	//"github.com/hashicorp/terraform/terraform"
)

// Witiz1932@teleworm.us is a fake email address created at fakemailgenerator.com
// view inbox: http://www.fakemailgenerator.com/inbox/teleworm.us/witiz1932/

func TestAccOktaGroupsCreate(t *testing.T) {
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOktaGroups_create(ri),
			},
			{
				Config: testOktaGroups_update(ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_group.test-"+strconv.Itoa(ri), "name", "testgroupdifferent")),
			},
		},
	})
}


func testOktaGroups_create(rInt int) string {
	return fmt.Sprintf(`
resource "okta_group" "test-%d" {
  name = "testgroup-%d"
  description = "testing, testing"
}
`, rInt, rInt)
}

func testOktaGroups_update(rInt int) string {
	return fmt.Sprintf(`
resource "okta_group" "test-%d" {
  name = "testgroupdifferent"
  description = "testing, testing"
}
`, rInt)
}

