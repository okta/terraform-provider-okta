package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaAppOAuthApplication_redirectCrud(t *testing.T) {
	mgr := newFixtureManager("resources", appOAuthRedirectURI, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", appOAuthRedirectURI)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(appOAuth, createDoesAppExist(sdk.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					createRedirectURIExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "uri", "http://google.com"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					createRedirectURIExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "http://google-updated.com"),
					resource.TestCheckResourceAttr(resourceName, "uri", "http://google-updated.com"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
		},
	})
}

func createRedirectURIExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", resourceName)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return missingErr
		}

		uri := rs.Primary.ID
		appID := rs.Primary.Attributes["app_id"]
		client := sdkV2ClientForTest()
		app := sdk.NewOpenIdConnectApplication()
		_, response, err := client.Application.GetApplication(context.Background(), appID, app, nil)

		// We don't want to consider a 404 an error in some cases and thus the delineation
		if response != nil && response.StatusCode == 404 {
			return missingErr
		} else if err != nil && contains(app.Settings.OauthClient.RedirectUris, uri) {
			return nil
		}

		return err
	}
}
