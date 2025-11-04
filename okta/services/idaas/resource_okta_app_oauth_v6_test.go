package idaas
package idaas

import (
	"testing"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/config"
)

func TestResourceAppOAuth_V6SDK_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOAuthAppConfig_v6_basic(),
				Check: resource.ComposeTestCheckFunc(
					testOAuthAppExists_v6("okta_app_oauth.test"),
					resource.TestCheckResourceAttr("okta_app_oauth.test", "type", "web"),
					resource.TestCheckResourceAttr("okta_app_oauth.test", "label", "Test OAuth App V6"),
					resource.TestCheckResourceAttrSet("okta_app_oauth.test", "client_id"),
				),
			},
			{
				Config: testOAuthAppConfig_v6_updated(),
				Check: resource.ComposeTestCheckFunc(
					testOAuthAppExists_v6("okta_app_oauth.test"),
					resource.TestCheckResourceAttr("okta_app_oauth.test", "label", "Updated OAuth App V6"),
				),
			},
		},
	})
}

func testOAuthAppConfig_v6_basic() string {
	return `
resource "okta_app_oauth" "test" {
  label          = "Test OAuth App V6"
  type           = "web"
  grant_types    = ["authorization_code"]
  response_types = ["code"]
  redirect_uris  = ["https://example.com/callback"]
  token_endpoint_auth_method = "client_secret_basic"
}
`
}

func testOAuthAppConfig_v6_updated() string {
	return `
resource "okta_app_oauth" "test" {
  label          = "Updated OAuth App V6"
  type           = "web"
  grant_types    = ["authorization_code", "refresh_token"]
  response_types = ["code"]
  redirect_uris  = ["https://example.com/callback", "https://example.com/callback2"]
  token_endpoint_auth_method = "client_secret_basic"
}
`
}

func testOAuthAppExists_v6(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client := testAccProvider.Meta().(*config.Config).OktaV6Client
		ctx := context.Background()
		
		_, _, err := client.ApplicationAPI.GetApplication(ctx, rs.Primary.ID).Execute()
		if err != nil {
			return fmt.Errorf("Failed to find OAuth app: %v", err)
		}

		return nil
	}
}