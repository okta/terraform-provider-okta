package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaBrandCRUD(t *testing.T) {
	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource okta_brand test{
					name = "test"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_brand.test", "name", "test"),
					resource.TestCheckResourceAttr("okta_brand.test", "agree_to_custom_privacy_policy", "false"),
					resource.TestCheckNoResourceAttr("okta_brand.test", "custom_privacy_policy_url"),
					resource.TestCheckNoResourceAttr("okta_brand.test", "email_domain_id"),
					resource.TestCheckResourceAttr("okta_brand.test", "is_default", "false"),
					resource.TestCheckNoResourceAttr("okta_brand.test", "locale"),
					resource.TestCheckResourceAttr("okta_brand.test", "remove_powered_by_okta", "false"),
				),
			},
			{
				Config: `					
				resource okta_brand test{
					name = "test2"
					agree_to_custom_privacy_policy = true
					custom_privacy_policy_url = "https://example.com"
					locale = "es"
					remove_powered_by_okta = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_brand.test", "name", "test2"),
					resource.TestCheckResourceAttr("okta_brand.test", "agree_to_custom_privacy_policy", "true"),
					resource.TestCheckResourceAttr("okta_brand.test", "custom_privacy_policy_url", "https://example.com"),
					resource.TestCheckResourceAttr("okta_brand.test", "locale", "es"),
					resource.TestCheckNoResourceAttr("okta_brand.test", "email_domain_id"),
					resource.TestCheckResourceAttr("okta_brand.test", "remove_powered_by_okta", "true"),
					resource.TestCheckNoResourceAttr("okta_brand.test", "default_app_app_instance_id"),
				),
			},
		},
	})
}
