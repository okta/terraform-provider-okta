package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaAppUserAssignments_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_app_user_assignments", t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.#"),
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.0.id"),
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.0.status"),
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.0.scope"),
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.0.sync_state"),
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.0.created"),
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.0.last_updated"),
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.0.status_changed"),
					// Test that credentials are present when available
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.0.credentials.#"),
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.0.credentials.0.user_name"),
					// Test that profile data is present
					resource.TestCheckResourceAttrSet("data.okta_app_user_assignments.test", "users.0.profile.%"),
				),
			},
		},
	})
}
