package okta

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOktaUser_customProfileAttributes(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaUser)
	config := mgr.GetFixtures("okta_user_custom_attributes.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_user_remove_custom_attributes.tf", ri, t)
	resourceName := buildResourceFQN(oktaUser, ri)
	email := fmt.Sprintf("test-acc-%d@testing.com", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes.customAttribute123", "testing-custom-attribute"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
				),
			},
		},
	})
}

func TestAccOktaUser_groupMembership(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaUser)
	config := mgr.GetFixtures("okta_user_group_assigned.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_user_group_unassigned.tf", ri, t)
	resourceName := buildResourceFQN(oktaUser, ri)
	email := fmt.Sprintf("test-acc-%d@testing.com", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "group_memberships.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "group_memberships.#", "0"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "group_memberships.#", "1"),
				),
			},
		},
	})
}

func TestAccOktaUser_invalidCustomProfileAttribute(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testOktaUserConfig_invalidCustomProfileAttribute(rName),
				ExpectError: regexp.MustCompile("Api validation failed: newUser"),
			},
		},
	})
}

func TestAccOktaUser_updateAllAttributes(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaUser)
	config := mgr.GetFixtures("okta_user_staged.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_user_all_attributes.tf", ri, t)
	minimalConfig := mgr.GetFixtures("okta_user.tf", ri, t)
	resourceName := buildResourceFQN(oktaUser, ri)
	email := fmt.Sprintf("test-acc-%d@testing.com", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "2"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", fmt.Sprintf("test1-%d@testing.com", ri)),
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "city", "New York"),
					resource.TestCheckResourceAttr(resourceName, "cost_center", "10"),
					resource.TestCheckResourceAttr(resourceName, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceName, "department", "IT"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Dr. TestAcc Smith"),
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
					resource.TestCheckResourceAttr(resourceName, "second_email", fmt.Sprintf("test2-%d@testing.com", ri)),
					resource.TestCheckResourceAttr(resourceName, "state", "NY"),
					resource.TestCheckResourceAttr(resourceName, "street_address", "5678 Testing Ave."),
					resource.TestCheckResourceAttr(resourceName, "timezone", "America/New_York"),
					resource.TestCheckResourceAttr(resourceName, "title", "Director"),
					resource.TestCheckResourceAttr(resourceName, "user_type", "Employee"),
					resource.TestCheckResourceAttr(resourceName, "zip_code", "11111"),
				),
			},
			{
				Config: minimalConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
				),
			},
		},
	})
}

func TestAccOktaUser_statusDeprovisioned(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaUser)
	statusChanged := mgr.GetFixtures("okta_user_deprovisioned.tf", ri, t)
	config := mgr.GetFixtures("okta_user_staged.tf", ri, t)
	resourceName := buildResourceFQN(oktaUser, ri)
	email := fmt.Sprintf("test-acc-%d@testing.com", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: statusChanged,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "status", "DEPROVISIONED"),
				),
			},
		},
	})
}

func TestAccOktaUser_updateDeprovisioned(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oktaUser)
	config := mgr.GetFixtures("okta_user_deprovisioned.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config:      testOktaUserConfig_updateDeprovisioned(strconv.Itoa(ri)),
				ExpectError: regexp.MustCompile(".*Only the status of a DEPROVISIONED user can be updated, we detected other change"),
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

func testOktaUserConfig_invalidCustomProfileAttribute(r string) string {
	return fmt.Sprintf(`
resource okta_user "testAcc_%[1]s" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "%[1]s"
  login       = "test-acc-%[1]s@testing.com"
  email       = "test-acc-%[1]s@testing.com"

  custom_profile_attributes {
    notValid = "this-isnt-valid"
  }
}
`, r)
}

func testOktaUserConfig_updateDeprovisioned(r string) string {
	return fmt.Sprintf(`
resource okta_user "testAcc_%[1]s" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "%[1]s"
  login       = "test-acc-%[1]s@testing.com"
  status      = "DEPROVISIONED"
  email       = "hello@testing.com"
}
`, r)
}

func testOktaUserConfig_validRole(r string) string {
	return fmt.Sprintf(`
resource okta_user "testAcc_%[1]s" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN", "GROUP_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "Smith"
  login       = "test-acc-%[1]s@testing.com"
  email       = "test-acc-%[1]s@testing.com"
}
`, r)
}
