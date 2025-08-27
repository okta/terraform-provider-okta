package idaas_test

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaAuthServerScope_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthServerScope)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthServerScope, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	importConfig := mgr.GetFixtures("import.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAuthServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "consent", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "name", "test:something"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test display name"),
					resource.TestCheckResourceAttr(resourceName, "system", "false"),
					resource.TestCheckResourceAttr(resourceName, "optional", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "consent", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "name", "test:something"),
					resource.TestCheckResourceAttr(resourceName, "description", "test_updated"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "test display name updated"),
					resource.TestCheckResourceAttr(resourceName, "system", "false"),
					resource.TestCheckResourceAttr(resourceName, "optional", "true"),
				),
			},
			{
				Config: importConfig,
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("failed to find %s", resourceName)
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["auth_server_id"], rs.Primary.Attributes["id"]), nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) (err error) {
					if len(s) != 1 {
						err = errors.New("failed to import into resource into state")
						return
					}

					id := s[0].Attributes["id"]

					if strings.Contains(id, "@") {
						err = fmt.Errorf("user resource id incorrectly set, %s", id)
					}
					return
				},
			},
			{
				// Addresses
				// https://github.com/okta/terraform-provider-okta/issues/1759
				// but benefits all resource imports that are compound input by
				// concatenating input with slashes.
				//
				// Before fixing 1759 this step would cause the panic
				// experienced in 1759. Now, it illustrates the provider will
				// error if input was incorrect as just `auth_server_id` and not
				// the expected `auth_server_id/id`.
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("failed to find %s", resourceName)
					}
					return rs.Primary.Attributes["auth_server_id"], nil
				},
				ExpectError: regexp.MustCompile(`expected 2 import fields "auth_server_id/id", got 1 fields "(\w*)"`),
			},
		},
	})
}
