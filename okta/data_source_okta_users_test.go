package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceUsers_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(users)
	users := mgr.GetFixtures("users.tf", ri, t)
	config := mgr.GetFixtures("basic.tf", ri, t)
	dataSource := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				// Ensure users are created
				Config: users,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test1", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test2", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test3", "id"),
				),
			},
			{
				// Ensure data source props are set
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_users.test", "users.#"),
				),
			},
			{
				Config: dataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_users.compound_search", "compound_search_operator"),
					resource.TestCheckResourceAttr("data.okta_users.compound_search", "compound_search_operator", "and"),
					resource.TestCheckResourceAttrSet("data.okta_users.compound_search", "users.#"),
					resource.TestCheckResourceAttr("data.okta_users.compound_search", "users.#", "1"),
				),
			},
		},
	})
}

func TestAccOktaDataSourceUsers_readWithGroupId(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(users)
	users := mgr.GetFixtures("users_with_group.tf", ri, t)
	config := mgr.GetFixtures("group.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				// Ensure user and group are created
				Config: users,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test1", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test2", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test3", "id"),
					resource.TestCheckResourceAttrSet("okta_group.test", "id"),
					resource.TestCheckResourceAttr("okta_group_memberships.test", "users.#", "2"),
				),
			},
			{
				// Ensure data source props are set
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_users.test", "users.#", "2"),
				),
			},
		},
	})
}

func TestAccOktaDataSourceUsers_readWithGroupIdIncludingGroups(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(users)
	users := mgr.GetFixtures("users_with_group.tf", ri, t)
	config := mgr.GetFixtures("group_with_groups.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				// Ensure user and group are created
				Config: users,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user.test", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test1", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test2", "id"),
					resource.TestCheckResourceAttrSet("okta_user.test3", "id"),
					resource.TestCheckResourceAttrSet("okta_group.test", "id"),
					resource.TestCheckResourceAttr("okta_group_memberships.test", "users.#", "2"),
				),
			},
			{
				// Ensure data source props are set
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_users.test", "users.#", "2"),
				),
			},
		},
	})
}
