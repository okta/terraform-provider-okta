package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceUser_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(user)
	baseConfig := mgr.GetFixtures("datasource.tf", ri, t)
	createUserConfig := mgr.GetFixtures("datasource_create_user.tf", ri, t)

	// NOTE: The ACC tests on the datasource.tf can flap as sometimes these
	// tests can run faster than the Okta org becoming eventually consistent.
	//
	// TF_ACC=1 go test -tags unit -mod=readonly -test.v -run ^TestAccOktaDataSourceUser_read$
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
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

					resource.TestCheckResourceAttrSet("data.okta_user.read_by_id_with_skip", "id"),
					resource.TestCheckResourceAttrSet("data.okta_user.read_by_id_with_skip", "skip_groups"),
					resource.TestCheckResourceAttrSet("data.okta_user.read_by_id_with_skip", "skip_roles"),
					resource.TestCheckResourceAttr("data.okta_user.read_by_id_with_skip", "skip_groups", "true"),
					resource.TestCheckResourceAttr("data.okta_user.read_by_id_with_skip", "skip_roles", "true"),

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

// TestSkipAdminRoles pertains to https://github.com/okta/terraform-provider-okta/pull/1137
func TestSkipAdminRoles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: testRolesGroupsConfig(false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_user.test", "admin_roles.#", "0"),       // skipped
					resource.TestCheckResourceAttr("data.okta_user.test", "group_memberships.#", "2"), // Everyone, A Group
				),
			},
		},
	})
}

// TestSkipGroups pertains to https://github.com/okta/terraform-provider-okta/pull/1137
func TestSkipGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: testRolesGroupsConfig(true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_user.test", "admin_roles.#", "1"),
					resource.TestCheckResourceAttr("data.okta_user.test", "group_memberships.#", "0"), // skipped
				),
			},
		},
	})
}

// TestSkipGroupsSkipRoles pertains to https://github.com/okta/terraform-provider-okta/pull/1137
func TestSkipGroupsSkipRoles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: testRolesGroupsConfig(true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_user.test", "admin_roles.#", "0"),       // skipped
					resource.TestCheckResourceAttr("data.okta_user.test", "group_memberships.#", "0"), // skipped
				),
			},
		},
	})
}

func TestNoSkips(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: testRolesGroupsConfig(false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_user.test", "admin_roles.#", "1"),
					resource.TestCheckResourceAttr("data.okta_user.test", "group_memberships.#", "2"), // Everyone, A Group
				),
			},
		},
	})
}

func testRolesGroupsConfig(skip_groups, skip_roles bool) string {
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
    "APP_ADMIN",
  ]
}
resource "okta_user_group_memberships" "test" {
  user_id = okta_user.test.id
  groups = [
    okta_group.test.id,
  ]
}
data "okta_user" "test" {
  user_id = okta_user.test.id
`

	var clause string
	if skip_groups {
		clause = "  skip_groups = true"
	}
	if skip_roles {
		clause = fmt.Sprintf("%s\nskip_roles = true\n", clause)
	}

	append := `
  depends_on = [
    okta_user_admin_roles.test,
    okta_user_group_memberships.test,
  ]
}`

	return fmt.Sprintf("%s%s%s", prepend, clause, append)
}
