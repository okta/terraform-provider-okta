package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceOktaLogStream_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", logStream, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_log_stream.test", "id"),
					resource.TestCheckResourceAttr("data.okta_log_stream.test", "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttrSet("data.okta_log_stream.test", "type"),
					resource.TestCheckResourceAttr("data.okta_log_stream.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("data.okta_log_stream.test", "settings.region", "eu-west-3"),
				),
			},
		},
	})
}
