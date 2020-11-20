package okta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOktaUserType_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", userType)
	mgr := newFixtureManager(userType)
	config := mgr.GetFixtures("okta_user_type.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_user_type_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(userType, doesUserTypeExist),
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
	client := getSupplementFromMetadata(testAccProvider.Meta())
	_, response, err := client.GetUserType(id)

	return doesResourceExist(response, err)
}
