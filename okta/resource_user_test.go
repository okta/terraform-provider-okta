package okta

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOktaUser_emailError(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testOktaUserConfig_emailError(rName),
				ExpectError: regexp.MustCompile("login field not a valid email address"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccOktaUser_updateDeprovisioned(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testOktaUserConfig_deprovisioned(rName),
			},
			{
				Config:      testOktaUserConfig_updateDeprovisioned(rName),
				ExpectError: regexp.MustCompile("Cannot update a DEPROVISIONED user"),
			},
		},
	})
}

func TestAccOktaUser_validRole(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testOktaUserConfig_validRole(rName),
				ExpectError: regexp.MustCompile("GROUP_ADMIN is not a valid Okta role"),
			},
		},
	})
}

func TestAccOktaUser_updateAllAttributes(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "email", "test1-"+rName+"@testing.com"),
					resource.TestCheckResourceAttr(resourceName, "city", "New York"),
					resource.TestCheckResourceAttr(resourceName, "cost_center", "10"),
					resource.TestCheckResourceAttr(resourceName, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceName, "department", "IT"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Dr. TestAcc "+rName),
					resource.TestCheckResourceAttr(resourceName, "division", "Acquisitions"),
					resource.TestCheckResourceAttr(resourceName, "employee_number", "111111"),
					resource.TestCheckResourceAttr(resourceName, "honorific_prefix", "Dr."),
					resource.TestCheckResourceAttr(resourceName, "honorific_suffix", "Jr."),
					resource.TestCheckResourceAttr(resourceName, "locale", "en_US"),
					resource.TestCheckResourceAttr(resourceName, "manager", "Jimbo"),
					resource.TestCheckResourceAttr(resourceName, "manager_id", "222222"),
					resource.TestCheckResourceAttr(resourceName, "middle_name", "John"),
					resource.TestCheckResourceAttr(resourceName, "mobile_phone", "1112223333"),
					resource.TestCheckResourceAttr(resourceName, "nick_name", "Johnny"),
					resource.TestCheckResourceAttr(resourceName, "organization", "Testing Inc."),
					resource.TestCheckResourceAttr(resourceName, "postal_address", "1234 Testing St."),
					resource.TestCheckResourceAttr(resourceName, "preferred_language", "en-us"),
					resource.TestCheckResourceAttr(resourceName, "primary_phone", "4445556666"),
					resource.TestCheckResourceAttr(resourceName, "profile_url", "http://www.example.com/profile"),
					resource.TestCheckResourceAttr(resourceName, "second_email", "test2-"+rName+"@testing.com"),
					resource.TestCheckResourceAttr(resourceName, "state", "NY"),
					resource.TestCheckResourceAttr(resourceName, "street_address", "5678 Testing Ave."),
					resource.TestCheckResourceAttr(resourceName, "timezone", "America/New_York"),
					resource.TestCheckResourceAttr(resourceName, "title", "Director"),
					resource.TestCheckResourceAttr(resourceName, "user_type", "Employee"),
					resource.TestCheckResourceAttr(resourceName, "zip_code", "11111"),
				),
			},
		},
	})
}

func testAccCheckUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).oktaClient

	for _, r := range s.RootModule().Resources {
		if _, resp, err := client.User.GetUser(r.Primary.ID); err != nil {
			if strings.Contains(resp.Response.Status, "404") {
				continue
			}
			return fmt.Errorf("[ERROR] Error Getting User in Okta: %v", err)
		}
		return fmt.Errorf("User still exists")
	}

	return nil
}

func testOktaUserConfig_emailError(r string) string {
	return fmt.Sprintf(`
resource "okta_user" "test_%s" {
  firstname = "testAcc"
  lastname  = "%s"
  login     = "notavalidemail"
  role      = ["SUPER_ADMIN"]
}
`, r, r)
}

func testOktaUserConfig_deprovisioned(r string) string {
	return fmt.Sprintf(`
resource "okta_user" "test_acc_%s" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "%s"
  login       = "test-acc-%s@testing.com"
  status      = "DEPROVISIONED"
}
`, r, r, r)
}

func testOktaUserConfig_updateDeprovisioned(r string) string {
	return fmt.Sprintf(`
resource "okta_user" "test_acc_%s" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "%s"
  login       = "test-acc-%s@testing.com"
  status      = "DEPROVISIONED"
  email       = "hello@testing.com"
}
`, r, r, r)
}

func testOktaUserConfig_validRole(r string) string {
	return fmt.Sprintf(`
resource "okta_user" "test_acc_%s" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN", "GROUP_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "%s"
  login       = "test-acc-%s@testing.com"
}
`, r, r, r)
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
  admin_roles        = ["ORG_ADMIN"]
  first_name         = "TestAcc"
  last_name          = "%s"
  login              = "test-acc-%s@testing.com"
  email              = "test1-%s@testing.com"
  city               = "New York"
  cost_center        = "10"
  country_code       = "US"
  department         = "IT"
  display_name       = "Dr. TestAcc %s"
  division           = "Acquisitions"
  employee_number    = "111111"
  honorific_prefix   = "Dr."
  honorific_suffix   = "Jr."
  locale             = "en_US"
  manager            = "Jimbo"
  manager_id         = "222222"
  middle_name        = "John"
  mobile_phone       = "1112223333"
  nick_name          = "Johnny"
  organization       = "Testing Inc."
  postal_address     = "1234 Testing St."
  preferred_language = "en-us"
  primary_phone      = "4445556666"
  profile_url        = "http://www.example.com/profile"
  second_email       = "test2-%s@testing.com"
  state              = "NY"
  street_address     = "5678 Testing Ave."
  timezone           = "America/New_York"
  title              = "Director"
  user_type          = "Employee"
  zip_code           = "11111"
}
`, r, r, r, r, r, r)
}
