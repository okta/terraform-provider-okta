package idaas_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func doesGroupRuleExist(id string) (bool, error) {
	client := IDaaSClientForTest(&testing.T{}).OktaSDKClientV2()

	var exists bool
	var lastErr error
	err := retry.RetryContext(context.Background(), time.Minute, func() *retry.RetryError {
		_, response, err := client.Group.GetGroupRule(context.Background(), id, nil)
		existsCheck, checkErr := utils.DoesResourceExist(response, err)
		if checkErr != nil {
			// If the error is a 404 (not found), that's success for destroy check
			if response != nil && response.StatusCode == 404 {
				exists = false
				return nil // Stop retrying, resource is gone
			}
			lastErr = checkErr
			return retry.NonRetryableError(checkErr)
		}
		if !existsCheck {
			// Defensive: treat as not found
			exists = false
			return nil
		}
		// Resource still exists, so retry
		return retry.RetryableError(fmt.Errorf("group rule %s still exists", id))
	})
	if err != nil {
		return false, err
	}
	if lastErr != nil {
		return false, lastErr
	}
	return exists, nil
}

func TestAccResourceOktaGroupRule_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupRule)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupRule, t.Name())

	name := acctest.BuildResourceName(mgr.Seed)
	name2 := acctest.BuildResourceName(mgr.Seed)

	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	groupUpdate := mgr.GetFixtures("basic_group_update.tf", t)
	deactivated := mgr.GetFixtures("basic_deactivated.tf", t)

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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"remove_assigned_users", "expression_validation",
				},
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

func TestAccResourceOktaGroupRule_complex(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupRule, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupRule)
	config := mgr.GetFixtures("complex.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroupRule, doesGroupRuleExist),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{

				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "expression_value",
						"(user.firstName == \"John\" AND user.lastName == \"Doe\") OR user.email == \"john@example.com\""),
				),
			},
		},
	})
}

func TestAccResourceOktaGroupRule_directoryFunctions(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupRule, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupRule)
	config := mgr.GetFixtures("directory_functions.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroupRule, doesGroupRuleExist),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "expression_value",
						"hasDirectoryUser() AND findDirectoryUser().managerUpn == \"manager@example.com\""),
				),
			},
		},
	})
}

func TestAccResourceOktaGroupRule_validConditions(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupRule, t.Name())
	// resourceName := fmt.Sprintf("%s.valid_expression_examples", resources.OktaIDaaSGroupRule)
	validConfig := mgr.GetFixtures("expression_examples_valid.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroupRule, doesGroupRuleExist),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: validConfig,
			},
		},
	})
}

// TODO: capture a 423 response in a VCR cassette for this test
// This is difficult to do with a Okta dev instance due to the lower rate limits
//
// func TestAccResourceOktaGroupRule_423response(t *testing.T) {
// 	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupRule, t.Name())
// 	// resourceName := fmt.Sprintf("%s.valid_expression_examples", resources.OktaIDaaSGroupRule)

// 	// configuration containing many valid rules
// 	testStep1 := mgr.GetFixtures("test_423response_step1.tf", t)
// 	testStep2 := mgr.GetFixtures("test_423response_step2.tf", t)
// 	testStep3 := mgr.GetFixtures("test_423response_step3.tf", t)

// 	acctest.OktaResourceTest(t, resource.TestCase{
// 		PreCheck:                 acctest.AccPreCheck(t),
// 		ErrorCheck:               testAccErrorChecks(t),
// 		CheckDestroy:             nil,
// 		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
// 		Steps: []resource.TestStep{
// 			{Config: testStep1},
// 			{Config: testStep2},
// 			{Config: testStep3},
// 		},
// 	})
// }

func TestAccResourceOktaGroupRule_invalidSyntax(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupRule, t.Name())
	config := mgr.GetFixtures("invalid_syntax.tf", t)
	operatorConfig := mgr.GetFixtures("invalid_syntax_operator.tf", t)
	trailingConfig := mgr.GetFixtures("invalid_syntax_trailing.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroupRule, doesGroupRuleExist),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Invalid expression"),
			},
			{
				Config:      operatorConfig,
				ExpectError: regexp.MustCompile("Invalid expression"),
			},
			{
				Config:      trailingConfig,
				ExpectError: regexp.MustCompile("Invalid expression"),
			},
		},
	})
}

func TestAccResourceOktaGroupRule_stringFunction(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroupRule, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupRule)
	config := mgr.GetFixtures("string_function.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "expression_value",
						"String.stringContains(user.email, \"@example.com\")"),
				),
			},
		},
	})
}

func TestAccResourceOktaGroupRule_nameLengthVerification_Issue2396(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroupRule)
	mgr := newFixtureManager("resources", "okta_group_rule", t.Name())
	config := mgr.GetFixtures("test_group_rule_name_length_verify.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroupRule, doesGroupRuleExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttr(resourceName, "name", "[xx]ZZZ_ああああ_yyyyyyあいうw1w1えおかきくけ1"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc_[xx]ZZZ_ああああ_yyyyyyあいうw1w"),
				),
			},
		},
	})
}

func TestAccResourceOktaGroupRule_nameLengthVerification_Issue2396_Fail(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_group_rule", t.Name())
	failConfig := mgr.GetFixtures("test_group_rule_name_length_fail.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroupRule, doesGroupRuleExist),
		Steps: []resource.TestStep{
			{
				Config:      failConfig,
				ExpectError: regexp.MustCompile(`\[{{} name\}\] cannot be longer than 50 runes`),
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
