package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaApps_read(t *testing.T) {
	mgr := newFixtureManager("datasources", apps, t.Name())
	appsCreate := appsResources
	appsRead := fmt.Sprintf("%s%s", appsResources, appsDataSources)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(appsCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("okta_app_oauth.test1", "id"),
					resource.TestCheckResourceAttrSet("okta_app_oauth.test2", "id"),
					resource.TestCheckResourceAttrSet("okta_app_oauth.test3", "id"),
				),
			},
			{
				Config: mgr.ConfigReplace(appsRead),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.okta_apps.test_by_exact_match", "apps.#", "1"),
					resource.TestCheckResourceAttrSet("data.okta_apps.test_by_exact_match", "apps.#.id"),
					resource.TestCheckResourceAttr("data.okta_apps.test_by_exact_match", "apps.#.label", fmt.Sprintf("testApp_%s_one", buildResourceName(mgr.Seed))),
					resource.TestCheckResourceAttr("data.okta_apps.test_by_exact_match", "apps.#.status", statusActive),

					resource.TestCheckResourceAttr("data.okta_apps.test_by_prefix", "apps.#", "2"),

					resource.TestCheckResourceAttr("data.okta_apps.test_by_no_match", "apps.#", "0"),
				),
			},
		},
	})
}

const appsResources = `
resource "okta_app_oauth" "test1" {
	label          = "testApp_testAcc_replace_with_uuid_one"
	type           = "web"
	grant_types    = ["implicit", "authorization_code"]
	redirect_uris  = ["http://a.com/"]
	response_types = ["code", "token", "id_token"]
	issuer_mode    = "ORG_URL"
	consent_method = "TRUSTED"
}
resource "okta_app_oauth" "test2" {
	label          = "testApp_testAcc_replace_with_uuid_two"
	type           = "web"
	grant_types    = ["implicit", "authorization_code"]
	redirect_uris  = ["http://b.com/"]
	response_types = ["code", "token", "id_token"]
	issuer_mode    = "ORG_URL"
	consent_method = "TRUSTED"
}
resource "okta_app_oauth" "test3" {
	label          = "testAppInvalid_testAcc_replace_with_uuid"
	type           = "web"
	grant_types    = ["implicit", "authorization_code"]
	redirect_uris  = ["http://c.com/"]
	response_types = ["code", "token", "id_token"]
	issuer_mode    = "ORG_URL"
	consent_method = "TRUSTED"
}
`

const appsDataSources = `

data "okta_apps" "test_by_exact_match" {
	label = "testApp_testAcc_replace_with_uuid_one"
}
  
data "okta_apps" "test_by_prefix" {
	label_prefix = "testApp_testAcc_replace_with_uuid_"
}

data "okta_apps" "test_by_no_match" {
	label = "invalidApp_replace_with_uuid"
}
`
