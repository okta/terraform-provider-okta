package okta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

// Tests a standard OAuth application with an updated type. This tests the ForceNew on type and tests creating an
// ACTIVE and INACTIVE application via the create action.
func TestAccOktaAppOauthBasic(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("oauth_app.tf", ri, t)
	updatedConfig := mgr.GetFixtures("oauth_app_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_secret", "something_from_somewhere"),
					resource.TestCheckResourceAttr(resourceName, "client_id", "something_from_somewhere"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "browser"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "client_secret"),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
				),
			},
		},
	})
}

// Tests creation of service app and updates it to native
func TestAccOktaAppOauthServiceNative(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("service.tf", ri, t)
	updatedConfig := mgr.GetFixtures("native.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "service"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "native"),
				),
			},
		},
	})
}

// Tests ACTIVE to INACTIVE OAuth application via the update action.
func TestAccOktaAppOauthUpdateStatus(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("oauth_app.tf", ri, t)
	updatedConfig := mgr.GetFixtures("oauth_app_updated_status.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
				),
			},
		},
	})
}

// Add and remove groups/users
func TestAccOktaAppOauthUserGroups(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("oauth_app_groups_and_users.tf", ri, t)
	updatedConfig := mgr.GetFixtures("oauth_app_remove_groups_and_users.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "login_uri", "http://test.com"),
					resource.TestCheckResourceAttr(resourceName, "post_logout_redirect_uris.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "client_secret", ""),
				),
			},
		},
	})
}

// Tests properly errors on conditional requirements.
func TestAccOktaAppOauthBadGrantTypes(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestOAuthConfigBadGrantTypes(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`failed conditional validation for field "grant_types" of type "service", it can contain client_credentials, implicit and must contain client_credentials, received implicit`),
			},
		},
	})
}

func createDoesAppExist(app okta.App) func(string) (bool, error) {
	return func(id string) (bool, error) {
		client := getOktaClientFromMetadata(testAccProvider.Meta())
		_, response, err := client.Application.GetApplication(id, app, &query.Params{})

		// We don't want to consider a 404 an error in some cases and thus the delineation
		if response.StatusCode == 404 {
			return false, nil
		}

		if err != nil {
			return false, responseErr(response, err)
		}

		return true, err
	}
}

func buildTestOAuthConfigBadGrantTypes(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  status      = "ACTIVE"
  label       = "%s"
  type		  = "service"
  grant_types = [ "implicit" ]
  redirect_uris = ["http://d.com/"]
}
`, appOAuth, name, name)
}
