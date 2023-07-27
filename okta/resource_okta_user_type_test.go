package okta

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOktaUserType_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", userType)
	mgr := newFixtureManager(userType, t.Name())
	config := mgr.GetFixtures("okta_user_type.tf", t)
	updatedConfig := mgr.GetFixtures("okta_user_type_updated.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(userType, doesUserTypeExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Terraform Acceptance Test User Type"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test User Type")),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Terraform Acceptance Test User Type Updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test User Type Updated")),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import schema into state")
					}

					return nil
				},
			},
		},
	})
}

func doesUserTypeExist(id string) (bool, error) {
	client := sdkV2ClientForTest()
	_, response, err := client.UserType.GetUserType(context.Background(), id)
	return doesResourceExist(response, err)
}
