package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func createRedirectURIExists(name string) resource.TestCheckFunc {
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
		if response.StatusCode == 404 {
			return missingErr
		} else if err != nil && contains(app.Settings.OauthClient.RedirectUris, uri) {
			return nil
		}

		return err
	}
}

func TestAccAppOAuthApplication_redirectCrud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuthRedirectURI)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuthRedirectURI)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
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
