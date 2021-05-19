package okta

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

// Tests a standard OAuth application with an updated type. This tests the ForceNew on type and tests creating an
// ACTIVE and INACTIVE application via the create action.
func TestAccAppOauth_basic(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "hide_ios", "true"),
					resource.TestCheckResourceAttr(resourceName, "hide_web", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_submit_toolbar", "false"),
					resource.TestCheckResourceAttr(resourceName, "redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_secret", "something_from_somewhere"),
					resource.TestCheckResourceAttr(resourceName, "client_id", "something_from_somewhere"),
					resource.TestCheckResourceAttr(resourceName, "wildcard_redirect", "DISABLED"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "type", "browser"),
					resource.TestCheckResourceAttr(resourceName, "hide_ios", "true"),
					resource.TestCheckResourceAttr(resourceName, "hide_web", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_submit_toolbar", "false"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "client_secret"),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttr(resourceName, "wildcard_redirect", "SUBDOMAIN"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import schema into state")
					}
					return nil
				},
			},
		},
	})
}

// TestAccAppOauth_refreshToken enables refresh token for browser type oauth app
func TestAccAppOauth_refreshToken(t *testing.T) {
	// TODO: This is an "Early Access Feature" and needs to be enabled by Okta
	//       Skipping for now assuming that the okta account doesn't have this feature enabled.
	//       If this feature is enabled or Okta releases this to all this test should be enabled.
	//       SEE https://help.okta.com/en/prod/Content/Topics/Apps/apps-fbm-enable.htm
	t.Skip("This is an 'Early Access Feature' and needs to be enabled by Okta, skipping this test as it fails when this feature is not available")
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("refresh.tf", ri, t)
	update := mgr.GetFixtures("refresh_update.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "type", "browser"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_rotation", "STATIC"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_leeway", "0"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "type", "browser"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_rotation", "ROTATE"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_leeway", "30"),
				),
			},
		},
	})
}

// Tests creation of service app and updates it to native
func TestAccAppOauth_serviceNative(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("service.tf", ri, t)
	updatedConfig := mgr.GetFixtures("native.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "service"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "native"),
				),
			},
		},
	})
}

// Tests creation of service app and updates it to turn on federated broker
func TestAccAppOauth_federationBroker(t *testing.T) {
	// TODO: This is an "Early Access Feature" and needs to be enabled by Okta
	//       Skipping for now assuming that the okta account doesn't have this feature enabled.
	//       If this feature is enabled or Okta releases this to all this test should be enabled.
	//       SEE https://help.okta.com/en/prod/Content/Topics/Apps/apps-fbm-enable.htm
	t.Skip("This is an 'Early Access Feature' and needs to be enabled by Okta, skipping this test as it fails when this feature is not available")

	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("federation_broker_off.tf", ri, t)
	updatedConfig := mgr.GetFixtures("federation_broker_on.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "implicit_assignment", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "implicit_assignment", "true"),
				),
			},
		},
	})
}

// Tests properly errors on conditional requirements.
func TestAccAppOauth_badGrantTypes(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestOAuthConfigBadGrantTypes(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`failed conditional validation for field "grant_types" of type "service", it can contain client_credentials, implicit and must contain client_credentials, received implicit`),
			},
		},
	})
}

// Tests an OAuth application with profile attributes. This tests with a nested JSON object as well as an array.
func TestAccAppOauth_customProfileAttributes(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("custom_attributes.tf", ri, t)
	groupWhitelistConfig := mgr.GetFixtures("group_for_groups_claim.tf", ri, t)
	updatedConfig := mgr.GetFixtures("remove_custom_attributes.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "profile", "{\"customAttribute123\":\"testing-custom-attribute\"}"),
				),
			},
			{
				Config: groupWhitelistConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "profile", fmt.Sprintf("{\"groups\":{\"whitelist\":[\"%s_%d\"]}}", testResourcePrefix, ri)),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "profile", ""),
				),
			},
		},
	})
}

// Tests various expected properties of client_id and custom_client_id
// TODO: remove when custom_client_id is removed
func TestAccAppOauth_customClientID(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				// Create App with custom_client_id set
				Config: buildTestOAuthAppCustomClientID(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "custom_client_id", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "client_id", buildResourceName(ri)),
				),
			},
			{
				// Replace custom_client_id with client_id
				Config: buildTestOAuthAppClientID(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "custom_client_id", ""),
					resource.TestCheckResourceAttr(resourceName, "client_id", buildResourceName(ri)),
				),
			},
		},
	})
}

// TODO: remove when custom_client_id is removed
func TestAccAppOauth_customClientIDError(t *testing.T) {
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config:      buildTestOAuthAppCustomClientIDBadConfig(ri),
				ExpectError: regexp.MustCompile(`"custom_client_id": conflicts with client_id`),
			},
		},
	})
}

// Tests an OAuth application with profile attributes. This tests with a nested JSON object as well as an array.
func TestAccAppOauth_serviceWithJWKS(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appOAuth)
	config := mgr.GetFixtures("service_with_jwks.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appOAuth)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "jwks.0.kty", "RSA"),
					resource.TestCheckResourceAttr(resourceName, "jwks.0.kid", "SIGNING_KEY"),
					resource.TestCheckResourceAttr(resourceName, "jwks.0.e", "AQAB"),
					resource.TestCheckResourceAttr(resourceName, "jwks.0.n", "owfoXNHcAlAVpIO41840ZU2tZraLGw3yEr3xZvAti7oEZPUKCytk88IDgH7440JOuz8GC_D6vtduWOqnEt0j0_faJnhKHgfj7DTWBOCxzSdjrM-Uyj6-e_XLFvZXzYsQvt52PnBJUV15G1W9QTjlghT_pFrW0xrTtbO1c281u1HJdPd5BeIyPb0pGbciySlx53OqGyxrAxPAt5P5h-n36HJkVsSQtNvgptLyOwWYkX50lgnh2szbJ0_O581bqkNBy9uqlnVeK1RZDQUl4mk8roWYhsx_JOgjpC3YyeXA6hHsT5xWZos_gNx98AHivNaAjzIzvyVItX2-hP0Aoscfff"),
				),
			},
		},
	})
}

func createDoesAppExist(app okta.App) func(string) (bool, error) {
	return func(id string) (bool, error) {
		client := getOktaClientFromMetadata(testAccProvider.Meta())
		_, response, err := client.Application.GetApplication(context.Background(), id, app, &query.Params{})

		// We don't want to consider a 404 an error in some cases and thus the delineation
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}

		if err != nil {
			return false, responseErr(response, err)
		}

		return true, err
	}
}

func buildTestOAuthAppClientID(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "test" {
  label          = "%s"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["implicit", "client_credentials"]
  redirect_uris  = ["http://test.com"]
  client_id      = "%s"
}`, appOAuth, name, name)
}

func buildTestOAuthAppCustomClientID(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "test" {
  label            = "%s"
  type             = "service"
  response_types   = ["token"]
  grant_types      = ["implicit", "client_credentials"]
  redirect_uris    = ["http://test.com"]
  custom_client_id = "%s"
}`, appOAuth, name, name)
}

func buildTestOAuthAppCustomClientIDBadConfig(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "test" {
  label            = "%s"
  type             = "service"
  response_types   = ["token"]
  grant_types      = ["implicit", "client_credentials"]
  redirect_uris    = ["http://test.com"]
  custom_client_id = "%s"
  client_id        = "%s"
}`, appOAuth, name, name, name)
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
  response_types   = ["token"]
}
`, appOAuth, name, name)
}
