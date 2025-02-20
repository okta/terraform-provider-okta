package idaas_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/stretchr/testify/require"
)

func TestAppUpdateStatus(t *testing.T) {
	before := func() func() {
		o := os.Getenv("OKTA_ORG_NAME")
		u := os.Getenv("OKTA_BASE_URL")
		t := os.Getenv("OKTA_API_TOKEN")
		os.Setenv("OKTA_ORG_NAME", "test")
		os.Setenv("OKTA_BASE_URL", "example.com")
		os.Setenv("OKTA_API_TOKEN", "token")

		return func() {
			os.Setenv("OKTA_ORG_NAME", o)
			os.Setenv("OKTA_BASE_URL", u)
			os.Setenv("OKTA_API_TOKEN", t)
		}
	}
	after := before()
	defer after()

	apiHostname := fmt.Sprintf("https://%s.%s", os.Getenv("OKTA_ORG_NAME"), os.Getenv("OKTA_BASE_URL"))
	activateURI := fmt.Sprintf("%s/api/v1/apps/123/lifecycle/activate", apiHostname)
	deactivateURI := fmt.Sprintf("%s/api/v1/apps/123/lifecycle/deactivate", apiHostname)
	mockActivateURI := fmt.Sprintf("POST %s", activateURI)
	mockDeactivateURI := fmt.Sprintf("POST %s", deactivateURI)

	testAppSchema := idaas.BuildAppSchema(map[string]*schema.Schema{
		"test": {
			Type:     schema.TypeString,
			Optional: true,
		},
	})
	cases := []struct {
		name                   string
		state                  *terraform.InstanceState
		diff                   *terraform.InstanceDiff
		expectAddtionalChanges bool
		expectAPIcall          string
	}{
		{
			name: "deactivate an app only",
			// NOTE: diff represents new state
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"status": {
						Old: "ACTIVE",
						New: "INACTIVE",
					},
				},
			},
			state: &terraform.InstanceState{
				Attributes: map[string]string{
					"id":     "123",
					"status": "ACTIVE",
					"test":   "test",
				},
			},
			expectAddtionalChanges: false,
			expectAPIcall:          mockDeactivateURI,
		},
		{
			name: "activate an app only",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"status": {
						Old: "INACTIVE",
						New: "ACTIVE",
					},
				},
			},
			state: &terraform.InstanceState{
				Attributes: map[string]string{
					"id":     "123",
					"status": "INACTIVE",
					"test":   "test",
				},
			},
			expectAPIcall:          mockActivateURI,
			expectAddtionalChanges: false,
		},
		{
			//  deactivate app and update app values
			name: "deactivate an app and update app values",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"status": {
						Old: "ACTIVE",
						New: "INACTIVE",
					},
					"test": {
						Old: "old",
						New: "new",
					},
				},
			},
			state: &terraform.InstanceState{
				Attributes: map[string]string{
					"id":     "123",
					"status": "ACTIVE",
					"test":   "old",
				},
			},
			expectAPIcall:          mockDeactivateURI,
			expectAddtionalChanges: true,
		},
		{
			name: "activate an app and update app values",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"status": {
						Old: "INACTIVE",
						New: "ACTIVE",
					},
					"test": {
						Old: "old",
						New: "new",
					},
				},
			},
			state: &terraform.InstanceState{
				Attributes: map[string]string{
					"id":     "123",
					"status": "INACTIVE",
					"test":   "old",
				},
			},
			expectAPIcall:          mockActivateURI,
			expectAddtionalChanges: true,
		},
		{
			name: "update inactive app",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"test": {
						Old: "old",
						New: "new",
					},
				},
			},
			state: &terraform.InstanceState{
				Attributes: map[string]string{
					"id":     "123",
					"status": "INACTIVE",
					"test":   "old",
				},
			},
			expectAddtionalChanges: true,
		},
		{
			name: "update active app",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"test": {
						Old: "old",
						New: "new",
					},
				},
			},
			state: &terraform.InstanceState{
				Attributes: map[string]string{
					"id":     "123",
					"status": "ACTIVE",
					"test":   "old",
				},
			},
			expectAddtionalChanges: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := schema.InternalMap(testAppSchema).Data(tc.state, tc.diff)
			require.NoError(t, err)

			var m interface{}
			c := config.NewConfig(d)
			err = c.LoadClients()
			require.NoError(t, err)
			m = c

			defer httpmock.DeactivateAndReset()
			httpmock.ActivateNonDefault(c.OktaSDKClientV2.GetConfig().HttpClient)
			httpmock.RegisterResponder("POST", activateURI, httpmock.NewStringResponder(200, ""))
			httpmock.RegisterResponder("POST", deactivateURI, httpmock.NewStringResponder(200, ""))

			additionalChanges, err := idaas.AppUpdateStatus(context.Background(), d, m)
			require.NoError(t, err)
			require.Equal(t, tc.expectAddtionalChanges, additionalChanges)
			if tc.expectAPIcall != "" {
				apiCallInfo := httpmock.GetCallCountInfo()
				count, found := apiCallInfo[tc.expectAPIcall]
				if !found || count != 1 {
					t.Fatalf("expected api info\n%+v\nto contain %q\nwith count 1, got %d", apiCallInfo, tc.expectAPIcall, count)
				}
			}
		})
	}
}

