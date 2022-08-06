package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func sweepBehaviors(client *testClient) error {
	var errorList []error
	behaviors, _, err := client.apiSupplement.ListBehaviors(context.Background(), &query.Params{Q: testResourcePrefix})
	if err != nil {
		return err
	}
	for _, b := range behaviors {
		if _, err := client.apiSupplement.DeleteBehavior(context.Background(), b.ID); err != nil {
			errorList = append(errorList, err)
		}
	}
	return condenseError(errorList)
}

func TestAccOktaBehavior(t *testing.T) {
	mgr := newFixtureManager(behavior, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	inactive := mgr.GetFixtures("inactive.tf", t)
	resourceName := fmt.Sprintf("%s.test", behavior)
	oktaResourceTest(
		t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      createCheckResourceDestroy(behavior, doesBehaviorExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(resourceName, "number_of_authentications", "50"),
						resource.TestCheckResourceAttr(resourceName, "location_granularity_type", "LAT_LONG"),
						resource.TestCheckResourceAttr(resourceName, "radius_from_location", "20"),
						resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)+"_updated"),
						resource.TestCheckResourceAttr(resourceName, "number_of_authentications", "100"),
						resource.TestCheckResourceAttr(resourceName, "location_granularity_type", "LAT_LONG"),
						resource.TestCheckResourceAttr(resourceName, "radius_from_location", "5"),
						resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					),
				},
				{
					Config: inactive,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)+"_updated"),
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
	_, response, err := getSupplementFromMetadata(testAccProvider.Meta()).GetBehavior(context.Background(), id)
	return doesResourceExist(response, err)
}
