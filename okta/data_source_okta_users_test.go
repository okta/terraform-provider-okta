package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	allAdminRolesRegexp, _       = regexp.Compile("APP_ADMIN, SUPER_ADMIN")
	allGroupMembershipsRegexp, _ = regexp.Compile("00g[a-z,A-Z,0-9]{17}, 00g[a-z,A-Z,0-9]{17}")
)

func TestAccOktaDataSourceUsers_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(users)
	users := mgr.GetFixtures("users.tf", ri, t)
	config := mgr.GetFixtures("basic.tf", ri, t)
	dataSource := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				// Ensure users are created
				Config: users,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test1", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test2", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test3", "id"),
				),
			},
			{
				// Ensure data source props are set
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_users.test", "users.#"),
				),
			},
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_users.compound_search", "compound_search_operator"),
					resource.TestCheckResourceAttr("data.okta_users.compound_search", "compound_search_operator", "and"),
					resource.TestCheckResourceAttrSet("data.okta_users.compound_search", "users.#"),
					resource.TestCheckResourceAttr("data.okta_users.compound_search", "users.#", "1"),
				),
			},
		},
	})
}

func TestAccOktaDataSourceUsers_readWithGroupId(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(users)
	users := mgr.GetFixtures("users_with_group.tf", ri, t)
	config := mgr.GetFixtures("group.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				// Ensure user and group are created
				Config: users,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test1", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test2", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test3", "id"),
					resource.TestCheckResourceAttrSet("okta_group.test", "id"),
					resource.TestCheckResourceAttr("okta_group_memberships.test", "users.#", "2"),
				),
			},
			{
				// Ensure data source props are set
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_users.test", "users.#", "2"),
				),
			},
		},
	})
}

func TestAccOktaDataSourceUsers_readWithGroupIdIncludingGroups(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(users)
	users := mgr.GetFixtures("users_with_group.tf", ri, t)
	config := mgr.GetFixtures("group_with_groups.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				// Ensure user and group are created
				Config: users,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test1", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test2", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test3", "id"),
					resource.TestCheckResourceAttrSet("okta_group.test", "id"),
					resource.TestCheckResourceAttr("okta_group_memberships.test", "users.#", "2"),
				),
			},
			{
				// Ensure data source props are set
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_users.test", "users.#", "2"),
				),
			},
		},
	})
}

// TestAccDataSourceOktaUsers_IncludeNone pertains to https://github.com/okta/terraform-provider-okta/pull/1137 and https://github.com/okta/terraform-provider-okta/issues/1014
func TestAccDataSourceOktaUsers_IncludeNone(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: testOktaUsersRolesGroupsConfig(false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.okta_users.test", "users.0.admin_roles.#", "0"),       // skipped
					resource.TestCheckResourceAttr("data.okta_users.test", "users.0.group_memberships.#", "0"), // skipped
					resource.TestCheckOutput("output_admin_roles", ""),
					resource.TestCheckOutput("output_group_memberships", ""),
				),
			},
		},
	})
}

// TestAccDataSourceOktaUsers_IncludeGroups pertains to https://github.com/okta/terraform-provider-okta/pull/1137 and https://github.com/okta/terraform-provider-okta/issues/1014
func TestAccDataSourceOktaUsers_IncludeGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: testOktaUsersRolesGroupsConfig(true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.okta_users.test", "users.0.admin_roles.#", "0"),       // skipped
					resource.TestCheckResourceAttr("data.okta_users.test", "users.0.group_memberships.#", "2"), // Everyone, A Group
					resource.TestCheckOutput("output_admin_roles", ""),
					resource.TestMatchOutput("output_group_memberships", allGroupMembershipsRegexp),
				),
			},
		},
	})
}

// TestAccDataSourceOktaUsers_IncludeRoles pertains to https://github.com/okta/terraform-provider-okta/pull/1137 and https://github.com/okta/terraform-provider-okta/issues/1014
func TestAccDataSourceOktaUsers_IncludeRoles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: testOktaUsersRolesGroupsConfig(false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.okta_users.test", "users.0.admin_roles.#", "2"),       // SUPER_ADMIN, APP_ADMIN
					resource.TestCheckResourceAttr("data.okta_users.test", "users.0.group_memberships.#", "0"), // not included
					resource.TestMatchOutput("output_admin_roles", allAdminRolesRegexp),
					resource.TestCheckOutput("output_group_memberships", ""),
				),
			},
		},
	})
}

// TestAccDataSourceOktaUsers_IncludeAll pertains to https://github.com/okta/terraform-provider-okta/pull/1137 and https://github.com/okta/terraform-provider-okta/issues/1014
func TestAccDataSourceOktaUsers_IncludeAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: testOktaUsersRolesGroupsConfig(true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.okta_users.test", "users.0.admin_roles.#", "2"),       // SUPER_ADMIN, APP_ADMIN
					resource.TestCheckResourceAttr("data.okta_users.test", "users.0.group_memberships.#", "2"), // Everyone, A Group
					resource.TestMatchOutput("output_admin_roles", allAdminRolesRegexp),
					resource.TestMatchOutput("output_group_memberships", allGroupMembershipsRegexp),
				),
			},
		},
	})
}

func testOktaUsersRolesGroupsConfig(includeGroups, includeRoles bool) string {
	prepend := `
resource "okta_group" "test" {
  name        = "Example"
  description = "A Group"
}
resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
  lifecycle {
    ignore_changes = [admin_roles]
  }
}
resource "okta_user_admin_roles" "test" {
  user_id     = okta_user.test.id
  admin_roles = [
    "SUPER_ADMIN",
    "APP_ADMIN",
  ]
}
resource "okta_user_group_memberships" "test" {
  user_id = okta_user.test.id
  groups = [
    okta_group.test.id,
  ]
}
data "okta_users" "test" {
  search {
    name = "profile.email"
    comparison = "eq"
    value = okta_user.test.email
  }
  
  delay_read_seconds = 2
`

	var clause string
	if includeGroups {
		clause = "  include_groups = true"
	}
	if includeRoles {
		clause = fmt.Sprintf("%s\n  include_roles = true\n", clause)
	}

	append := `
  depends_on = [
    okta_user.test,
    okta_user_admin_roles.test,
    okta_user_group_memberships.test
  ]
}
output "output_admin_roles" {
  value = join(", ", data.okta_users.test.users.0.admin_roles)
  depends_on = [
    okta_user.test,
    okta_user_admin_roles.test,
    okta_user_group_memberships.test
  ]
}
output "output_group_memberships" {
  value = join(", ", data.okta_users.test.users.0.group_memberships)
  depends_on = [
    okta_user.test,
    okta_user_admin_roles.test,
    okta_user_group_memberships.test
  ]
}
`

	return fmt.Sprintf("%s%s%s", prepend, clause, append)
}
