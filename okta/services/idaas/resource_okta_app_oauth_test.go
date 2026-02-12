package idaas_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

// Tests a standard OAuth application with an updated type. This tests the ForceNew on type and tests creating an
// ACTIVE and INACTIVE application via the create action.
func TestAccResourceOktaAppOauth_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "hide_ios", "true"),
					resource.TestCheckResourceAttr(resourceName, "hide_web", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_submit_toolbar", "false"),
					resource.TestCheckResourceAttr(resourceName, "redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "response_types.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "client_secret"),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
					resource.TestCheckResourceAttr(resourceName, "wildcard_redirect", "DISABLED"),
					resource.TestCheckResourceAttr(resourceName, "participate_slo", "true"),
					resource.TestCheckResourceAttr(resourceName, "frontchannel_logout_uri", "https://example.com/logout"),
					resource.TestCheckResourceAttr(resourceName, "issuer_mode", "ORG_URL"),
					resource.TestCheckResourceAttr(resourceName, "groups_claim.0.type", "EXPRESSION"),
					resource.TestCheckResourceAttr(resourceName, "groups_claim.0.value", "aa"),
					resource.TestCheckResourceAttr(resourceName, "groups_claim.0.name", "bb"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "type", "browser"),
					resource.TestCheckResourceAttr(resourceName, "hide_ios", "true"),
					resource.TestCheckResourceAttr(resourceName, "hide_web", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_submit_toolbar", "false"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "client_secret"),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
					resource.TestCheckResourceAttr(resourceName, "groups_claim.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "wildcard_redirect", "SUBDOMAIN"),
					resource.TestCheckResourceAttr(resourceName, "participate_slo", "true"),
					resource.TestCheckResourceAttr(resourceName, "frontchannel_logout_uri", "https://*.example.com/logout"),
					resource.TestCheckResourceAttr(resourceName, "frontchannel_logout_session_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "issuer_mode", "ORG_URL"),
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

// TestAccResourceOktaAppOauth_refreshToken enables refresh token for browser type oauth app
func TestAccResourceOktaAppOauth_refreshToken(t *testing.T) {
	// TODO: This is an "Early Access Feature" and needs to be enabled by Okta
	//       Skipping for now assuming that the okta account doesn't have this feature enabled.
	//       If this feature is enabled or Okta releases this to all this test should be enabled.
	//       SEE https://help.okta.com/en/prod/Content/Topics/Apps/apps-fbm-enable.htm
	t.Skip("This is an 'Early Access Feature' and needs to be enabled by Okta, skipping this test as it fails when this feature is not available")
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	config := mgr.GetFixtures("refresh.tf", t)
	update := mgr.GetFixtures("refresh_update.tf", t)
	secondUpdate := mgr.GetFixtures("refresh_second_update.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "browser"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_rotation", "STATIC"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_leeway", "0"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "browser"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_rotation", "ROTATE"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_leeway", "0"),
				),
			},
			{
				Config: secondUpdate,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "type", "browser"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_rotation", "ROTATE"),
					resource.TestCheckResourceAttr(resourceName, "refresh_token_leeway", "30"),
				),
			},
		},
	})
}

// Tests creation of service app and updates it to native
func TestAccResourceOktaAppOauth_serviceNative(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	config := mgr.GetFixtures("service.tf", t)
	updatedConfig := mgr.GetFixtures("native.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "service"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "native"),
				),
			},
		},
	})
}

// Tests creation of service app and updates it to turn on federated broker
func TestAccResourceOktaAppOauth_federationBroker(t *testing.T) {
	// TODO: This is an "Early Access Feature" and needs to be enabled by Okta
	//       Skipping for now assuming that the okta account doesn't have this feature enabled.
	//       If this feature is enabled or Okta releases this to all this test should be enabled.
	//       SEE https://help.okta.com/en/prod/Content/Topics/Apps/apps-fbm-enable.htm
	t.Skip("This is an 'Early Access Feature' and needs to be enabled by Okta, skipping this test as it fails when this feature is not available")

	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	config := mgr.GetFixtures("federation_broker_off.tf", t)
	updatedConfig := mgr.GetFixtures("federation_broker_on.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "implicit_assignment", "false"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "implicit_assignment", "true"),
				),
			},
		},
	})
}

