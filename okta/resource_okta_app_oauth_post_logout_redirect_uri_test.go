package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func createPostLogoutRedirectURIExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}

		uri := rs.Primary.ID
		appID := rs.Primary.Attributes["app_id"]
		client := getOktaClientFromMetadata(testAccProvider.Meta())
		app := okta.NewOpenIdConnectApplication()
		_, response, err := client.Application.GetApplication(context.Background(), appID, app, nil)

		// We don't want to consider a 404 an error in some cases and thus the delineation
		if response != nil && response.StatusCode == 404 {
			return missingErr
		} else if err != nil && contains(app.Settings.OauthClient.PostLogoutRedirectUris, uri) {
			return nil
		}

		return err
	}
}

func TestAccAppOAuthApplication_postLogoutRedirectCrud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuthPostLogoutRedirectURI)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuthPostLogoutRedirectURI)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					createPostLogoutRedirectURIExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "uri", "http://google.com"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					createPostLogoutRedirectURIExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "https://www.example-updated.com"),
					resource.TestCheckResourceAttr(resourceName, "uri", "https://www.example-updated.com"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
		},
	})
}
