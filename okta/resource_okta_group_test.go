package okta

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func sweepGroups(client *testClient) error {
	var errorList []error
	// Should never need to deal with pagination, limit is 10,000 by default
	groups, _, err := client.oktaClient.Group.ListGroups(context.Background(), &query.Params{Q: testResourcePrefix})
	if err != nil {
		return err
	}

	for _, s := range groups {
		if _, err := client.oktaClient.Group.DeleteGroup(context.Background(), s.Id); err != nil {
			errorList = append(errorList, err)
		}
	}
	return condenseError(errorList)
}

func TestAccOktaGroup_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", group)
	mgr := newFixtureManager(group)
	config := mgr.GetFixtures("okta_group.tf", ri, t)
	updatedConfig := mgr.GetFixtures("okta_group_updated.tf", ri, t)
	addUsersConfig := mgr.GetFixtures("okta_group_with_users.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
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
	updated1 := mgr.GetFixtures("okta_group_custom_updated1.tf", ri, t)
	updated2 := mgr.GetFixtures("okta_group_custom_updated2.tf", ri, t)
	updated3 := mgr.GetFixtures("okta_group_custom_updated3.tf", ri, t)
	cleanup := mgr.GetFixtures("okta_group_custom_cleanup.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(group, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(ri))),
				),
			},
			{
				Config: updated1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(ri))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"testing1234\",\"testSchema2_%s\":true,\"testSchema3_%s\":54321}", strconv.Itoa(ri), strconv.Itoa(ri), strconv.Itoa(ri))),
				),
			},
			{
				Config: updated2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(ri))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"moretesting1234\",\"testSchema2_%s\":false,\"testSchema3_%s\":12345}", strconv.Itoa(ri), strconv.Itoa(ri), strconv.Itoa(ri))),
				),
			},
			{
				Config: updated3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(ri))),
				),
			},
			{
				Config: cleanup,
			},
		},
	})
}

func doesGroupExist(id string) (bool, error) {
	_, response, err := getOktaClientFromMetadata(testAccProvider.Meta()).Group.GetGroup(context.Background(), id)
	return doesResourceExist(response, err)
}
