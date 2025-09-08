package idaas_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaBehavior_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSBehavior, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	inactive := mgr.GetFixtures("inactive.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSBehavior)
	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSBehavior, doesBehaviorExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(resourceName, "number_of_authentications", "50"),
						resource.TestCheckResourceAttr(resourceName, "location_granularity_type", "LAT_LONG"),
						resource.TestCheckResourceAttr(resourceName, "radius_from_location", "20"),
						resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)+"_updated"),
						resource.TestCheckResourceAttr(resourceName, "number_of_authentications", "100"),
						resource.TestCheckResourceAttr(resourceName, "location_granularity_type", "LAT_LONG"),
						resource.TestCheckResourceAttr(resourceName, "radius_from_location", "5"),
						resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					),
				},
				{
					Config: inactive,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)+"_updated"),
						resource.TestCheckResourceAttr(resourceName, "number_of_authentications", "100"),
						resource.TestCheckResourceAttr(resourceName, "location_granularity_type", "LAT_LONG"),
						resource.TestCheckResourceAttr(resourceName, "radius_from_location", "5"),
						resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					),
				},
			},
		})
}

func doesBehaviorExist(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV5()
	_, response, err := client.BehaviorAPI.GetBehaviorDetectionRule(context.Background(), id).Execute()
	if err != nil {
		if 200 <= response.StatusCode && response.StatusCode <= 299 {
			if strings.HasPrefix(err.Error(), "parsing time") {
				return utils.DoesResourceExistV5(response, nil)
			}
			if strings.Contains(err.Error(), "cannot unmarshal number") {
				return utils.DoesResourceExistV5(response, nil)
			}
		}
	}
	return utils.DoesResourceExistV5(response, err)
}
