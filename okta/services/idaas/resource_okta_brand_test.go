package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaBrand_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSBrand, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSBrand)
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
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.AccMergeProvidersFactoriesForTest(),
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
	mgr := newFixtureManager("resources", resources.OktaIDaaSBrand, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSBrand)
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
resource "okta_email_domain" "test" {
	brand_id     = okta_brand.test.id
	domain       = "testAcc-replace_with_uuid.example.com"
	display_name = "test"
	user_name    = "fff"
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.AccMergeProvidersFactoriesForTest(),
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
				// Create okta_email_domain.test with a brand_id from okta_brand.test.id.
				// okta_brand.test will have not have a computed email_domain_id value until it is refreshed.
				Config: mgr.ConfigReplace(step2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "email_domain_id"),
				),
			},
			{
				RefreshState: true,
				// Step 3
				// Upon refresh, okta_brand.test will have computed email_domain_id value.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "email_domain_id"),
				),
			},
			{
				// Step 4
				// Even though okta_email_domain.test was destroyed, okta_brand.test will have an email_domain_id
				// until the resource is refreshed.
				Config: mgr.ConfigReplace(step1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "email_domain_id"),
				),
			},
			{
				// Step 5
				// okta_brand.test resource shouldn't have an email_domain_id after refresh
				Config: mgr.ConfigReplace(step1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "email_domain_id"),
				),
			},
		},
	})
}

func TestAccResourceOktaBrand_Issue_1846_with_classic_application_uri(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSBrand, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSBrand)
	step1 := `
resource okta_trusted_origin test{
	name   = "testAcc-replace_with_uuid"
	origin = "https://examplesss.com"
	scopes = ["CORS", "REDIRECT"]
}

resource okta_brand test{
	name                                = "testAcc-replace_with_uuid"
	agree_to_custom_privacy_policy      = true
	custom_privacy_policy_url           = "https://example.com"
	default_app_classic_application_uri = "https://examplesss.com"
	locale                              = "en"
	remove_powered_by_okta              = true
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		CheckDestroy:             nil,
		ProtoV5ProviderFactories: acctest.AccMergeProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(step1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc-%d", mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "agree_to_custom_privacy_policy", "true"),
					resource.TestCheckResourceAttr(resourceName, "custom_privacy_policy_url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "default_app_classic_application_uri", "https://examplesss.com"),
					resource.TestCheckResourceAttr(resourceName, "is_default", "false"),
					resource.TestCheckResourceAttr(resourceName, "locale", "en"),
					resource.TestCheckResourceAttr(resourceName, "remove_powered_by_okta", "true"),
				),
			},
		},
	})
}
