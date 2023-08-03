package okta

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaAuthServerScope_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", authServerScope)
	mgr := newFixtureManager(authServerScope, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	importConfig := mgr.GetFixtures("import.tf", t)

	// NOTE this test will fail, see notes below
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "consent", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "name", "test:something"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					// NOTE there seems to be a bug in the API where it is not returning displayName
					// GET /api/v1/authorizationServers/{asID}/scopes/{scopeID}
					/*
						{
						 "id": "scp7arqqmiW9N2Yub1d7",
						 "name": "test:something",
						 "description": "test",
						 "system": false,
						 "metadataPublish": "ALL_CLIENTS",
						 "apiResourceId": null,
						 "default": false,
						 "_links": {
						  "self": {
						   "href": "https://test.oktapreview.com/api/v1/authorizationServers/aus7arnjxqGcftVoD1d7/scopes/scp7arqqmiW9N2Yub1d7",
						   "hints": {
						    "allow": [
						     "GET",
						     "PUT",
						     "DELETE"
						    ]
						   }
						  }
						 }
						}
					*/
					resource.TestCheckResourceAttr(resourceName, "display_name", "test display name"),
					resource.TestCheckResourceAttr(resourceName, "system", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "consent", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceName, "name", "test:something"),
					resource.TestCheckResourceAttr(resourceName, "description", "test_updated"),
					// NOTE there seems to be a bug in the API where it is not returning displayName
					resource.TestCheckResourceAttr(resourceName, "display_name", "test display name updated"),
					resource.TestCheckResourceAttr(resourceName, "system", "false"),
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
		},
	})
}
