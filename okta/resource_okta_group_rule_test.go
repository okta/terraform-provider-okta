package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func sweepGroupRules(client *testClient) error {
	var errorList []error
	// Should never need to deal with pagination
	rules, _, err := client.oktaClient.Group.ListGroupRules(context.Background(), &query.Params{Limit: 300})
	if err != nil {
		return err
	}

	for _, s := range rules {
		if s.Status == statusActive {
			if _, err := client.oktaClient.Group.DeactivateGroupRule(context.Background(), s.Id); err != nil {
				errorList = append(errorList, err)
				continue
			}
		}
		if _, err := client.oktaClient.Group.DeleteGroupRule(context.Background(), s.Id); err != nil {
			errorList = append(errorList, err)
		}
	}
	return condenseError(errorList)
}

func TestAccOktaGroupRule_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", groupRule)
	mgr := newFixtureManager("okta_group_rule")
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	name := buildResourceName(ri)
	ri = acctest.RandInt()
	groupUpdate := mgr.GetFixtures("basic_group_update.tf", ri, t)
	deactivated := mgr.GetFixtures("basic_deactivated.tf", ri, t)
	name2 := buildResourceName(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(groupRule, doesGroupRuleExist),
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
				),
			},
		},
	})
}

func TestAccOktaGroupRule_invalidHandle(t *testing.T) {
	ri := acctest.RandInt()
	groupResource := fmt.Sprintf("%s.test", oktaGroup)
	ruleResource := fmt.Sprintf("%s.inval", groupRule)
	testName := buildResourceName(ri)
	mgr := newFixtureManager(groupRule)
	testSetup := mgr.GetFixtures("inval_setup.tf", ri, t)
	testBuild := mgr.GetFixtures("inval_build.tf", ri, t)
	testRun := mgr.GetFixtures("inval_test.tf", ri, t)
	testUpdate := mgr.GetFixtures("inval_update.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(groupRule, doesGroupRuleExist),
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
				Config: testRun,
				Check:  resource.TestCheckResourceAttr(ruleResource, "status", statusActive),
			},
			{
				Config:   testRun,
				PlanOnly: true,
				Check:    resource.TestCheckResourceAttr(ruleResource, "status", statusInvalid),
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

func doesGroupRuleExist(id string) (bool, error) {
	client := getOktaClientFromMetadata(testAccProvider.Meta())
	_, response, err := client.Group.GetGroupRule(context.Background(), id, nil)

	return doesResourceExist(response, err)
}
