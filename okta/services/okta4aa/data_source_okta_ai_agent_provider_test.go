package okta4aa_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// Reads a pre-existing AI Agent provider. A provider can't be created
// inline: registering one needs a source app that's set up for agent import.
func TestAccDataSourceOktaAiAgentProvider_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaAIAgentProvider, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	dataSourceName := fmt.Sprintf("data.%s.test", resources.OktaAIAgentProvider)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "aipued6hgSvnXYSVPao5"),
					resource.TestCheckResourceAttr(dataSourceName, "source_id", "0oaued621rdzmWhadao5"),
					resource.TestCheckResourceAttr(dataSourceName, "source_name", "amazon_aws_sso"),
					resource.TestCheckResourceAttr(dataSourceName, "source_type", "apps"),
					resource.TestCheckResourceAttr(dataSourceName, "source_orn", "orn:okta:idp:00otxk5i5T0Rk59Z3ao5:apps:amazon_aws_sso:0oaued621rdzmWhadao5"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(dataSourceName, "stage"),
					resource.TestCheckResourceAttrSet(dataSourceName, "de_dup_expression"),
					resource.TestCheckResourceAttrSet(dataSourceName, "integration_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "import_in_progress"),
					// Non-string configuration values are flattened with Go's %v formatting, not JSON.
					resource.TestCheckResourceAttr(dataSourceName, "configuration.awsRegions", "[us-east-1]"),
					resource.TestCheckResourceAttr(dataSourceName, "configuration.platforms", "[AWS_BEDROCK_AGENTS]"),
				),
			},
		},
	})
}
