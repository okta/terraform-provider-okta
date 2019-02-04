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
func TestAccOktaOAuthApplication(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oAuthApp)
	config := mgr.GetFixtures("oauth_app.tf", ri, t)
	updatedConfig := mgr.GetFixtures("oauth_app_updated.tf", ri, t)
	resourceName := buildResourceFQN(oAuthApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(oAuthApp, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					testCheckResourceSliceAttr(resourceName, "grant_types", []string{implicit, authorizationCode}),
					testCheckResourceSliceAttr(resourceName, "redirect_uris", []string{"http://d.com/"}),
					testCheckResourceSliceAttr(resourceName, "response_types", []string{"code", "token", "id_token"}),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "browser"),
					testCheckResourceSliceAttr(resourceName, "grant_types", []string{implicit}),
				),
			},
		},
	})
}

// Tests creation of service app and updates it to native
func TestAccOktaOAuthApplicationServiceNative(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestOAuthConfigService(ri)
	updatedConfig := buildTestOAuthConfigNative(ri)
	resourceName := buildResourceFQN(oAuthApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(oAuthApp, createDoesAppExist(okta.NewOpenIdConnectApplication())),
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
func TestAccOktaOAuthApplicationUpdateStatus(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oAuthApp)
	config := mgr.GetFixtures("oauth_app.tf", ri, t)
	updatedConfig := mgr.GetFixtures("oauth_app_updated_status.tf", ri, t)
	resourceName := buildResourceFQN(oAuthApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(oAuthApp, createDoesAppExist(okta.NewOpenIdConnectApplication())),
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
func TestAccOktaOAuthApplicationUserGroups(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(oAuthApp)
	config := mgr.GetFixtures("oauth_app_groups_and_users.tf", ri, t)
	updatedConfig := mgr.GetFixtures("oauth_app_remove_groups_and_users.tf", ri, t)
	resourceName := buildResourceFQN(oAuthApp, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(oAuthApp, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "user.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "group.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "login_uri", "http://test.com"),
					testCheckResourceSliceAttr(resourceName, "post_logout_redirect_uris", []string{"http://d.com/post"}),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "user.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "group.#", "0"),
				),
			},
		},
	})
}

// Tests properly errors on conditional requirements.
func TestAccOktaOAuthApplicationBadGrantTypes(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestOAuthConfigBadGrantTypes(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`failed conditional validation for field "grant_types" of type "service", it can contain client_credentials, received implicit`),
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
			return false, err
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
`, oAuthApp, name, name)
}

func buildTestOAuthConfigService(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       = "%s"
  type		  = "service"
}
`, oAuthApp, name, name)
}

func buildTestOAuthConfigNative(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label       = "%s"
  type		  = "native"
  grant_types = [ "authorization_code" ]
  redirect_uris = ["http://d.com/"]
}
`, oAuthApp, name, name)
}
