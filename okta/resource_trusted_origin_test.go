package okta

import (
  "fmt"
  "testing"
  "github.com/hashicorp/terraform/helper/acctest"
  "github.com/hashicorp/terraform/helper/resource"
  "github.com/hashicorp/terraform/terraform"
)

func TestAccTrustedOrigin(t *testing.T) {

  // generate a random name for each widget test run, to avoid
  // collisions from multiple concurrent tests.
  // the acctest package includes many helpers such as RandStringFromCharSet
  // See https://godoc.org/github.com/hashicorp/terraform/helper/acctest
  rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

  resource.Test(t, resource.TestCase{
    PreCheck:  func() { testAccPreCheck(t) },
    Providers: testAccProviders,
    CheckDestroy: testAccCheckTrustedOriginDestroy,
    Steps: []resource.TestStep{
      {
        Config: testAccTrustedOriginCreate(rName),
      },
      {
        Config: testAccTrustedOriginUpdate(rName),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("okta_trusted_origin.test_"+rName, "origin", "https://example2-"+rName+".com"),
        ),
      },
    },
  })
}

func testAccTrustedOriginCreate(name string) string {
  return fmt.Sprintf(`
resource "okta_trusted_origin" "test_%s" {
  name = "test-%s"
  origin = "https://example-%s.com"
  scopes = ["CORS"]
}`, name, name, name)
}

func testAccTrustedOriginUpdate(name string) string {
  return fmt.Sprintf(`
resource "okta_trusted_origin" "test_%s" {
  name = "test-%s"
  active = false
  origin = "https://example2-%s.com"
  scopes = ["CORS", "REDIRECT"]
}`, name, name, name)
}

func testAccCheckTrustedOriginDestroy(s *terraform.State) error {
  client := testAccProvider.Meta().(*Config).oktaClient

  for _, r := range s.RootModule().Resources {
    if _, _, err := client.TrustedOrigins.GetTrustedOrigin(r.Primary.ID); err != nil {
      if client.OktaErrorCode == "E0000007" {
        continue
      }
      return fmt.Errorf("[ERROR] Error Getting Trusted Origin in Okta: %v", err)
    }
    return fmt.Errorf("Trusted Origin still exists")
  }

  return nil
}
