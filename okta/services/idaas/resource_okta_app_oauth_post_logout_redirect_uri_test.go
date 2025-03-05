package idaas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaAppOAuthApplication_postLogoutRedirectCrud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuthPostLogoutRedirectURI, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuthPostLogoutRedirectURI)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesAppExist(sdk.NewOpenIdConnectApplication())),
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
					resource.TestCheckResourceAttr(resourceName, "id", "https://www.google-updated.com"),
					resource.TestCheckResourceAttr(resourceName, "uri", "https://www.google-updated.com"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
				),
			},
		},
	})
}

func createPostLogoutRedirectURIExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", resourceName)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return missingErr
		}

		uri := rs.Primary.ID
		appID := rs.Primary.Attributes["app_id"]
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
		app := sdk.NewOpenIdConnectApplication()
		_, response, err := client.Application.GetApplication(context.Background(), appID, app, nil)

		// We don't want to consider a 404 an error in some cases and thus the delineation
		if response != nil && response.StatusCode == 404 {
			return missingErr
		} else if err != nil && utils.Contains(app.Settings.OauthClient.PostLogoutRedirectUris, uri) {
			return nil
		}

		return err
	}
}
