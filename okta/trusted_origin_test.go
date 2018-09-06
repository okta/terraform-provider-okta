package okta

import (
  "fmt"
  "testing"

  "github.com/hashicorp/terraform/helper/acctest"
  "github.com/hashicorp/terraform/helper/resource"
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
    Steps: []resource.TestStep{
      {
        Config: testAccTrustedOriginCreate(rName),
      },
      {
        Config: testAccTrustedOriginUpdate(rName),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("okta_trusted_origin.test-"+rName, "origin", "https://example2-"+rName+".com"),
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
  origin = "https://example2-%s.com"
  scopes = ["CORS"]
}`, name, name, name)
}


