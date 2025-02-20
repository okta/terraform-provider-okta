package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaUserType_read(t *testing.T) {
	resourceName := fmt.Sprintf("data.%s.test", resources.OktaIDaaSUserType)
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSUserType, t.Name())
	createUserType := mgr.GetFixtures("okta_user_type.tf", t)
	readNameConfig := mgr.GetFixtures("read_name.tf", t)
	readIdConfig := mgr.GetFixtures("read_id.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.AccMergeProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config: createUserType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user_type.test", "id"),
				),
			},
			{
				Config: readNameConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: readIdConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(fmt.Sprintf("data.%s.test2", resources.OktaIDaaSUserType), "name"),
				),
			},
		},
	})
}
