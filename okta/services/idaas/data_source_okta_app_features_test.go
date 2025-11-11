package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaAppFeatures_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSAppFeatures, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_features.test", "id"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "app_id", "0oarblaf7hWdLawNg1d7"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "name", "INBOUND_PROVISIONING"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "description", "In-bound provisioning settings for provisioning users from an application to Okta"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "status", "ENABLED"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "capabilities.import_rules.user_create_and_match.allow_partial_match", "true"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "capabilities.import_rules.user_create_and_match.auto_activate_new_users", "false"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "capabilities.import_rules.user_create_and_match.autoconfirm_exact_match", "false"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "capabilities.import_rules.user_create_and_match.autoconfirm_new_users", "false"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "capabilities.import_rules.user_create_and_match.autoconfirm_partial_match", "false"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "capabilities.import_rules.user_create_and_match.exact_match_criteria", "USERNAME"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "capabilities.import_settings.schedule.status", "DISABLED"),
					resource.TestCheckResourceAttr("data.okta_app_features.test", "capabilities.import_settings.username.username_format", "EMAIL"),
				),
			},
		},
	})
}