// TestAppUpdateStatus_all_apps DRY up testing all the Okta app's implementation
// of update in one place
func TestAppUpdateStatus_all_apps(t *testing.T) {
	cases := []struct {
		name         string
		resource     string
		attrName     string
		value1       string
		value2       string
		value3       string
		value4       string
		value5       string
		body         string
		appPrototype sdk.App
	}{
		{
			name:     "update auto login app",
			resource: resources.OktaIDaaSAppAutoLogin,
			attrName: "app_settings_json",
			value1:   `{"baseUrl":"https://example.com/1","subdomain":"articulate"}`,
			value2:   `{"baseUrl":"https://example.com/2","subdomain":"articulate"}`,
			value3:   `{"baseUrl":"https://example.com/3","subdomain":"articulate"}`,
			value4:   `{"baseUrl":"https://example.com/4","subdomain":"articulate"}`,
			value5:   `{"baseUrl":"https://example.com/5","subdomain":"articulate"}`,
			body:     `preconfigured_app="pagerduty"`,

			appPrototype: sdk.NewAutoLoginApplication(),
		},
		{
			name:         "update basic auth app",
			resource:     resources.OktaIDaaSAppBasicAuth,
			attrName:     "auth_url",
			value1:       "https://example.com/login-1",
			value2:       "https://example.com/login-2",
			value3:       "https://example.com/login-3",
			value4:       "https://example.com/login-4",
			value5:       "https://example.com/login-5",
			body:         `url = "https://example.com/login.html"`,
			appPrototype: sdk.NewBasicAuthApplication(),
		},
		{
			name:         "update bookmark app",
			resource:     resources.OktaIDaaSAppBookmark,
			attrName:     "url",
			value1:       "https://example.com/1",
			value2:       "https://example.com/2",
			value3:       "https://example.com/3",
			value4:       "https://example.com/4",
			value5:       "https://example.com/5",
			body:         "",
			appPrototype: sdk.NewBookmarkApplication(),
		},
		{
			name:         "update oauth app",
			resource:     resources.OktaIDaaSAppOAuth,
			attrName:     "logo_uri",
			value1:       "https://example.com/logo-1.png",
			value2:       "https://example.com/logo-2.png",
			value3:       "https://example.com/logo-3.png",
			value4:       "https://example.com/logo-4.png",
			value5:       "https://example.com/logo-5.png",
			body:         testAppBodyOAuth,
			appPrototype: sdk.NewOpenIdConnectApplication(),
		},
		{
			name:         "update saml app",
			resource:     resources.OktaIDaaSAppSaml,
			attrName:     "app_settings_json",
			value1:       `{"baseUrl":"https://example.com/1","subdomain":"articulate"}`,
			value2:       `{"baseUrl":"https://example.com/2","subdomain":"articulate"}`,
			value3:       `{"baseUrl":"https://example.com/3","subdomain":"articulate"}`,
			value4:       `{"baseUrl":"https://example.com/4","subdomain":"articulate"}`,
			value5:       `{"baseUrl":"https://example.com/5","subdomain":"articulate"}`,
			body:         `preconfigured_app="pagerduty"`,
			appPrototype: sdk.NewSamlApplication(),
		},
		{
			name:         "update secure password store app",
			resource:     resources.OktaIDaaSAppSecurePasswordStore,
			attrName:     "url",
			value1:       "https://exmaple.com/1",
			value2:       "https://exmaple.com/2",
			value3:       "https://exmaple.com/3",
			value4:       "https://exmaple.com/4",
			value5:       "https://exmaple.com/5",
			body:         testAppSecurePasswordStore,
			appPrototype: sdk.NewSecurePasswordStoreApplication(),
		},
		{
			name:         "update shared credentials app",
			resource:     resources.OktaIDaaSAppSharedCredentials,
			attrName:     "admin_note",
			value1:       "note 1",
			value2:       "note 2",
			value3:       "note 3",
			value4:       "note 4",
			value5:       "note 5",
			body:         testAppSharedCredentials,
			appPrototype: sdk.NewBrowserPluginApplication(),
		},
		{
			name:         "update SWA app",
			resource:     resources.OktaIDaaSAppSwa,
			attrName:     "url",
			value1:       "https://exmaple.com/1",
			value2:       "https://exmaple.com/2",
			value3:       "https://exmaple.com/3",
			value4:       "https://exmaple.com/4",
			value5:       "https://exmaple.com/5",
			body:         testAppSWA,
			appPrototype: sdk.NewSwaApplication(),
		},
		{
			name:         "update three field app",
			resource:     resources.OktaIDaaSAppThreeField,
			attrName:     "url",
			value1:       "https://exmaple.com/1",
			value2:       "https://exmaple.com/2",
			value3:       "https://exmaple.com/3",
			value4:       "https://exmaple.com/4",
			value5:       "https://exmaple.com/5",
			body:         testAppThreeField,
			appPrototype: sdk.NewSwaThreeFieldApplication(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mgr := newFixtureManager("resources", tc.resource, t.Name())
			resourceName := fmt.Sprintf("%s.test", tc.resource)
			config := `
resource %q "test" {
  label = "testAcc_replace_with_uuid"
  status = %q
  %s = %q
  %s
}`

			acctest.OktaResourceTest(t, resource.TestCase{
				PreCheck:          acctest.AccPreCheck(t),
				ErrorCheck:        testAccErrorChecks(t),
				ProviderFactories: acctest.AccProvidersFactoriesForTest(),
				CheckDestroy:      checkResourceDestroy(tc.resource, createDoesAppExist(tc.appPrototype)),
				Steps: []resource.TestStep{
					{
						// 1 - create an active app
						Config: mgr.ConfigReplace(fmt.Sprintf(config, tc.resource, "ACTIVE", tc.attrName, tc.value1, tc.body)),
						Check: resource.ComposeTestCheckFunc(
							ensureResourceExists(resourceName, createDoesAppExist(tc.appPrototype)),
							resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
							resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
							resource.TestCheckResourceAttr(resourceName, tc.attrName, tc.value1),
						),
					},
					{
						// 2 - deactivate an app only
						Config: mgr.ConfigReplace(fmt.Sprintf(config, tc.resource, "INACTIVE", tc.attrName, tc.value1, tc.body)),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
							resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
							resource.TestCheckResourceAttr(resourceName, tc.attrName, tc.value1),
						),
					},
					{
						// 3 - activate an app only
						Config: mgr.ConfigReplace(fmt.Sprintf(config, tc.resource, "ACTIVE", tc.attrName, tc.value1, tc.body)),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
							resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
							resource.TestCheckResourceAttr(resourceName, tc.attrName, tc.value1),
						),
					},
					{
						// 4 - deactivate an app and update some arguments
						Config: mgr.ConfigReplace(fmt.Sprintf(config, tc.resource, "INACTIVE", tc.attrName, tc.value2, tc.body)),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
							resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
							resource.TestCheckResourceAttr(resourceName, tc.attrName, tc.value2),
						),
					},
					{
						// 5 - update some arguments on a deactivated app
						Config: mgr.ConfigReplace(fmt.Sprintf(config, tc.resource, "INACTIVE", tc.attrName, tc.value3, tc.body)),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
							resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
							resource.TestCheckResourceAttr(resourceName, tc.attrName, tc.value3),
						),
					},
					{
						// 6 - activate an app and update some arguments
						Config: mgr.ConfigReplace(fmt.Sprintf(config, tc.resource, "INACTIVE", tc.attrName, tc.value4, tc.body)),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
							resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
							resource.TestCheckResourceAttr(resourceName, tc.attrName, tc.value4),
						),
					},
					{
						// 7 - update some arguments on an active app
						Config: mgr.ConfigReplace(fmt.Sprintf(config, tc.resource, "ACTIVE", tc.attrName, tc.value5, tc.body)),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "label", acctest.BuildResourceName(mgr.Seed)),
							resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
							resource.TestCheckResourceAttr(resourceName, tc.attrName, tc.value5),
						),
					},
				},
			})
		})
	}
}

