package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaBrandCRUD(t *testing.T) {
	mgr := newFixtureManager("resources", brand, t.Name())
	resourceName := fmt.Sprintf("%s.test", brand)
	step1 := `
resource okta_brand test{
	name = "testAcc-replace_with_uuid"
	locale = "en"
}`
	step2 := `
resource okta_brand test{
	name = "testAcc-changed-replace_with_uuid"
	agree_to_custom_privacy_policy = true
	custom_privacy_policy_url = "https://example.com"
	locale = "es"
	remove_powered_by_okta = true
}`
	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(step1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc-%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "agree_to_custom_privacy_policy", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "custom_privacy_policy_url"),
					resource.TestCheckNoResourceAttr(resourceName, "email_domain_id"),
					resource.TestCheckResourceAttr(resourceName, "is_default", "false"),
					resource.TestCheckResourceAttr(resourceName, "locale", "en"),
					resource.TestCheckResourceAttr(resourceName, "remove_powered_by_okta", "false"),
				),
			},
			{
				Config: mgr.ConfigReplace(step2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc-changed-%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "agree_to_custom_privacy_policy", "true"),
					resource.TestCheckResourceAttr(resourceName, "custom_privacy_policy_url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "locale", "es"),
					resource.TestCheckNoResourceAttr(resourceName, "email_domain_id"),
					resource.TestCheckResourceAttr(resourceName, "remove_powered_by_okta", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "default_app_app_instance_id"),
				),
			},
		},
	})
}
