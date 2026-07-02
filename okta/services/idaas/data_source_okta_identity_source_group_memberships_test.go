package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaIdentitySourceGroupMemberships_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSIdentitySourceGroupMemberships, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_identity_source_group_memberships.test", "identity_source_id", "0oaxc95befZNgrJl71d7"),
					resource.TestCheckResourceAttr("data.okta_identity_source_group_memberships.test", "group_external_id", "GROUPEXT123456784C2IF"),
					resource.TestCheckResourceAttr("data.okta_identity_source_group_memberships.test", "id", "0oaxc95befZNgrJl71d7/GROUPEXT123456784C2IF"),
					// Verify members are populated and span both pages of the paginated response.
					resource.TestCheckResourceAttr("data.okta_identity_source_group_memberships.test", "member_external_ids.#", "3"),
					resource.TestCheckResourceAttr("data.okta_identity_source_group_memberships.test", "member_external_ids.0", "EXT001"),
					resource.TestCheckResourceAttr("data.okta_identity_source_group_memberships.test", "member_external_ids.1", "EXT002"),
					resource.TestCheckResourceAttr("data.okta_identity_source_group_memberships.test", "member_external_ids.2", "EXT003"),
				),
			},
		},
	})
}
