package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOktaFactor(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaFactor(ri)
	updatedConfig := testOktaFactorInactive(ri)
	resourceName := buildResourceFQN(factor, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "provider_id", "google_otp"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "provider_id", "google_otp"),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
				),
			},
		},
	})
}

func testOktaFactor(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
	provider_id  = "google_otp"
	status   	 = "ACTIVE"
}
`, factor, name)
}

func testOktaFactorInactive(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
	provider_id  = "google_otp"
	status   	 = "INACTIVE"
}
`, factor, name)
}
