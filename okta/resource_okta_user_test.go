package okta

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOktaUser_customProfileAttributes(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	config := mgr.GetFixtures("custom_attributes.tf", t)
	arrayAttrConfig := mgr.GetFixtures("custom_attributes_array.tf", t)
	ignoreConfig := mgr.GetFixtures("custom_attributes_to_ignore.tf", t)
	updatedConfig := mgr.GetFixtures("remove_custom_attributes.tf", t)
	importConfig := mgr.GetFixtures("import.tf", t)
	resourceName := fmt.Sprintf("%s.test", user)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config:  config,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", "{\"customAttribute123\":\"testing-custom-attribute\"}"),
				),
			},
			{
				Config:  arrayAttrConfig,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", "{\"array123\":[\"test\"],\"number123\":1}"),
				),
			},
			{
				Config:  ignoreConfig,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", "{\"customAttribute1234\":\"testing-custom-attribute\"}"), // Note: "customAttribute123" is ignored and should not be present
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
			{
				Config: importConfig,
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) (err error) {
					if len(s) != 1 {
						err = errors.New("failed to import into resource into state")
						return
					}

					id := s[0].Attributes["id"]

					if strings.Contains(id, "@") {
						err = fmt.Errorf("user resource id incorrectly set, %s", id)
					}
					return
				},
			},
		},
	})
}

func TestAccOktaUser_invalidCustomProfileAttribute(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testOktaUserConfigInvalidCustomProfileAttribute(rName),
				ExpectError: regexp.MustCompile("Api validation failed: newUser"),
			},
		},
	})
}

func TestAccOktaUser_updateAllAttributes(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	config := mgr.GetFixtures("staged.tf", t)
	updatedConfig := mgr.GetFixtures("all_attributes.tf", t)
	minimalConfig := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", user)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
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
					resource.TestCheckResourceAttr(resourceName, "second_email", fmt.Sprintf("test2-%d@example.com", mgr.Seed)),
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

func TestAccOktaUser_updateCredentials(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	config := mgr.GetFixtures("basic_with_credentials.tf", t)
	minimalConfigWithCredentials := mgr.GetFixtures("basic_with_credentials_updated.tf", t)
	minimalConfigWithCredentialsOldPassword := mgr.GetFixtures("basic_with_credentials_updated_old_password.tf", t)
	resourceName := fmt.Sprintf("%s.test", user)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "password", "Abcd1234"),
					resource.TestCheckResourceAttr(resourceName, "recovery_answer", "Forty Two"),
				),
			},
			{
				Config: minimalConfigWithCredentials,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "password", "SuperSecret007"),
					resource.TestCheckResourceAttr(resourceName, "recovery_answer", "Asterisk"),
				),
			},
			{
				Config: minimalConfigWithCredentialsOldPassword,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "password", "Super#Secret@007"),
					resource.TestCheckResourceAttr(resourceName, "old_password", "SuperSecret007"),
					resource.TestCheckResourceAttr(resourceName, "recovery_answer", "0010"),
				),
			},
		},
	})
}

func TestAccOktaUser_statusDeprovisioned(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	statusChanged := mgr.GetFixtures("deprovisioned.tf", t)
	config := mgr.GetFixtures("staged.tf", t)
	resourceName := fmt.Sprintf("%s.test", user)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
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
					resource.TestCheckResourceAttr(resourceName, "status", userStatusDeprovisioned),
				),
			},
		},
	})
}

func TestAccOktaUserHashedPassword(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	config := mgr.GetFixtures("password_hash.tf", t)
	configUpdated := mgr.GetFixtures("password_hash_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", user)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "status", userStatusStaged),
					resource.TestCheckResourceAttr(resourceName, "password_hash.0.algorithm", "SHA-512"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "status", userStatusStaged),
					resource.TestCheckResourceAttr(resourceName, "password_hash.0.algorithm", "BCRYPT"),
				),
			},
		},
	})
}

func TestAccOktaUser_updateDeprovisioned(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	config := mgr.GetFixtures("deprovisioned.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config:      testOktaUserConfigUpdateDeprovisioned(strconv.Itoa(mgr.Seed)),
				ExpectError: regexp.MustCompile(".*Only the status of a DEPROVISIONED user can be updated, we detected other change"),
			},
		},
	})
}

func TestAccOktaUser_loginUpdates(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedLogin := mgr.GetFixtures("login_changed.tf", t)

	resourceName := fmt.Sprintf("%s.test", user)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)
	updatedEmail := fmt.Sprintf("testAccUpdated-%d@example.com", mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(resourceName, "email", email),
			},
			{
				Config: updatedLogin,
				Check:  resource.TestCheckResourceAttr(resourceName, "email", updatedEmail),
			},
		},
	})
}

func checkUserDestroy(s *terraform.State) error {
	if isVCRPlayMode() {
		return nil
	}
	client := oktaClientForTest()
	for _, r := range s.RootModule().Resources {
		if _, resp, err := client.User.GetUser(context.Background(), r.Primary.ID); err != nil {
			if resp != nil && resp.Response.StatusCode == http.StatusNotFound {
				continue
			}
			return fmt.Errorf("[ERROR] Error Getting User in Okta: %v", err)
		}
		return fmt.Errorf("user still exists")
	}
	return nil
}

func testOktaUserConfigInvalidCustomProfileAttribute(r string) string {
	return fmt.Sprintf(`
resource okta_user "test" {
  first_name  = "TestAcc"
  last_name   = "%[1]s"
  login       = "testAcc-%[1]s@example.com"
  email       = "testAcc-%[1]s@example.com"

  custom_profile_attributes = <<JSON
  {
    "notValid": "this-isnt-valid"
  }
  JSON
}
`, r)
}

func testOktaUserConfigUpdateDeprovisioned(r string) string {
	return fmt.Sprintf(`
resource okta_user "test" {
  first_name  = "TestAcc"
  last_name   = "%[1]s"
  login       = "testAcc-%[1]s@example.com"
  status      = "DEPROVISIONED"
  email       = "hello@example.com"
}
`, r)
}

// TestIssue1216Suppress403Errors
// https://github.com/okta/terraform-provider-okta/issues/1216 When this test
// runs with an API token of Org Admin (not Super Admin) the resource will fail
// when the admin roles are gathered.
func TestIssue1216Suppress403Errors(t *testing.T) {
	if !orgAdminOnlyTest(t) {
		return
	}
	mgr := newFixtureManager(user, t.Name())
	config := `
resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}`
	config = mgr.ConfigReplace(config)
	resourceName := fmt.Sprintf("%s.test", user)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "admin_roles.#", "0"),
				),
			},
		},
	})
}
