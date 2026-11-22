package okta4aa_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaApiServersAuthorizationServer_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaApiServersAuthorizationServer, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaApiServersAuthorizationServer)
	apiServerName := fmt.Sprintf("%s.test", resources.OktaApiServer)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "easu7ekmjWxppv5PFao5"),
					resource.TestCheckResourceAttrPair(resourceName, "api_server_id", apiServerName, "id"),
					resource.TestCheckResourceAttr(resourceName, "orn", "orn:okta:directory:00otxk5i5T0Rk59Z3ao5:external-authorization-servers:easu7ekmjWxppv5PFao5"),
					resource.TestCheckResourceAttrPair(resourceName, "authorization_server.id", resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization_server.issuer"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization_server.status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: importStateIdForApiServersAuthorizationServer(resourceName),
				ImportStateVerify: true,
				// orn is a create-only input; the read endpoint doesn't return it.
				ImportStateVerifyIgnore: []string{"orn"},
			},
		},
	})
}

// importStateIdForApiServersAuthorizationServer builds the {api_server_id}/{id} import ID the resource expects.
func importStateIdForApiServersAuthorizationServer(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["api_server_id"], rs.Primary.ID), nil
	}
}
