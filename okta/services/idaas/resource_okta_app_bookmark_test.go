package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/sdk"
)

func TestAccResourceOktaAppBookmarkApplication_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppBookmark, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppBookmark)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppBookmarkApplication_timeouts(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppBookmark, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppBookmark)
	config := `
resource "okta_app_bookmark" "test" {
  label = "testAcc_replace_with_uuid"
  url   = "https://test.com"
  timeouts {
    create = "60m"
    read = "2h"
    update = "30m"
  }
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "30m"),
				),
			},
		},
	})
}

// TestAccAppBookmarkApplication_PR1366 Test for @jakezarobsky-8451 PR #1366
// https://github.com/okta/terraform-provider-okta/pull/1366
func TestAccResourceOktaAppBookmarkApplication_PR1366_authentication_policy(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppBookmark, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppBookmark)
	config := `
data "okta_policy" "test" {
  name = "Any two factors"
  type = "ACCESS_POLICY"
}
resource "okta_app_signon_policy" "test" {
  name        = "testAcc_Policy_replace_with_uuid"
  description = "Sign On Policy"
  depends_on  = [
    data.okta_policy.test
  ]
}
resource "okta_app_bookmark" "test" {
  label  = "testAcc_replace_with_uuid"
  url    = "https://test.com"
  authentication_policy = okta_app_signon_policy.test.id
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_policy"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppBookmarkApplication_authenticationPolicy_OIEonly(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppBookmark, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("app_with_authentication_policy.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppBookmark)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_policy"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_policy"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppBookmarkApplication_skipAuthenticationPolicy(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppBookmark, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppBookmark)
	config := `
resource "okta_app_bookmark" "test" {
  label  = "testAcc_replace_with_uuid"
  url    = "https://test.com"
  skip_authentication_policy = true
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "skip_authentication_policy", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppBookmarkApplication_skipAuthenticationPolicyUpdate(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppBookmark, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppBookmark)
	config := `
resource "okta_app_bookmark" "test" {
  label  = "testAcc_replace_with_uuid"
  url    = "https://test.com"
  skip_authentication_policy = false
}`
	updatedConfig := `
resource "okta_app_bookmark" "test" {
  label  = "testAcc_replace_with_uuid"
  url    = "https://test.com"
  skip_authentication_policy = true
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "skip_authentication_policy", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
			{
				Config: mgr.ConfigReplace(updatedConfig),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "skip_authentication_policy", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}

func TestAccResourceOktaAppBookmarkApplication_skipAuthenticationPolicyWithLogo(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSAppBookmark, t.Name())
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAppBookmark)
	config := `
resource "okta_app_bookmark" "test" {
  label  = "testAcc_replace_with_uuid"
  url    = "https://test.com"
  skip_authentication_policy = true
}`
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSAppBookmark, createDoesAppExist(sdk.NewBookmarkApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(sdk.NewBookmarkApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "skip_authentication_policy", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
				),
			},
		},
	})
}
