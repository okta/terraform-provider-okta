package okta

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
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
					resource.TestCheckResourceAttr(resourceName, "hide_ios", "true"),
					resource.TestCheckResourceAttr(resourceName, "hide_web", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_submit_toolbar", "false"),
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
					resource.TestCheckResourceAttr(resourceName, "hide_ios", "true"),
					resource.TestCheckResourceAttr(resourceName, "hide_web", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_submit_toolbar", "false"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "client_secret"),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("Failed to import schema into state")
					}

					return nil
				},
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

// Tests properly errors on conditional requirements.
func TestAccAppOauth_badGrantTypes(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestOAuthConfigBadGrantTypes(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
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
					resource.TestCheckResourceAttr(resourceName, "profile", "{\"customAttribute123\":\"testing-custom-attribute\"}"),
				),
			},
			{
				Config: groupWhitelistConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "profile", "{\"groups\":{\"whitelist\":[\"Whitelist Group\"]}}"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewOpenIdConnectApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appOAuth, createDoesAppExist(okta.NewOpenIdConnectApplication())),
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
