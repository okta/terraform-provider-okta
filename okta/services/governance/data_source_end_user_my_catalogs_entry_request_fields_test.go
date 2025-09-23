package governance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaEndUserMyCatalogsEntryRequestFields_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceEndUsersMyCatalogsEntryRequestFields, t.Name())
	config := mgr.GetFixtures("basic.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_end_user_my_catalogs_entry_request_fields.test", "data.0.type", "DURATION"),
					resource.TestCheckResourceAttr("data.okta_end_user_my_catalogs_entry_request_fields.test", "data.0.label", ""),
					resource.TestCheckResourceAttr("data.okta_end_user_my_catalogs_entry_request_fields.test", "data.0.value", "PT8H"),
					resource.TestCheckResourceAttr("data.okta_end_user_my_catalogs_entry_request_fields.test", "metadata.risk_assessment.request_submission_type", "ALLOWED_WITH_NO_OVERRIDES"),
				),
			},
		},
	})
}
