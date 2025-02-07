package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaLogStream_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", logStream, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	awsDataSource := fmt.Sprintf("data.%s.test_by_id", logStream)
	splunkDataSource := fmt.Sprintf("data.%s.test_by_name", logStream)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:                 testAccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: testAccMergeProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet(awsDataSource, "id"),
					resource.TestCheckResourceAttr(awsDataSource, "name", fmt.Sprintf("%s AWS", buildResourceName(mgr.Seed))),
					resource.TestCheckResourceAttr(awsDataSource, "type", "aws_eventbridge"),
					resource.TestCheckResourceAttr(awsDataSource, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(awsDataSource, "settings.account_id", "123456789012"),
					resource.TestCheckResourceAttr(awsDataSource, "settings.region", "eu-west-3"),
					resource.TestCheckResourceAttr(awsDataSource, "settings.event_source_name", fmt.Sprintf("%s_AWS", buildResourceName(mgr.Seed))),

					// resource.TestCheckResourceAttrSet(splunkDataSource, "id"),
					resource.TestCheckResourceAttr(splunkDataSource, "name", fmt.Sprintf("%s Splunk", buildResourceName(mgr.Seed))),
					resource.TestCheckResourceAttr(splunkDataSource, "type", "splunk_cloud_logstreaming"),
					resource.TestCheckResourceAttr(splunkDataSource, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(splunkDataSource, "settings.host", "acme.splunkcloud.com"),
					resource.TestCheckResourceAttr(splunkDataSource, "settings.edition", "aws"),
				),
			},
		},
	})
}
