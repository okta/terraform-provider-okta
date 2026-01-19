package idaas_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/stretchr/testify/require"
)

func TestAccResourceOktaGroupRule_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupRule)
	mgr := newFixtureManager("resources", "okta_group_rule", t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	name := acctest.BuildResourceName(mgr.Seed)
	groupUpdate := mgr.GetFixtures("basic_group_update.tf", t)
	deactivated := mgr.GetFixtures("basic_deactivated.tf", t)
	name2 := acctest.BuildResourceName(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroupRule, doesGroupRuleExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "expression_type", "urn:okta:expression:1.0"),
					resource.TestCheckResourceAttr(resourceName, "expression_value", "String.startsWith(user.firstName,\"andy\")"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
				),
			},
			{
				Config: groupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name2),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
				),
			},
			{
				Config: deactivated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name2),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceOktaGroupRule_invalidHandle(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupRule, t.Name())
	groupResource := fmt.Sprintf("%s.test", resources.OktaIDaaSGroup)
	ruleResource := fmt.Sprintf("%s.inval", resources.OktaIDaaSGroupRule)
	testName := acctest.BuildResourceName(mgr.Seed)
	testSetup := buildInvalidSetup(testName)
	testBuild := buildInvalidBuild(testName)
	testRun := buildInvalidTest(testName)
	testUpdate := buildInvalidUpdate(testName)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroupRule, doesGroupRuleExist),
		Steps: []resource.TestStep{
			{
				Config: testSetup,
				Check:  resource.TestCheckResourceAttr(groupResource, "name", testName),
			},
			{
				Config: testBuild,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ruleResource, "name", testName),
					resource.TestCheckResourceAttr(ruleResource, "status", idaas.StatusActive),
				),
			},
			{
				Config:      testRun,
				Check:       resource.TestCheckResourceAttr(ruleResource, "status", idaas.StatusActive),
				ExpectError: regexp.MustCompile(`group with name .+ does not exist`),
			},
			{
				Config: testUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(groupResource, "name", testName),
					resource.TestCheckResourceAttr(ruleResource, "name", testName),
					resource.TestCheckResourceAttr(ruleResource, "status", idaas.StatusActive),
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
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
	_, response, err := client.Group.GetGroupRule(context.Background(), id, nil)

	return utils.DoesResourceExist(response, err)
}

func TestAccResourceOktaGroupRule_statusIsInvalidDiffFn(t *testing.T) {
	cases := []struct {
		status   string
		expected bool
	}{
		{
			status:   "VALID",
			expected: false,
		},
		{
			status:   "INVALID",
			expected: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.status, func(t *testing.T) {
			result := idaas.StatusIsInvalidDiffFn(tc.status)
			require.EqualValues(t, tc.expected, result)
		})
	}
}

func TestAccResourceOktaGroupRule_nameLengthVerification_Issue2396(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupRule)
	mgr := newFixtureManager("resources", "okta_group_rule", t.Name())
	config := mgr.GetFixtures("basic_group_rule_name_length_verify.tf", t)
	failConfig := mgr.GetFixtures("basic_group_rule_name_length_fail.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroupRule, doesGroupRuleExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "[xx]ZZZ_ああああ_yyyyyyあいうw1w1えおかきくけ1"),
				),
			},
			{
				Config:      failConfig,
				ExpectError: regexp.MustCompile(`\[\{\{\} name\}\] cannot be longer than 50 runes`),
			},
		},
	})
}

func TestAccResourceOktaGroupRule_423ResponseCapture(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupRule)
	mgr := newFixtureManager("resources", "okta_group_rule", t.Name())

	// Test that captures 423 responses but handles them gracefully
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroupRule, doesGroupRuleExist),
		Steps: []resource.TestStep{
			{
				// Create a minimal set of resources that should succeed
				Config: mgr.GetFixtures("test_423_response_capture.tf", t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test_423_response_capture"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "expression_value", "String.startsWith(user.firstName,\"andy\")"),
				),
			},
		},
	})
}
