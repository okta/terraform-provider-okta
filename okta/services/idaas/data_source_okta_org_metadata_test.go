package idaas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaOrgMetadata_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_org_metadata", t.Name())
	resourceName := fmt.Sprintf("data.%s.test", "okta_org_metadata")
	var customDomain, customURI string
	customDomain = os.Getenv("OKTA_ACC_TEST_CUSTOM_DOMAIN")
	if customDomain != "" {
		customURI = fmt.Sprintf("https://%s", customDomain)
	}

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.AccMergeProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("datasource.tf", t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "pipeline"),
					resource.TestCheckResourceAttrSet(resourceName, "settings.analytics_collection_enabled"),
					resource.TestCheckResourceAttr(resourceName, "domains.organization", fmt.Sprintf("https://%s.%s", os.Getenv("OKTA_ORG_NAME"), os.Getenv("OKTA_BASE_URL"))),
					resource.TestCheckResourceAttr(resourceName, "domains.alternate", customURI),
				),
			},
		},
	})
}
