package okta

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaAdminRoleCustomAssignments(t *testing.T) {
	mgr := newFixtureManager(adminRoleCustomAssignments, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", adminRoleCustomAssignments)
	oktaResourceTest(
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
	client := apiSupplementForTest()
	parts := strings.Split(id, "/")
	_, response, err := client.GetResourceSetBinding(context.Background(), parts[0], parts[1])
	return doesResourceExist(response, err)
}
