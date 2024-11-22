package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaDeviceAssurancePolicy_read(t *testing.T) {
	resourceName := fmt.Sprintf("data.%s.test", "okta_device_assurance_policy")
	mgr := newFixtureManager("data-sources", "okta_device_assurance_policy", t.Name())
	createUserType := mgr.GetFixtures("okta_device_assurance_policy.tf", t)
	readNameConfig := mgr.GetFixtures("read_name.tf", t)
	readIdConfig := mgr.GetFixtures("read_id.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: createUserType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_policy_device_assurance_windows.test", "id"),
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
					resource.TestCheckResourceAttrSet(fmt.Sprintf("data.%s.test2", "okta_device_assurance_policy"), "name"),
				),
			},
		},
	})
}
