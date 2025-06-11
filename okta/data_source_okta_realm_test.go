package okta

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaRealm_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", realm, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	configInvalid := mgr.GetFixtures("datasource_not_found.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_realm.test", "id"),
					resource.TestCheckResourceAttr("data.okta_realm.test", "name", "AccTest Example Realm"),
					resource.TestCheckResourceAttr("data.okta_realm.test", "realm_type", "DEFAULT"),
				),
			},
			{
				Config:      configInvalid,
				ExpectError: regexp.MustCompile(`Realm with name "Unknown Example Realm" not found`),
			},
		},
	})
}
