package idaas_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaRealm_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSRealm, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	configInvalid := mgr.GetFixtures("datasource_not_found.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
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
