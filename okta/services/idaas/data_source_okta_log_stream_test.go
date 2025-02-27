package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceOktaLogStream_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSLogStream, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	awsDataSource := fmt.Sprintf("data.%s.test_by_id", resources.OktaIDaaSLogStream)
	splunkDataSource := fmt.Sprintf("data.%s.test_by_name", resources.OktaIDaaSLogStream)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttrSet(awsDataSource, "id"),
					resource.TestCheckResourceAttr(awsDataSource, "name", fmt.Sprintf("%s AWS", acctest.BuildResourceName(mgr.Seed))),
					resource.TestCheckResourceAttr(awsDataSource, "type", "aws_eventbridge"),
					resource.TestCheckResourceAttr(awsDataSource, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(awsDataSource, "settings.account_id", "123456789012"),
					resource.TestCheckResourceAttr(awsDataSource, "settings.region", "eu-west-3"),
					resource.TestCheckResourceAttr(awsDataSource, "settings.event_source_name", fmt.Sprintf("%s_AWS", acctest.BuildResourceName(mgr.Seed))),

					//resource.TestCheckResourceAttrSet(splunkDataSource, "id"),
					resource.TestCheckResourceAttr(splunkDataSource, "name", fmt.Sprintf("%s Splunk", acctest.BuildResourceName(mgr.Seed))),
					resource.TestCheckResourceAttr(splunkDataSource, "type", "splunk_cloud_logstreaming"),
					resource.TestCheckResourceAttr(splunkDataSource, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(splunkDataSource, "settings.host", "acme.splunkcloud.com"),
					resource.TestCheckResourceAttr(splunkDataSource, "settings.edition", "aws"),
				),
			},
		},
	})
}
