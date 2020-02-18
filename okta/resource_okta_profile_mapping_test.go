package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaProfileMapping_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", oktaProfileMapping)
	mgr := newFixtureManager(oktaProfileMapping)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)
	preventDelete := mgr.GetFixtures("prevent_delete.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(oktaProfileMapping, doesOktaProfileExist),
		Steps: []resource.TestStep{
			{
				Config: preventDelete,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "delete_when_absent", "false"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "delete_when_absent", "true"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "delete_when_absent", "true"),
				),
			},
		},
	})
}

func doesOktaProfileExist(id string) (bool, error) {
	client := getSupplementFromMetadata(testAccProvider.Meta())
	_, response, err := client.GetEmailTemplate(id)

	return doesResourceExist(response, err)
}
