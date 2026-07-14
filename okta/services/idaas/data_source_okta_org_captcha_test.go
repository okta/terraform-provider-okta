package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaOrgCaptcha_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_org_captcha", t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	resourceName := "data.okta_org_captcha.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "captcha_id"),
				),
			},
		},
	})
}
