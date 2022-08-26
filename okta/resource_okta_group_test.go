package okta

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaGroup_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", group)
	mgr := newFixtureManager(group)
	config := mgr.GetFixtures("okta_group.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_group_updated.tf", ri, t)
	addUsersConfig := mgr.GetFixtures("okta_group_with_users.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(group, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc")),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testAccDifferent")),
			},
			{
				Config: addUsersConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "4"),
				),
			},
		},
	})
}

func TestAccOktaGroup_customschema(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", group)
	mgr := newFixtureManager(group)
	base := mgr.GetFixtures("okta_group_custom_base.tf", ri, t)
	updated := mgr.GetFixtures("okta_group_custom_updated.tf", ri, t)
	removal := mgr.GetFixtures("okta_group_custom_removal.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(group, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(ri))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"testing1234\",\"testSchema2_%s\":true,\"testSchema3_%s\":54321}", strconv.Itoa(ri), strconv.Itoa(ri), strconv.Itoa(ri))),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(ri))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"moretesting1234\",\"testSchema2_%s\":false,\"testSchema3_%s\":12345}", strconv.Itoa(ri), strconv.Itoa(ri), strconv.Itoa(ri))),
				),
			},
			{
				Config: removal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(ri))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"moretesting1234\"}", strconv.Itoa(ri))),
				),
			},
		},
	})
}

func doesGroupExist(id string) (bool, error) {
	_, response, err := getOktaClientFromMetadata(testAccProvider.Meta()).Group.GetGroup(context.Background(), id)
	return doesResourceExist(response, err)
}
