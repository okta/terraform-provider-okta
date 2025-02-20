package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaEmailTemplates_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSEmailTemplates, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_email_templates.test", "email_templates.#"),
					resource.TestCheckResourceAttrSet("data.okta_email_templates.test", "email_templates.0.name"),
					resource.TestCheckResourceAttrSet("data.okta_email_templates.test", "email_templates.0.links"),
					resource.TestCheckResourceAttrSet("data.okta_email_templates.test", "email_templates.1.name"),
					resource.TestCheckResourceAttrSet("data.okta_email_templates.test", "email_templates.1.links"),
				),
			},
		},
	})
}
