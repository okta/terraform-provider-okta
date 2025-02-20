package idaas_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaInlineHook_crud(t *testing.T) {
	resourceName := "okta_inline_hook.test"
	mgr := newFixtureManager("resources", resources.OktaIDaaSInlineHook, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	activatedConfig := mgr.GetFixtures("basic_activated.tf", t)
	registration := mgr.GetFixtures("registration.tf", t)
	passwordImport := mgr.GetFixtures("password_import.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: acctest.AccProvidersFactoriesForTest(),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSInlineHook, inlineHookExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.oauth2.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "channel.method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),

					resource.TestCheckResourceAttr("okta_inline_hook.twilio", "type", "com.okta.telephony.provider"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.import.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test1"),
					resource.TestCheckResourceAttr(resourceName, "channel.method", "POST"),
				),
			},
			{
				Config: activatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.import.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test1"),
					resource.TestCheckResourceAttr(resourceName, "channel.method", "POST"),
				),
			},
			{
				Config: registration,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.user.pre-registration"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test1"),
					resource.TestCheckResourceAttr(resourceName, "channel.method", "POST"),
				),
			},
			{
				Config: passwordImport,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.user.credential.password.import"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test1"),
					resource.TestCheckResourceAttr(resourceName, "channel.method", "POST"),
				),
			},
		},
	})
}

func TestAccResourceOktaInlineHook_com_okta_saml_tokens_transform(t *testing.T) {
	resourceName := "okta_inline_hook.test"
	mgr := newFixtureManager("resources", inlineHook, t.Name())

	name1 := "One"
	name2 := "Two"
	config := `
resource "okta_inline_hook" "test" {
  name    = "testAcc_replace_with_uuid_%s"
  type     = "com.okta.saml.tokens.transform"
  version  = "1.0.2"
  status   = "ACTIVE"
  channel_json = <<JSON
{
        "type": "OAUTH",
        "version": "1.0.0",
        "config": {
            "headers": [
                {
                    "key": "Field 1",
                    "value": "Value 1"
                },
                {
                    "key": "Field 2",
                    "value": "Value 2"
                }
            ],
            "method": "POST",
            "authType": "client_secret_post",
            "uri": "https://example.com/service",
            "clientId": "abc123",
            "clientSecret": "def456",
            "tokenUrl": "https://example.com/token",
            "scope": "api"
        }
}
JSON
}
	`

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(inlineHook, inlineHookExists),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(fmt.Sprintf(config, name1)),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)+"_One"),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.saml.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttrSet(resourceName, "channel_json"),
					resource.TestCheckNoResourceAttr(resourceName, "channel"),
					resource.TestCheckNoResourceAttr(resourceName, "auth"),
				),
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf(config, name2)),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)+"_Two"),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.saml.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttrSet(resourceName, "channel_json"),
					resource.TestCheckNoResourceAttr(resourceName, "channel"),
					resource.TestCheckNoResourceAttr(resourceName, "auth"),
				),
			},
		},
	})
}

func inlineHookExists(id string) (bool, error) {
	client := provider.SdkV2ClientForTest()
	_, resp, err := client.InlineHook.GetInlineHook(context.Background(), id)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return resp.StatusCode != http.StatusNotFound, nil
}
