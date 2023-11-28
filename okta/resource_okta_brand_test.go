package okta

import (
	"fmt"
	"regexp"
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

// TestAccResourceOktaBrand_Issue_1824_with_email_domain addresses issue
// https://github.com/okta/terraform-provider-okta/issues/1824 . This test was
// broken, then the brand resource was fixed to make the test pass.
func TestAccResourceOktaBrand_Issue_1824_with_email_domain(t *testing.T) {
	mgr := newFixtureManager("resources", brand, t.Name())
	resourceName := fmt.Sprintf("%s.test", brand)
	step1 := `
resource okta_brand test{
	name = "testAcc-replace_with_uuid"
	locale = "en"
}`
	step2 := `
resource okta_brand test{
	name = "testAcc-replace_with_uuid"
	locale = "en"
}
# HERE when email domain is created, the next time the brand resource is
# refreshed it will now have an email_domain_id and will trigger change
# detection.
resource "okta_email_domain" "test" {
	brand_id     = okta_brand.test.id
	domain       = "testAcc-replace_with_uuid.example.com"
	display_name = "test"
	user_name    = "fff"
}`
	step3 := `
resource okta_brand test{
	name = "testAcc-replace_with_uuid"
	locale = "en"
	email_domain_id = okta_email_domain.test.id
}
resource "okta_email_domain" "test" {
	brand_id     = okta_brand.test.id
	domain       = "testAcc-replace_with_uuid.example.com"
	display_name = "test"
	user_name    = "fff"
}`

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				// Step 1
				// Create okta_brand.test
				Config: mgr.ConfigReplace(step1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "email_domain_id"),
				),
			},
			{
				// Step 2
				// Create okta_email_domain.test with a brand_id from okta_brand.test.id
				// Upon refresh, okta_brand.test will have change detection as the Okta API will
				// reflect the brand having an email_domain_id value.
				Config: mgr.ConfigReplace(step2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "email_domain_id"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				// Step 3
				// Here, when we destroy okta_email_domain.test Terraform runtime will have a cyclic error
				Config: mgr.ConfigReplace(step3),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "email_domain_id"),
				),
				ExpectError: regexp.MustCompile(`.*Cycle: okta_email_domain.test, okta_brand.test.*`),
			},
			{
				// Step 4
				// Rolling back to just having a brand resource after the operator unwound the cyclic error.
				Config: mgr.ConfigReplace(step1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "email_domain_id"),
				),
			},
		},
	})
}
