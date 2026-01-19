package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaCatalogEntryUserAccessRequesterFields_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceCatalogEntryUserAccessRequestFields, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_catalog_entry_user_access_request_fields.test", "id"),
					resource.TestCheckResourceAttr("data.okta_catalog_entry_user_access_request_fields.test", "data.0.read_only", "false"),
					resource.TestCheckResourceAttr("data.okta_catalog_entry_user_access_request_fields.test", "data.0.type", "TEXT"),
				),
			},
		},
	})
}
