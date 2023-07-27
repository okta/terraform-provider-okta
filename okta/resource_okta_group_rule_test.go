package okta

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaGroupRule_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", groupRule)
	mgr := newFixtureManager("okta_group_rule", t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	name := buildResourceName(mgr.Seed)
	groupUpdate := mgr.GetFixtures("basic_group_update.tf", t)
	deactivated := mgr.GetFixtures("basic_deactivated.tf", t)
	name2 := buildResourceName(mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(groupRule, doesGroupRuleExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "expression_type", "urn:okta:expression:1.0"),
					resource.TestCheckResourceAttr(resourceName, "expression_value", "String.startsWith(user.firstName,\"andy\")"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
				),
			},
			{
				Config: groupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name2),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
				),
			},
			{
				Config: deactivated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name2),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "1"),
				),
			},
		},
	})
}

func TestAccOktaGroupRule_invalidHandle(t *testing.T) {
	mgr := newFixtureManager(groupRule, t.Name())
	groupResource := fmt.Sprintf("%s.test", group)
	ruleResource := fmt.Sprintf("%s.inval", groupRule)
	testName := buildResourceName(mgr.Seed)
	testSetup := buildInvalidSetup(testName)
	testBuild := buildInvalidBuild(testName)
	testRun := buildInvalidTest(testName)
	testUpdate := buildInvalidUpdate(testName)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(groupRule, doesGroupRuleExist),
		Steps: []resource.TestStep{
			{
				Config: testSetup,
				Check:  resource.TestCheckResourceAttr(groupResource, "name", testName),
			},
			{
				Config: testBuild,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ruleResource, "name", testName),
					resource.TestCheckResourceAttr(ruleResource, "status", statusActive),
				),
			},
			{
				Config:      testRun,
				Check:       resource.TestCheckResourceAttr(ruleResource, "status", statusActive),
				ExpectError: regexp.MustCompile(`group with name .+ does not exist`),
			},
			{
				Config: testUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(groupResource, "name", testName),
					resource.TestCheckResourceAttr(ruleResource, "name", testName),
					resource.TestCheckResourceAttr(ruleResource, "status", statusActive),
				),
			},
		},
	})
}

func buildInvalidBuild(n string) string {
	return fmt.Sprintf(`
resource "okta_group" "test" {
  name = "%s"
}

data "okta_group" "test" {
  name = "%s"
}

resource "okta_group_rule" "inval" {
  name              = "%s"
  status            = "ACTIVE"
  group_assignments = [data.okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,String.toLowerCase(\"bob\"))"
}
`, n, n, n)
}

func buildInvalidSetup(n string) string {
	return fmt.Sprintf(`
resource "okta_group" "test" {
  name = "%s"
}
`, n)
}

func buildInvalidTest(n string) string {
	return fmt.Sprintf(`
data "okta_group" "test" {
  name = "%s"
}

resource "okta_group_rule" "inval" {
  name              = "%s"
  status            = "ACTIVE"
  group_assignments = [data.okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,String.toLowerCase(\"bob\"))"
}
`, n, n)
}

func buildInvalidUpdate(n string) string {
	return fmt.Sprintf(`
resource "okta_group" "test" {
  name = "%s"
}

resource "okta_group_rule" "inval" {
  name              = "%s"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,String.toLowerCase(\"bob\"))"
}
`, n, n)
}

func doesGroupRuleExist(id string) (bool, error) {
	client := sdkV2ClientForTest()
	_, response, err := client.Group.GetGroupRule(context.Background(), id, nil)

	return doesResourceExist(response, err)
}
