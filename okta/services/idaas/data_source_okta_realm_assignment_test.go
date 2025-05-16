package idaas_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaRealmAssignment_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSRealmAssignment, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	configInvalid := mgr.GetFixtures("datasource_not_found.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_realm_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_realm_assignment.test", "profile_source_id"),
					resource.TestCheckResourceAttrSet("data.okta_realm_assignment.test", "realm_id"),
					resource.TestCheckResourceAttr("data.okta_realm_assignment.test", "name", "AccTest Example Realm Assignment"),
					resource.TestCheckResourceAttr("data.okta_realm_assignment.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("data.okta_realm_assignment.test", "priority", "55"),
					resource.TestCheckResourceAttr("data.okta_realm_assignment.test", "condition_expression", "user.profile.login.contains(\"@acctest.com\")"),
				),
			},
			{
				Config:      configInvalid,
				ExpectError: regexp.MustCompile(`Realm assignment with name "Unknown Example Realm Assignment" not found`),
			},
		},
	})
}
