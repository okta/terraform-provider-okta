package okta

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceOktaOrgMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceOktaOrgMetadata,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_org_metadata.test", "id"),
					resource.TestCheckResourceAttrSet("data.okta_org_metadata.test", "pipeline"),
					resource.TestCheckResourceAttrSet("data.okta_org_metadata.test", "settings"),
					resource.TestMatchResourceAttr("data.okta_org_metadata.test", "domains.organization", regexp.MustCompile("^.*okta.com$")),
				),
			},
		},
	})
}

const testAccExampleDataSourceOktaOrgMetadata = `
data "okta_org_metadata" "test" {}
`
