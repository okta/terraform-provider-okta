package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccOktaUserFactorQuestion_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(userFactorQuestion)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", userFactorQuestion)
	resource.Test(
		t, resource.TestCase{
			PreCheck:          testClassicOnlyAccPreCheck(t),
			ErrorCheck:        testAccErrorChecks(t),
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      createUserFactorCheckDestroy(userFactorQuestion),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "key", "disliked_food"),
						resource.TestCheckResourceAttr(resourceName, "answer", "meatball"),
						resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "key", "name_of_first_plush_toy"),
						resource.TestCheckResourceAttr(resourceName, "answer", "meatball"),
						resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					),
				},
			},
		})
}

func createUserFactorCheckDestroy(factorType string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != factorType {
				continue
			}
			userID := rs.Primary.Attributes["user_id"]
			ID := rs.Primary.ID
			exists, err := doesUserFactorExistsUpstream(userID, ID)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("user factor still exists, user ID: %s, factor ID: %s", userID, ID)
			}
		}
		return nil
	}
}

func doesUserFactorExistsUpstream(userId, factorId string) (bool, error) {
	var uf *okta.SecurityQuestionUserFactor
	_, resp, err := getOktaClientFromMetadata(testAccProvider.Meta()).UserFactor.GetFactor(context.Background(), userId, factorId, uf)
	return doesResourceExist(resp, err)
}
