package idaas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaRealm_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.example", resources.OktaIDaaSRealm)
	mgr := newFixtureManager("resources", resources.OktaIDaaSRealm, t.Name())
	config := mgr.GetFixtures("okta_realm.tf", t)
	updatedConfig := mgr.GetFixtures("okta_realm_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSRealm, doesRealmExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "TestAcc Example Realm"),
					resource.TestCheckResourceAttr(resourceName, "realm_type", "DEFAULT"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "TestAcc Example Realm Updated"),
					resource.TestCheckResourceAttr(resourceName, "realm_type", "PARTNER"),
				),
			},
		},
	})
}

func doesRealmExist(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV5()
	_, response, err := client.RealmAPI.GetRealm(context.Background(), id).Execute()
	return utils.DoesResourceExistV5(response, err)
}
