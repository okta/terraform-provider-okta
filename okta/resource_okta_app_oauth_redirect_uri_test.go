package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/okta/okta-sdk-golang/okta"
)

func createRedirectUriExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}

		uri := rs.Primary.ID
		appId := rs.Primary.Attributes["app_id"]
		client := getOktaClientFromMetadata(testAccProvider.Meta())
		app := okta.NewOpenIdConnectApplication()
		_, response, err := client.Application.GetApplication(appId, app, nil)

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
	mgr := newFixtureManager(appOAuthRedirectUri)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuthRedirectUri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					createRedirectUriExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "uri", "http://google.com"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					createRedirectUriExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", "http://google-updated.com"),
					resource.TestCheckResourceAttr(resourceName, "uri", "http://google-updated.com"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
		},
	})
}
