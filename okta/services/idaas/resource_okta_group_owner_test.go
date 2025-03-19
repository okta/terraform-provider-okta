package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccResourceOktaGroupOwner_crud(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_group_owner", t.Name())
	config := mgr.GetFixtures("resource.tf", t)

	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			CheckDestroy:             nil,
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("okta_user.test", "first_name", "TestAcc"),
						resource.TestCheckResourceAttr("okta_user.test", "last_name", "Smith"),
						resource.TestCheckResourceAttr("okta_group_owner.test", "type", "USER"),
					),
				},
			},
		})
}