// Tests an OAuth application with profile attributes. This tests with a nested JSON object as well as an array.
func TestAccResourceOktaAppOauth_customProfileAttributes(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	configBlankCustomAttributes := mgr.GetFixtures("blank_custom_attributes.tf", t)
	configCustomAttributes := mgr.GetFixtures("custom_attributes.tf", t)
	groupWhitelistConfig := mgr.GetFixtures("group_for_groups_claim.tf", t)
	updatedConfig := mgr.GetFixtures("remove_custom_attributes.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: configBlankCustomAttributes,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "profile", ""),
				),
			},
			{
				Config: configCustomAttributes,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "profile", "{\"customAttribute123\":\"testing-custom-attribute\"}"),
				),
			},
			{
				Config: groupWhitelistConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "profile", fmt.Sprintf("{\"groups\":{\"whitelist\":[\"%s_%d\"]}}", acctest.ResourcePrefixForTest, mgr.Seed)),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "profile", ""),
				),
			},
		},
	})
}

// Tests an OAuth application with profile attributes. This tests with a nested JSON object as well as an array.
func TestAccResourceOktaAppOauth_serviceWithJWKS(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	config := mgr.GetFixtures("service_with_jwks.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)

	ecResourceName := fmt.Sprintf("%s.test_ec", resources.OktaIDaaSAppOAuth)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "jwks.0.kty", "RSA"),
					resource.TestCheckResourceAttr(resourceName, "jwks.0.kid", "SIGNING_KEY_RSA"),
					resource.TestCheckResourceAttr(resourceName, "jwks.0.e", "AQAB"),
					resource.TestCheckResourceAttr(resourceName, "jwks.0.n", "owfoXNHcAlAVpIO41840ZU2tZraLGw3yEr3xZvAti7oEZPUKCytk88IDgH7440JOuz8GC_D6vtduWOqnEt0j0_faJnhKHgfj7DTWBOCxzSdjrM-Uyj6-e_XLFvZXzYsQvt52PnBJUV15G1W9QTjlghT_pFrW0xrTtbO1c281u1HJdPd5BeIyPb0pGbciySlx53OqGyxrAxPAt5P5h-n36HJkVsSQtNvgptLyOwWYkX50lgnh2szbJ0_O581bqkNBy9uqlnVeK1RZDQUl4mk8roWYhsx_JOgjpC3YyeXA6hHsT5xWZos_gNx98AHivNaAjzIzvyVItX2-hP0Aoscfff"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(ecResourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(ecResourceName, "jwks.0.kty", "EC"),
					resource.TestCheckResourceAttr(ecResourceName, "jwks.0.kid", "SIGNING_KEY_EC"),
					resource.TestCheckResourceAttr(ecResourceName, "jwks.0.x", "K37X78mXJHHldZYMzrwipjKR-YZUS2SMye0KindHp6I"),
					resource.TestCheckResourceAttr(ecResourceName, "jwks.0.y", "8IfvsvXWzbFWOZoVOMwgF5p46mUj3kbOVf9Fk0vVVHo"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppOauth_serviceWithJWKSURI(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	config := mgr.GetFixtures("service_with_jwks_uri.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "jwks_uri", "https://example.com"),
				),
			},
		},
	})
}

// createDoesAppExist is a compatibility wrapper for other test files
func createDoesAppExist(app interface{}) func(string) (bool, error) {
	// For OAuth apps, ignore the app parameter and use the v6 SDK helper
	return createDoesOAuthAppExist()
}

// createDoesOAuthAppExist checks if an OAuth application exists using v6 SDK
func createDoesOAuthAppExist() func(string) (bool, error) {
	return func(id string) (bool, error) {
		// Use v6 SDK for consistency with the resource implementation
		client := iDaaSAPIClientForTestUtil.OktaSDKClientV6()
		_, response, err := client.ApplicationAPI.GetApplication(context.Background(), id).Execute()

		// Check if it's a 404 error (app doesn't exist)
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}

		if err != nil {
			return false, err
		}

		return true, nil
	}
}