const (
	testAppBodyOAuth = `
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["http://d.com/"]
  response_types             = ["code"]
  client_basic_secret        = "something_from_somewhere"
  client_id                  = "something_from_somewhere"
  token_endpoint_auth_method = "client_secret_basic"
  consent_method             = "TRUSTED"
  wildcard_redirect          = "DISABLED"
  groups_claim {
    type  = "EXPRESSION"
    value = "aa"
    name  = "bb"
  }
	`

	testAppSecurePasswordStore = `
  username_field     = "user"
  password_field     = "pass"
  credentials_scheme = "ADMIN_SETS_CREDENTIALS"
	`

	testAppSharedCredentials = `
  button_field                     = "btn-login"
  username_field                   = "txtbox-username"
  password_field                   = "txtbox-password"
  url                              = "https://example.com/login-updated.html"
  redirect_url                     = "https://example.com/redirect_url"
  checkbox                         = "checkbox_red"
  user_name_template               = "user.firstName"
  user_name_template_type          = "CUSTOM"
  user_name_template_suffix        = "hello"
  shared_password                  = "sharedpass"
  shared_username                  = "sharedusername"
  accessibility_self_service       = false
  accessibility_error_redirect_url = "https://example.com/redirect_url_1"
  auto_submit_toolbar = true
  hide_ios            = true
	`
	testAppSWA = `
  button_field   = "btn-login"
  password_field = "txtbox-password"
  username_field = "txtbox-username"
	`

	testAppThreeField = `
  button_selector      = "btn"
  username_selector    = "user"
  password_selector    = "pass"
  extra_field_selector = "third"
  extra_field_value    = "third"
`
)
