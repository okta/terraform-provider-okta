package idaas_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaAdminRoleCustomAssignments_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAdminRoleCustomAssignments, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAdminRoleCustomAssignments)
	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:          acctest.AccPreCheck(t),
			ErrorCheck:        testAccErrorChecks(t),
			ProviderFactories: acctest.AccProvidersFactoriesForTest(),
			CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAdminRoleCustomAssignments, doesAdminRoleCustomAssignmentExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					),
				},
			},
		})
}

func doesAdminRoleCustomAssignmentExist(id string) (bool, error) {
	client := provider.SdkSupplementClientForTest()
	parts := strings.Split(id, "/")
	_, response, err := client.GetResourceSetBinding(context.Background(), parts[0], parts[1])
	return utils.DoesResourceExist(response, err)
}
