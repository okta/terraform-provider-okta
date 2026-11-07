package okta4ai_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAiAgentsSchemaDefault_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaAIAgentsSchemaDefault, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaAIAgentsSchemaDefault)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "default"),
					resource.TestCheckResourceAttr(resourceName, "description", "Default AI agent schema"),
					resource.TestCheckResourceAttrSet(resourceName, "definitions.custom.properties"),
				),
			},
			{
				// Singleton: the import ID is ignored and replaced by the fixed "default" ID.
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "default",
				ImportStateVerify: true,
				// Custom properties are written to the API but not read back, so import cannot restore them.
				ImportStateVerifyIgnore: []string{"definitions.base", "definitions.custom.properties"},
			},
		},
	})
}
