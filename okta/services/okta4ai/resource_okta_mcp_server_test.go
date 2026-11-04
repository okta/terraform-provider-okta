package okta4ai_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaMcpServer_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaMcpServer, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaMcpServer)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkMcpServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "orn"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization_server_count"),
					resource.TestCheckResourceAttr(resourceName, "resource_url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "metadata.display_name", "Test MCP Server"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// registration_type and app_name are request-only inputs; Read does not map them back.
				ImportStateVerifyIgnore: []string{"registration_type", "app_name"},
			},
		},
	})
}

func checkMcpServerDestroy(s *terraform.State) error {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return nil
	}
	return nil
}
