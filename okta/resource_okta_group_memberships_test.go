package okta

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaGroupMemberships_crud(t *testing.T) {
	ri := acctest.RandInt()

	mgr := newFixtureManager(groupMemberships)
	start := mgr.GetFixtures("basic.tf", ri, t)
	update := mgr.GetFixtures("basic_update.tf", ri, t)
	remove := mgr.GetFixtures("basic_removal.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: start,
			},
			{
				Config: update,
			},
			{
				Config: remove,
			},
		},
	})
}

// TestAccOktaGroupMembershipsIssue1094 addresses https://github.com/okta/terraform-provider-okta/issues/1094
func TestAccOktaGroupMembershipsIssue1094(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				// Before the apply the state will be:
				//   Group A without users.
				//   Users 1, 2 without a group association.
				//   Group B will have users 3, 4, 5 from their creation.
				//   There is a group memberships that will place users 1, 2 into Group A.
				//   There is a group rule that will assign all group B users to Group A.
				// After the apply:
				//   Group A has users 1, 2 by okta_group_memberships resource.
				//   The rule that okta_group_memberships.test_a_direct
				//   describes has been run at Okta associating users 3, 4, and 5
				//   with Group A.
				// Upon the next plan:
				//   The state of okta_group_memberships.test_a_direct
				//   will appear to have drifed from having only two
				//   users to five users becuase
				//   okta_group_rule.group_b_to_a_rule will have run and
				//   associated the three users from Group B to aslo be in
				//   Group A.
				ExpectNonEmptyPlan: true,

				Config: configIssue1094,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_group.test_b", "users.#", "3"),
					resource.TestCheckResourceAttr("data.okta_group.test_a", "users.#", "5"),
				),
			},
		},
	})
}

var configIssue1094 = `
# Given group a, b
resource "okta_group" "test_a" {
  name = "TestACC Group A"
  description = "Group A"
}
resource "okta_group" "test_b" {
  name = "TestACC Group B"
  description = "Group B"
}
# Given users 1, 2, not assigned to any group
resource "okta_user" "test_1" {
  first_name = "TestAcc"
  last_name  = "Smith1"
  login      = "testAcc1-replace_with_uuid@example.com"
  email      = "testAcc1-replace_with_uuid@example.com"
}
resource "okta_user" "test_2" {
  first_name = "TestAcc"
  last_name  = "Smith2"
  login      = "testAcc2-replace_with_uuid@example.com"
  email      = "testAcc2-replace_with_uuid@example.com"
}
# Given users 3, 4, 5 assigned to group B
resource "okta_user" "test_3" {
  first_name        = "TestAcc"
  last_name         = "Smith3"
  login             = "testAcc3-replace_with_uuid@example.com"
  email             = "testAcc3-replace_with_uuid@example.com"
  group_memberships = [okta_group.test_b.id]
  lifecycle {
    ignore_changes = [group_memberships]
  }
}
resource "okta_user" "test_4" {
  first_name        = "TestAcc"
  last_name         = "Smith4"
  login             = "testAcc4-replace_with_uuid@example.com"
  email             = "testAcc4-replace_with_uuid@example.com"
  group_memberships = [okta_group.test_b.id]
  lifecycle {
    ignore_changes = [group_memberships]
  }
}
resource "okta_user" "test_5" {
  first_name        = "TestAcc"
  last_name         = "Smith5"
  login             = "testAcc5-replace_with_uuid@example.com"
  email             = "testAcc5-replace_with_uuid@example.com"
  group_memberships = [okta_group.test_b.id]
  lifecycle {
    ignore_changes = [group_memberships]
  }
}
# Group A should have users 1, 2 assigned via okta_group_memberships
resource "okta_group_memberships" "test_a_direct" {
  group_id = okta_group.test_a.id
  users = [okta_user.test_1.id, okta_user.test_2.id]
  depends_on = [okta_user.test_1, okta_user.test_2, okta_group.test_a]

  # Ignore changes on users if group members will be changed outside of it
  # lifecycle {
  #   ignore_changes = [users]
  # }
}
# Group A should have users 3, 4, 5 assigned via okta_group_rule
resource "okta_group_rule" "group_b_to_a_rule" {
  name                  = "Group B -> A rule"
  status                = "ACTIVE"
  group_assignments     = [okta_group.test_a.id]
  expression_type       = "urn:okta:expression:1.0"
  expression_value      = "isMemberOfAnyGroup(\"${okta_group.test_b.id}\")"
  remove_assigned_users = true
  depends_on = [okta_user.test_3, okta_user.test_4, okta_user.test_5, okta_group.test_a, okta_group.test_b, okta_group_memberships.test_a_direct]
}
# Use a data source to read back in the state of each gorup for testing
# After the group rules run, users 3, 4, 5 should now be in Group A in addition to users 1, 2
data "okta_group" "test_a" {
  id = okta_group.test_a.id
  # There can be eventual consistency issues running a group rule so let's give ours chance to catch up adding group B users to group A
  delay_read_seconds = 5
  include_users = true
  depends_on = [okta_group_memberships.test_a_direct, okta_group_rule.group_b_to_a_rule]
}
data "okta_group" "test_b" {
  id = okta_group.test_b.id
  include_users = true
  depends_on = [okta_group_memberships.test_a_direct, okta_group_rule.group_b_to_a_rule]
}
`

// TestAccOktaGroupMembershipsIssue1119 addresses https://github.com/okta/terraform-provider-okta/issues/1119
func TestAccOktaGroupMembershipsIssue1119(t *testing.T) {
	configUsers := ""
	for i := 0; i < 250; i++ {
		configUsers = fmt.Sprintf("%s%s", configUsers, configUser(i))
	}
	args := []interface{}{
		`
resource "okta_group" "test" {
  name = "TestACC Group"
  description = "Test Group"
}
    `,
		configUsers,
		configGroupMemberships(250),
	}
	strFmt := ""
	for i := 0; i < len(args); i++ {
		strFmt = fmt.Sprintf("%s%s", strFmt, "%s")
	}

	config := fmt.Sprintf(strFmt, args...)
	if !allowLongRunningACCTest(t) {
		t.SkipNow()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_group_memberships.test", "users.#", "250"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("okta_group_memberships.test", "users.#", "250"),
				),
			},
		},
	})
}

func allowLongRunningACCTest(t *testing.T) bool {
	envVar := "OKTA_ALLOW_LONG_RUNNING_ACC_TEST"
	allow := (os.Getenv(envVar) != "")
	if !allow {
		t.Logf("%q not present, skipping test", envVar)
	}
	return allow
}

func configGroupMemberships(n int) string {
	users := []string{}
	for i := 0; i < n; i++ {
		users = append(users, fmt.Sprintf("okta_user.test-%03d.id", i))
	}
	return fmt.Sprintf(`
resource "okta_group_memberships" "test" {
  group_id = okta_group.test.id
  users = [%s]
}
`, strings.Join(users, ", "))
}

func configUser(i int) string {
	return fmt.Sprintf(`
resource "okta_user" "test-%03d" {
  first_name = "TestAcc"
  last_name  = "Smith-%03d"
  login      = "testAcc-%03d-replace_with_uuid@example.com"
  email      = "testAcc-%03d-replace_with_uuid@example.com"
}
`, i, i, i, i)
}
