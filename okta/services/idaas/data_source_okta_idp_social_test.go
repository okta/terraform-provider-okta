package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaIdpSocial_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSIdpSocial, t.Name())
	preConfig := mgr.GetFixtures("basic.tf", t)
	config := mgr.GetFixtures("datasource.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: preConfig,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_idp_social.test_facebook", "id"),
					resource.TestCheckResourceAttrSet("data.okta_idp_social.test_google", "name"),
					resource.TestCheckResourceAttrSet("data.okta_idp_social.test_microsoft", "id"),
				),
			},
		},
	})
}
