package idaas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: mgr.GetFixtures("datasource.tf", t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "pipeline"),
					resource.TestCheckResourceAttrSet(resourceName, "settings.analytics_collection_enabled"),
					// this check doesn't play well on VCR playback, but it works live
					// resource.TestCheckResourceAttr(resourceName, "domains.organization", fmt.Sprintf("https://%s.%s", oktaOrgNameForTest(), oktaBaseUrlForTest())),
					resource.TestCheckResourceAttrSet(resourceName, "domains.organization"),
					checkResourceIfEnabled(resourceName, "domains.alternate", customURI, customDomain != ""),
				),
			},
		},
	})
}

func checkResourceIfEnabled(resourceName, field, value string, check bool) resource.TestCheckFunc {
	if check {
		return resource.TestCheckResourceAttr(resourceName, field, value)
	}
	return func(s *terraform.State) error {
		return nil
	}
}