// TestAccResourceOktaAppOauth_redirect_uris relates to issue 1170
//
//	Enable terraform to maintain order of redirect_uris
//
// https://github.com/okta/terraform-provider-okta/issues/1170
func TestAccResourceOktaAppOauth_redirect_uris(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "okta_app_oauth" "test" {
					label = "example"
					type = "web"
					grant_types = ["authorization_code"]
					wildcard_redirect = "SUBDOMAIN"
					redirect_uris = [
					  "https://one.example.com/",
					  "https://two.example.com/",
					  "https://*.example.com/"
					]
					response_types = ["code"]
				  }
				`,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "redirect_uris.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "wildcard_redirect", "SUBDOMAIN"),
					resource.TestCheckResourceAttr(resourceName, "redirect_uris.0", "https://one.example.com/"),
					resource.TestCheckResourceAttr(resourceName, "redirect_uris.1", "https://two.example.com/"),
					resource.TestCheckResourceAttr(resourceName, "redirect_uris.2", "https://*.example.com/"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppOauth_groups_claim(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	config := `
resource "okta_app_oauth" "test" {
    label                      = "testAcc_replace_with_uuid"
	type                      = "web"
	grant_types                = ["authorization_code"]
	redirect_uris              = ["https://example.com/"]
	response_types             = ["code"]
	issuer_mode                = "ORG_URL"
	groups_claim {
	  type        = "FILTER" # required
	  filter_type = "REGEX"
	  name        = "groups" # required
	  value       = ".*" # required
	}
  }
`

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "issuer_mode", "ORG_URL"),
					resource.TestCheckResourceAttr(resourceName, "groups_claim.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "groups_claim.0.type", "FILTER"),
					resource.TestCheckResourceAttr(resourceName, "groups_claim.0.filter_type", "REGEX"),
					resource.TestCheckResourceAttr(resourceName, "groups_claim.0.value", ".*"),
					resource.TestCheckResourceAttr(resourceName, "groups_claim.0.name", "groups"),
					resource.TestCheckResourceAttr(resourceName, "groups_claim.0.issuer_mode", "ORG_URL"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import schema into state")
					}
					// issue 1536 check if the groups_claim is imported
					rs := s[0]
					if expected, ok := rs.Attributes["groups_claim.#"]; !ok || expected != "1" {
						return errors.New("expected groups_claim to be imported")
					}
					if expected, ok := rs.Attributes["groups_claim.0.type"]; !ok || expected != "FILTER" {
						return errors.New("expected imported groups_claim to have correct type")
					}

					return nil
				},
			},
		},
	})
}

func TestAccResourceOktaAppOauth_timeouts(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)
	config := `
resource "okta_app_oauth" "test" {
  label                      = "testAcc_replace_with_uuid"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["http://d.com/"]
  response_types             = ["code"]
  timeouts {
    create = "60m"
    read = "2h"
    update = "30m"
  }
}
`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "30m"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppOauth_pkce_required(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)
	config := `
resource "okta_app_oauth" "test" {
  label = "testAcc_replace_with_uuid"
  type  = "native"
  pkce_required  = true
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code"]
}
`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "pkce_required", "true"),
				),
			},
		},
	})
}

// TestAccResourceOktaAppOauth_config_combinations
// W.R.T. https://github.com/okta/terraform-provider-okta/issues/1325
// Documentation of the the API behavior of pkce_required when the app type is
// "browser" or "native"
//
// https://developer.okta.com/docs/reference/api/apps/#username-template-object
func TestAccResourceOktaAppOauth_config_combinations(t *testing.T) {
	if acctest.SkipVCRTest(t) {
		// the way this is table tested is not friendly w/ VCR
		return
	}
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())

	cases := []struct {
		name               string
		config             string
		attrPairs          [][]string
		expectErrorMessage string
	}{
		{
			name: "native-pkce-not-set",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "native"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "native"},
				{"pkce_required", "true"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "client_secret_basic"},
			},
		},
		{
			name: "native-pkce-set-true",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "native"
  pkce_required  = true
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "native"},
				{"pkce_required", "true"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "client_secret_basic"},
			},
		},
		{
			name: "native-pkce-set-false",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "native"
  pkce_required  = false
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "native"},
				{"pkce_required", "false"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "client_secret_basic"},
			},
		},
		{
			name: "native-pkce-not-set-token-none",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "native"
  token_endpoint_auth_method = "none"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "native"},
				{"pkce_required", "true"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "none"},
			},
		},
		{
			name: "native-pkce-set-true-token-none",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "native"
  pkce_required  = true
  token_endpoint_auth_method = "none"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "native"},
				{"pkce_required", "true"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "none"},
			},
		},
		{
			name: "native-pkce-set-false-token-none",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "native"
  pkce_required  = false
  token_endpoint_auth_method = "none"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			expectErrorMessage: `''pkce_required'' must be set to true when ''token_endpoint_auth_method'' is ''none''`,
			attrPairs: [][]string{
				{"should-not", "get-here"},
			},
		},
		{
			name: "browser-pkce-not-set",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "browser"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "browser"},
				{"pkce_required", "true"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "client_secret_basic"},
			},
		},
		{
			name: "browser-pkce-set-true",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "browser"
  pkce_required  = true
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "browser"},
				{"pkce_required", "true"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "client_secret_basic"},
			},
		},
		{
			name: "browser-pkce-set-false",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "browser"
  pkce_required  = false
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "browser"},
				{"pkce_required", "false"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "client_secret_basic"},
			},
		},
		{
			name: "browser-pkce-not-set-token-none",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "browser"
  token_endpoint_auth_method = "none"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "browser"},
				{"pkce_required", "true"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "none"},
			},
		},
		{
			name: "browser-pkce-set-true-token-none",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "browser"
  pkce_required  = true
  token_endpoint_auth_method = "none"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "browser"},
				{"pkce_required", "true"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "none"},
			},
		},
		{
			name: "browser-pkce-set-false-token-none",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "browser"
  pkce_required  = false
  token_endpoint_auth_method = "none"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			expectErrorMessage: `''pkce_required'' must be set to true when ''token_endpoint_auth_method'' is ''none''`,
			attrPairs: [][]string{
				{"should-not", "get-here"},
			},
		},
		{
			name: "web-pkce-not-set",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "web"},
				{"pkce_required", "false"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "client_secret_basic"},
			},
		},
		{
			name: "web-pkce-set-true",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  pkce_required  = true
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "web"},
				{"pkce_required", "true"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "client_secret_basic"},
			},
		},
		{
			name: "web-pkce-set-false",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  pkce_required  = false
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			attrPairs: [][]string{
				{"type", "web"},
				{"pkce_required", "false"},
				{"auto_key_rotation", "true"},
				{"token_endpoint_auth_method", "client_secret_basic"},
			},
		},
		{
			name: "web-pkce-not-set-token-none",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  token_endpoint_auth_method = "none"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			expectErrorMessage: `''pkce_required'' must be set to true when ''token_endpoint_auth_method'' is ''none''`,
			attrPairs: [][]string{
				{"should-not", "get-here"},
			},
		},
		{
			name: "web-pkce-set-true-token-none",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  pkce_required  = true
  token_endpoint_auth_method = "none"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			expectErrorMessage: `token_endpoint_auth_method: 'token_endpoint_auth_method' is invalid. Valid values: [client_secret_basic, client_secret_post, client_secret_jwt, private_key_jwt]`,
			attrPairs: [][]string{
				{"should-not", "get-here"},
			},
		},
		{
			name: "web-pkce-set-false-token-none",
			config: `resource "okta_app_oauth" "%s" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  pkce_required  = false
  token_endpoint_auth_method = "none"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://example.com/"]
  response_types = ["code"]
}`,
			expectErrorMessage: `''pkce_required'' must be set to true when ''token_endpoint_auth_method'' is ''none''`,
			attrPairs: [][]string{
				{"should-not", "get-here"},
			},
		},
	}
	for _, test := range cases {
		resourceName := fmt.Sprintf("%s.%s", resources.OktaIDaaSAppOAuth, test.name)
		config := fmt.Sprintf(test.config, test.name)
		testFuncs := []resource.TestCheckFunc{
			ensureResourceExists(resourceName, createDoesOAuthAppExist()),
		}
		for _, pair := range test.attrPairs {
			testFuncs = append(testFuncs, resource.TestCheckResourceAttr(resourceName, pair[0], pair[1]))
		}
		errorCheck := testAccErrorChecks(t)
		if test.expectErrorMessage != "" {
			errorCheck = func(err error) error {
				if err == nil {
					return errors.New("expected an error")
				}
				if !strings.Contains(err.Error(), test.expectErrorMessage) {
					return fmt.Errorf("expected error %q, got %q", test.expectErrorMessage, err.Error())
				}
				return nil
			}
		}

		acctest.OktaResourceTest(t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               errorCheck,
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
			Steps: []resource.TestStep{
				{
					Config: mgr.ConfigReplace(config),
					Check:  resource.ComposeTestCheckFunc(testFuncs...),
				},
			},
		})
	}
}

func TestAccResourceOktaAppOauth_omitSecretSafeEnable(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	omit_secret_off := mgr.GetFixtures("omit_secret_off.tf", t)
	omit_secret_on := mgr.GetFixtures("omit_secret_on.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: omit_secret_off,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
					resource.TestCheckResourceAttrSet(resourceName, "client_secret"),
				),
			},
			{
				Config: omit_secret_on,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
					resource.TestCheckResourceAttr(resourceName, "client_secret", ""),
				),
			},
		},
	})
}

// TestAccResourceOktaAppOauth_2659 covers auto_key_rotation: explicit true/false and default (omit).
// See https://github.com/okta/terraform-provider-okta/issues/2659
func TestAccResourceOktaAppOauth_2659(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)

	// Minimal config: no groups_claim to keep VCR cassette small. Label uses replace_with_uuid for VCR seed.
	createTrue := `
resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  auto_key_rotation = true
  grant_types    = ["authorization_code"]
  redirect_uris = ["https://example.com/callback"]
  response_types = ["code"]
}
`
	explicitFalse := `
resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  auto_key_rotation = false
  grant_types    = ["authorization_code"]
  redirect_uris = ["https://example.com/callback"]
  response_types = ["code"]
}
`
	omitDefault := `
resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["authorization_code"]
  redirect_uris = ["https://example.com/callback"]
  response_types = ["code"]
}
`

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			// Create with explicit true; API returns true.
			{
				Config: mgr.ConfigReplace(createTrue),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "auto_key_rotation", "true"),
				),
			},
			// Case 1: explicit true when existing is true -> no change
			{
				Config: mgr.ConfigReplace(createTrue),
				Check:  resource.TestCheckResourceAttr(resourceName, "auto_key_rotation", "true"),
			},
			// Case 2: explicit false when existing is true -> update to false
			{
				Config: mgr.ConfigReplace(explicitFalse),
				Check:  resource.TestCheckResourceAttr(resourceName, "auto_key_rotation", "false"),
			},
			// Case 4: explicit true when existing is false -> update to true
			{
				Config: mgr.ConfigReplace(createTrue),
				Check:  resource.TestCheckResourceAttr(resourceName, "auto_key_rotation", "true"),
			},
			// Case 3: omit (default) when existing is true -> no change
			{
				Config: mgr.ConfigReplace(omitDefault),
				Check:  resource.TestCheckResourceAttr(resourceName, "auto_key_rotation", "true"),
			},
			// Apply false so next step has existing false
			{
				Config: mgr.ConfigReplace(explicitFalse),
				Check:  resource.TestCheckResourceAttr(resourceName, "auto_key_rotation", "false"),
			},
			// Case 5: explicit false when existing is false -> no change
			{
				Config: mgr.ConfigReplace(explicitFalse),
				Check:  resource.TestCheckResourceAttr(resourceName, "auto_key_rotation", "false"),
			},
			// Case 6: omit (default) when existing is false -> update to true
			{
				Config: mgr.ConfigReplace(omitDefault),
				Check:  resource.TestCheckResourceAttr(resourceName, "auto_key_rotation", "true"),
			},
		},
	})
}

func TestAccResourceOktaAppOauth_1952(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppOAuth, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppOAuth)
	config := `
resource "okta_app_oauth" "test" {
	label         = "MyApp"
	type          = "web"
	redirect_uris = ["http://d.com/"]
	hide_ios      = true
	hide_web      = true
	omit_secret   = true
	}
`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppOAuth, createDoesOAuthAppExist()),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesOAuthAppExist()),
					resource.TestCheckResourceAttr(resourceName, "grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "grant_types.0", "authorization_code"),
				),
			},
		},
	})
}
