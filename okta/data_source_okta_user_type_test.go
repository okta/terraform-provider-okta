package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDataSourceUserType_read(t *testing.T) {
	resourceName := fmt.Sprintf("data.%s.test", userType)
	mgr := newFixtureManager(userType, t.Name())
	createUserType := mgr.GetFixtures("okta_user_type.tf", t)
	config := mgr.GetFixtures("datasource.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
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
