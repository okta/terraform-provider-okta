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

func TestAccResourceOktaA2aServersAuthorizationServer_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaA2aServersAuthorizationServer, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaA2aServersAuthorizationServer)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkA2aServersAuthorizationServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "a2a_server_id", "wlpu37lfb7LQOikYDao5"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization_server.id"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization_server.status"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       importStateIdForA2aServersAuthorizationServer(resourceName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"orn"},
			},
		},
	})
}

// importStateIdForA2aServersAuthorizationServer builds the {a2a_server_id}/{id} import ID the resource expects.
func importStateIdForA2aServersAuthorizationServer(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["a2a_server_id"], rs.Primary.ID), nil
	}
}

func checkA2aServersAuthorizationServerDestroy(s *terraform.State) error {
	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return nil
	}
	return nil
}
