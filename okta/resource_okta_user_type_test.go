package okta

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func sweepUserTypes(client *testClient) error {
	userTypeList, _, _ := client.oktaClient.UserType.ListUserTypes(context.Background())
	var errorList []error
	for _, ut := range userTypeList {
		if strings.HasPrefix(ut.Name, testResourcePrefix) {
			if _, err := client.oktaClient.UserType.DeleteUserType(context.Background(), ut.Id); err != nil {
				errorList = append(errorList, err)
			}
		}
	}
	return condenseError(errorList)
}

func TestAccOktaUserType_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", userType)
	mgr := newFixtureManager(userType)
	config := mgr.GetFixtures("okta_user_type.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_user_type_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(userType, doesUserTypeExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Terraform Acceptance Test User Type"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test User Type")),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
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
	_, response, err := getOktaClientFromMetadata(testAccProvider.Meta()).UserType.GetUserType(context.Background(), id)
	return doesResourceExist(response, err)
}
