package okta

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaOrgMetadata_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_org_metadata", t.Name())
	resourceName := fmt.Sprintf("data.%s.test", "okta_org_metadata")

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("datasource.tf", t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "pipeline"),
					resource.TestCheckResourceAttrSet(resourceName, "settings.analytics_collection_enabled"),
					resource.TestCheckResourceAttr(resourceName, "domains.organization", fmt.Sprintf("https://%s.%s", os.Getenv("OKTA_ORG_NAME"), os.Getenv("OKTA_BASE_URL"))),
				),
			},
		},
	})
}
