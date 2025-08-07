package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaRealm_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.example", realm)
	mgr := newFixtureManager("resources", realm, t.Name())
	config := mgr.GetFixtures("okta_realm.tf", t)
	updatedConfig := mgr.GetFixtures("okta_realm_updated.tf", t)
	buildResourceName(mgr.Seed)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		CheckDestroy:             checkResourceDestroy(realm, doesRealmExist),
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
	client := sdkV5ClientForTest()
	_, response, err := client.RealmAPI.GetRealm(context.Background(), id).Execute()
	return doesResourceExistV5(response, err)
}
