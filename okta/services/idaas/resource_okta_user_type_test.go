package idaas_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaUserType_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSUserType)
	mgr := newFixtureManager("resources", resources.OktaIDaaSUserType, t.Name())
	config := mgr.GetFixtures("okta_user_type.tf", t)
	updatedConfig := mgr.GetFixtures("okta_user_type_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSUserType, doesUserTypeExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Terraform Acceptance Test User Type"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test User Type")),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
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
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
	_, response, err := client.UserType.GetUserType(context.Background(), id)
	return utils.DoesResourceExist(response, err)
}
