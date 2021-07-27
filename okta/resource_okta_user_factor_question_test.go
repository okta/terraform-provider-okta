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

	mgr := newFixtureManager("okta_user_factor_question")
	config := mgr.GetFixtures("okta_user_factor_question.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", userFactorQuestion)
	resource.Test(
		t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      createUserFactorCheckDestroy(userFactorQuestion),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						// ensureUserFactorExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "security_question_key", "disliked_food"),
						resource.TestCheckResourceAttr(resourceName, "security_answer", "okta"),
					),
				},
			},
		})
}

func createUserFactorCheckDestroy(FactorType string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != FactorType {
				continue
			}
			userID := rs.Primary.Attributes["user_id"]
			ID := rs.Primary.ID
			exists, err := doesUserFactorExistsUpstream(userID, ID)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("user factor still exists,userID: %s, user factor ID: %s", userID, ID)
			}
		}
		return nil
	}
}

func doesUserFactorExistsUpstream(userId string, factorId string) (bool, error) {
	var uf *okta.SecurityQuestionUserFactor
	_, resp, err := getOktaClientFromMetadata(testAccProvider.Meta()).UserFactor.GetFactor(context.Background(), userId, factorId, uf)
	return doesResourceExist(resp, err)
}

func ensureUserFactorExists(name string) resource.TestCheckFunc {
	return nil
}
