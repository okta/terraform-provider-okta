package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaUser_read(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	baseConfig := mgr.GetFixtures("datasource.tf", t)
	createUserConfig := mgr.GetFixtures("datasource_create_user.tf", t)

	// NOTE: eliminated previous flapping issues when delay_read_seconds was added to okta_user
	// TF_ACC=1 go test -tags unit -mod=readonly -test.v -run ^TestAccOktaDataSourceUser_read$
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: createUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),
				),
			},
			{
				Config: baseConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),

					resource.TestCheckResourceAttrSet("data.okta_user.first_and_last", "id"),
					resource.TestCheckResourceAttr("data.okta_user.first_and_last", "first_name", "TestAcc"),
					resource.TestCheckResourceAttr("data.okta_user.first_and_last", "last_name", "Smith"),

					resource.TestCheckResourceAttr("data.okta_user.read_by_id", "first_name", "TestAcc"),
					resource.TestCheckResourceAttr("data.okta_user.read_by_id", "last_name", "Smith"),

					resource.TestCheckResourceAttrSet("data.okta_user.compound_search", "id"),
					resource.TestCheckResourceAttrSet("data.okta_user.compound_search", "compound_search_operator"),
					resource.TestCheckResourceAttr("data.okta_user.compound_search", "compound_search_operator", "or"),
					resource.TestCheckResourceAttr("data.okta_user.compound_search", "first_name", "Some"),
					resource.TestCheckResourceAttr("data.okta_user.compound_search", "last_name", "One"),

					resource.TestCheckResourceAttrSet("data.okta_user.expression_search", "id"),
					resource.TestCheckResourceAttr("data.okta_user.expression_search", "first_name", "Some"),
					resource.TestCheckResourceAttr("data.okta_user.expression_search", "last_name", "One"),
					resource.TestCheckResourceAttr("data.okta_user.expression_search", "custom_profile_attributes", `{"array123":["cool","feature"]}`),
				),
			},
		},
	})
}

// TestAccDataSourceOktaUser_SkipAdminRoles pertains to https://github.com/okta/terraform-provider-okta/pull/1137 and https://github.com/okta/terraform-provider-okta/issues/1014
func TestAccDataSourceOktaUser_SkipAdminRoles(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(testOktaUserRolesGroupsConfig(false, true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.okta_user.test", "admin_roles.#"),          // skipped
					resource.TestCheckNoResourceAttr("data.okta_user.test", "roles.#"),                // skipped
					resource.TestCheckResourceAttr("data.okta_user.test", "group_memberships.#", "2"), // Everyone, A Group
				),
			},
		},
	})
}

// TestAccDataSourceOktaUser_SkipGroups pertains to https://github.com/okta/terraform-provider-okta/pull/1137 and https://github.com/okta/terraform-provider-okta/issues/1014
func TestAccDataSourceOktaUser_SkipGroups(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(testOktaUserRolesGroupsConfig(true, false)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_user.test", "admin_roles.#", "2"),       // SUPER_ADMIN, APP_ADMIN
					resource.TestCheckResourceAttr("data.okta_user.test", "roles.#", "2"),             // SUPER_ADMIN, APP_ADMIN
					resource.TestCheckResourceAttr("data.okta_user.test", "group_memberships.#", "0"), // skipped
				),
			},
		},
	})
}

// TestAccDataSourceOktaUser_SkipGroupsSkipRoles pertains to https://github.com/okta/terraform-provider-okta/pull/1137 and https://github.com/okta/terraform-provider-okta/issues/1014
func TestAccDataSourceOktaUser_SkipGroupsSkipRoles(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(testOktaUserRolesGroupsConfig(true, true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_user.test", "admin_roles.#", "0"),       // skipped
					resource.TestCheckResourceAttr("data.okta_user.test", "roles.#", "0"),             // skipped
					resource.TestCheckResourceAttr("data.okta_user.test", "group_memberships.#", "0"), // skipped
				),
			},
		},
	})
}

// TestAccDataSourceOktaUser_NoSkips pertains to https://github.com/okta/terraform-provider-okta/pull/1137 and https://github.com/okta/terraform-provider-okta/issues/1014
func TestAccDataSourceOktaUser_NoSkips(t *testing.T) {
	mgr := newFixtureManager(user, t.Name())
	allAdminRolesRegexp, _ := regexp.Compile("APP_ADMIN, SUPER_ADMIN")
	allGroupMembershipsRegexp, _ := regexp.Compile("00g[a-z,A-Z,0-9]{17}, 00g[a-z,A-Z,0-9]{17}")
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(testOktaUserRolesGroupsConfig(false, false)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_user.test", "admin_roles.#", "2"),       // SUPER_ADMIN, APP_ADMIN
					resource.TestCheckResourceAttr("data.okta_user.test", "roles.#", "2"),             // SUPER_ADMIN, APP_ADMIN
					resource.TestCheckResourceAttr("data.okta_user.test", "group_memberships.#", "2"), // Everyone, A Group
					resource.TestMatchOutput("output_admin_roles", allAdminRolesRegexp),
					resource.TestMatchOutput("output_group_memberships", allGroupMembershipsRegexp),
				),
			},
		},
	})
}

func testOktaUserRolesGroupsConfig(skipGroups, skipRoles bool) string {
	prepend := `

resource "okta_group" "testAcc-replace_with_uuid" {
  name        = "testAcc-replace_with_uuid"
  description = "A Group"
}
resource "okta_user" "testAcc-replace_with_uuid" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}
resource "okta_user_admin_roles" "test" {
  user_id     = okta_user.testAcc-replace_with_uuid.id
  admin_roles = [
    "SUPER_ADMIN",
    "APP_ADMIN",
  ]
}
resource "okta_user_group_memberships" "test" {
  user_id = okta_user.testAcc-replace_with_uuid.id
  groups = [
    okta_group.testAcc-replace_with_uuid.id,
  ]
}
data "okta_user" "test" {
  user_id = okta_user.testAcc-replace_with_uuid.id
`

	var clause string
	if skipGroups {
		clause = "  skip_groups = true"
	}
	if skipRoles {
		clause = fmt.Sprintf("%s\n  skip_roles = true\n", clause)
	}

	append := `
  depends_on = [
    okta_user_admin_roles.test,
    okta_user_admin_roles.test,
    okta_user_group_memberships.test,
  ]
}
output "output_admin_roles" {
  value = join(", ", data.okta_user.test.admin_roles)
}
output "output_group_memberships" {
  value = join(", ", data.okta_user.test.group_memberships)
}
`

	return fmt.Sprintf("%s%s%s", prepend, clause, append)
}
