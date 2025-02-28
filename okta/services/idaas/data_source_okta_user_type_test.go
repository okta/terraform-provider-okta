package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaUserType_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSUserType, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	resourceUserTypeName := acctest.BuildResourceName(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_user_type.test", "id"),

					resource.TestCheckResourceAttr("data.okta_user_type.test-find-by-id", "name", resourceUserTypeName),
					resource.TestCheckResourceAttr("data.okta_user_type.test-find-by-id", "display_name", fmt.Sprintf("%s Name", resourceUserTypeName)),
					resource.TestCheckResourceAttr("data.okta_user_type.test-find-by-id", "description", fmt.Sprintf("%s Description", resourceUserTypeName)),

					resource.TestCheckResourceAttr("data.okta_user_type.test-find-by-name", "name", resourceUserTypeName),
					resource.TestCheckResourceAttr("data.okta_user_type.test-find-by-name", "display_name", fmt.Sprintf("%s Name", resourceUserTypeName)),
					resource.TestCheckResourceAttr("data.okta_user_type.test-find-by-name", "description", fmt.Sprintf("%s Description", resourceUserTypeName)),
				),
			},
		},
	})
}
