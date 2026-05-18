package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaIdentitySourceGroups_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSIdentitySourceGroups, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_identity_source_groups.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_identity_source_groups.test", "external_id"),
					resource.TestCheckResourceAttr("data.okta_identity_source_groups.test", "profile.display_name", "West Coast users-2"),
					resource.TestCheckResourceAttr("data.okta_identity_source_groups.test", "profile.description", "All users West of The Rockies"),
				),
			},
		},
	})
}
