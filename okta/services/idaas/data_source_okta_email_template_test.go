package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaEmailTemplate_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSEmailTemplate, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config:  config,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_email_template.forgot_password", "brand_id"),
					resource.TestCheckResourceAttrSet("data.okta_email_template.forgot_password", "name"),
					resource.TestCheckResourceAttr("data.okta_email_template.forgot_password", "name", "ForgotPassword"),
					resource.TestCheckResourceAttrSet("data.okta_email_template.forgot_password", "links"),
				),
			},
		},
	})
}
