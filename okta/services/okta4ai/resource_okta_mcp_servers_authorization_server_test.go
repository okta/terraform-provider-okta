package okta4ai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaMcpServersAuthorizationServer_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaMcpServersAuthorizationServer, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaMcpServersAuthorizationServer)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "mcp_server_id", "emsu6mjxY5YyMHF0zao5"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization_server.id"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization_server.status"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       importStateIdForMcpServersAuthorizationServer(resourceName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"orn", "protocol_configurations"},
			},
		},
	})
}

// importStateIdForMcpServersAuthorizationServer builds the {mcp_server_id}/{id} import ID the resource expects.
func importStateIdForMcpServersAuthorizationServer(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["mcp_server_id"], rs.Primary.ID), nil
	}
}
