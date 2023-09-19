package okta

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaOrgMetadata_read(t *testing.T) {
	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceOktaOrgMetadata,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_org_metadata.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_org_metadata.test", "pipeline"),
					resource.TestCheckResourceAttrSet("data.okta_org_metadata.test", "settings.analytics_collection_enabled"),
					resource.TestCheckResourceAttr("data.okta_org_metadata.test", "domains.organization", fmt.Sprintf("https://%s", os.Getenv("TF_VAR_hostname"))),
				),
			},
		},
	})
}

const testAccExampleDataSourceOktaOrgMetadata = `
data "okta_org_metadata" "test" {}
`
