package okta

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccAppOAuthApplication_apiScope(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuthAPIScope)
	plainConfig := mgr.GetFixtures("basic.tf", ri, t)
	plainUpdatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test_app_scopes", appOAuthAPIScope)

	// Replace example org url with actual url to prevent API error
	config := strings.ReplaceAll(plainConfig, "https://your.okta.org", getOktaDomainName())
	updatedConfig := strings.ReplaceAll(plainUpdatedConfig, "https://your.okta.org", getOktaDomainName())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, apiScopeExists()),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "issuer"),
					resource.TestCheckTypeSetElemAttr(resourceName, "scopes.*", "okta.users.read"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, apiScopeExists()),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "issuer"),
					resource.TestCheckTypeSetElemAttr(resourceName, "scopes.*", "okta.users.read"),
					resource.TestCheckTypeSetElemAttr(resourceName, "scopes.*", "okta.users.manage"),
				),
			},
		},
	})
}

func apiScopeExists() func(string) (bool, error) {
	return func(id string) (bool, error) {
		scopes, _, err := getOktaClientFromMetadata(testAccProvider.Meta()).Application.ListScopeConsentGrants(context.Background(), id, nil)
		if err != nil {
			return false, fmt.Errorf("failed to get application scope consent grants: %v", err)
		}
		if len(scopes) > 0 {
			return true, nil
		}
		return false, nil
	}
}

func getOktaDomainName() string {
	c, err := oktaConfig()
	if err != nil {
		return ""
	}
	domain := ""
	if c.domain == "" {
		domain = "okta.com"
	} else {
		domain = c.domain
	}
	return fmt.Sprintf("https://%v.%v", c.orgName, domain)
}
