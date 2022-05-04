package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaEmailTemplates_read(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(emailTemplates)
	config := mgr.GetFixtures("datasource.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_email_templates.test", "email_templates.#"),
					resource.TestCheckResourceAttr("data.okta_email_templates.test", "email_templates.#", "10"),
					resource.TestCheckResourceAttrSet("data.okta_email_templates.test", "email_templates.0.name"),
					resource.TestCheckResourceAttrSet("data.okta_email_templates.test", "email_templates.0.links"),
				),
			},
		},
	})
}
