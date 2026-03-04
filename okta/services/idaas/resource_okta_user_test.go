package idaas_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccResourceOktaUser_customProfileAttributes(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	config := mgr.GetFixtures("custom_attributes.tf", t)
	arrayAttrConfig := mgr.GetFixtures("custom_attributes_array.tf", t)
	ignoreConfig := mgr.GetFixtures("custom_attributes_to_ignore.tf", t)
	updatedConfig := mgr.GetFixtures("remove_custom_attributes.tf", t)
	importConfig := mgr.GetFixtures("import.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUser)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
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

func TestAccResourceOktaUser_invalidCustomProfileAttribute(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testOktaUserConfigInvalidCustomProfileAttribute(mgr.SeedStr()),
				ExpectError: regexp.MustCompile("Api validation failed: newUser"),
			},
		},
	})
}

func TestAccResourceOktaUser_updateAllAttributes(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	config := mgr.GetFixtures("staged.tf", t)
	updatedConfig := mgr.GetFixtures("all_attributes.tf", t)
	minimalConfig := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUser)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
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

func TestAccResourceOktaUser_updateCredentials(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	config := mgr.GetFixtures("basic_with_credentials.tf", t)
	minimalConfigWithCredentials := mgr.GetFixtures("basic_with_credentials_updated.tf", t)
	minimalConfigWithCredentialsOldPassword := mgr.GetFixtures("basic_with_credentials_updated_old_password.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUser)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
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

func TestAccResourceOktaUser_statusDeprovisioned(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	statusChanged := mgr.GetFixtures("deprovisioned.tf", t)
	config := mgr.GetFixtures("staged.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUser)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
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
					resource.TestCheckResourceAttr(resourceName, "status", idaas.UserStatusDeprovisioned),
				),
			},
		},
	})
}

func TestAccResourceOktaUser_hashed_password_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	config := mgr.GetFixtures("password_hash.tf", t)
	configUpdated := mgr.GetFixtures("password_hash_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUser)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.UserStatusStaged),
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
					resource.TestCheckResourceAttr(resourceName, "status", idaas.UserStatusStaged),
					resource.TestCheckResourceAttr(resourceName, "password_hash.0.algorithm", "BCRYPT"),
				),
			},
		},
	})
}

func TestAccResourceOktaUser_updateDeprovisioned(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	config := mgr.GetFixtures("deprovisioned.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
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

func TestAccResourceOktaUser_loginUpdates(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedLogin := mgr.GetFixtures("login_changed.tf", t)

	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUser)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)
	updatedEmail := fmt.Sprintf("testAccUpdated-%d@example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
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
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
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

// TestAccResourceOktaUser_issue_1216_Suppress403Errors
// https://github.com/okta/terraform-provider-okta/issues/1216 When this test
// runs with an API token of Org Admin (not Super Admin) the resource will fail
// when the admin roles are gathered.
func TestAccResourceOktaUser_issue_1216_Suppress403Errors(t *testing.T) {
	if !orgAdminOnlyTest(t) {
		return
	}
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	config := `
resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}`
	config = mgr.ConfigReplace(config)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUser)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
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

func TestAccResourceOktaUser_withUserType(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	config := mgr.GetFixtures("user_with_type.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUser)
	userTypeResourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUserType)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "type.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "type.0.id", userTypeResourceName, "id"),
				),
			},
		},
	})
}

func TestAccResourceOktaUser_addUserTypeLater(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSUser, t.Name())
	configWithoutType := mgr.GetFixtures("basic.tf", t)
	configWithType := mgr.GetFixtures("user_with_type.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUser)
	userTypeResourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUserType)
	email := fmt.Sprintf("testAcc-%d@example.com", mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithoutType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "type.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "type.0.id"),
				),
			},
			{
				Config: configWithType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "first_name", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Smith"),
					resource.TestCheckResourceAttr(resourceName, "login", email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
					resource.TestCheckResourceAttr(resourceName, "type.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "type.0.id", userTypeResourceName, "id"),
				),
			},
		},
	})
}
