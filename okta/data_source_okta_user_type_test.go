package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceUserType_read(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("data.%s.test", userType)
	mgr := newFixtureManager(userType)
	createUserType := mgr.GetFixtures("okta_user_type.tf", ri, t)
	config := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: createUserType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user_type.test", "id"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
