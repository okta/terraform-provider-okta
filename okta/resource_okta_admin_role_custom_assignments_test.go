package okta

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaAdminRoleCustomAssignments(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(adminRoleCustomAssignments)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", adminRoleCustomAssignments)
	resource.Test(
		t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      createCheckResourceDestroy(adminRoleCustomAssignments, doesAdminRoleCustomAssignmentExist),
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
	parts := strings.Split(id, "/")
	_, response, err := getSupplementFromMetadata(testAccProvider.Meta()).GetResourceSetBinding(context.Background(), parts[0], parts[1])
	return doesResourceExist(response, err)
}
