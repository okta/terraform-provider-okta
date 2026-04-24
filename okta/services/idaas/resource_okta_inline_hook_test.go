package idaas_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
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
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSInlineHook, inlineHookExists),
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
	mgr := newFixtureManager("resources", resources.OktaIDaaSInlineHook, t.Name())
	name1 := "One"
	name2 := "Two"
	secret := secretHolder{"def456"}
	updatedSecret := secretHolder{"def456_UPDATED"}
	config := mgr.GetFixtures("okta_inline_hook_saml_clientsecret.tf", t)
	updatedConfig := mgr.GetFixtures("okta_inline_hook_saml_clientsecret_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSInlineHook, inlineHookExists),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(fmt.Sprintf(config, name1)),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)+"_One"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.saml.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttrSet(resourceName, "channel_json"),
					resource.TestCheckResourceAttrWith(resourceName, "channel_json", secret.CheckSecret),
					resource.TestCheckNoResourceAttr(resourceName, "channel"),
					resource.TestCheckNoResourceAttr(resourceName, "auth"),
				),
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf(updatedConfig, name2)),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)+"_Two"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.saml.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttrSet(resourceName, "channel_json"),
					resource.TestCheckResourceAttrWith(resourceName, "channel_json", updatedSecret.CheckSecret),
					resource.TestCheckNoResourceAttr(resourceName, "channel"),
					resource.TestCheckNoResourceAttr(resourceName, "auth"),
				),
			},
			{
				Config: mgr.ConfigReplace(fmt.Sprintf(updatedConfig, name2)),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)+"_Two"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.saml.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttrSet(resourceName, "channel_json"),
					resource.TestCheckResourceAttrWith(resourceName, "channel_json", updatedSecret.CheckSecret),
					resource.TestCheckNoResourceAttr(resourceName, "channel"),
					resource.TestCheckNoResourceAttr(resourceName, "auth"),
				),
			},
		},
	})
}

func TestAccResourceOktaInlineHook_token_channeljson(t *testing.T) {
	resourceName := fmt.Sprintf("%s.token_channeljson", resources.OktaIDaaSInlineHook)
	mgr := newFixtureManager("resources", resources.OktaIDaaSInlineHook, t.Name())
	config := mgr.GetFixtures("okta_inline_hook_token_channeljson.tf", t)
	updatedConfig := mgr.GetFixtures("okta_inline_hook_token_channeljson_updated.tf", t)
	secret := secretHolder{"some_secret"}
	updatedSecret := secretHolder{"some_secret_UPDATED"}
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSInlineHook, inlineHookExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", "Inline Hook Channel JSON"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.oauth2.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
					resource.TestCheckResourceAttrSet(resourceName, "channel_json"),
					resource.TestCheckResourceAttrWith(resourceName, "channel_json", secret.CheckSecret),
					resource.TestCheckNoResourceAttr(resourceName, "channel.version"),
					resource.TestCheckNoResourceAttr(resourceName, "auth"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", "Inline Hook Channel JSON"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.oauth2.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
					resource.TestCheckResourceAttrSet(resourceName, "channel_json"),
					resource.TestCheckResourceAttrWith(resourceName, "channel_json", updatedSecret.CheckSecret),
					resource.TestCheckNoResourceAttr(resourceName, "channel.version"),
					resource.TestCheckNoResourceAttr(resourceName, "auth"),
				),
			},
			{ // No change whatsoever, this is akin to running 'terraform plan' with no changes planned
				Config: updatedConfig, // same as previous step with no changes
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", "Inline Hook Channel JSON"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.oauth2.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
					resource.TestCheckResourceAttrSet(resourceName, "channel_json"),
					resource.TestCheckResourceAttrWith(resourceName, "channel_json", updatedSecret.CheckSecret),
					resource.TestCheckNoResourceAttr(resourceName, "channel.version"),
					resource.TestCheckNoResourceAttr(resourceName, "auth"),
				),
			},
		},
	})
}

func TestAccResourceOktaInlineHook_token_channel_auth(t *testing.T) {
	resourceName := fmt.Sprintf("%s.token_channel_auth", resources.OktaIDaaSInlineHook)
	mgr := newFixtureManager("resources", resources.OktaIDaaSInlineHook, t.Name())
	config := mgr.GetFixtures("okta_inline_hook_token_channel_auth.tf", t)
	updatedConfig := mgr.GetFixtures("okta_inline_hook_token_channel_auth_updated.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSInlineHook, inlineHookExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", "Inline Hook Channel Auth"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.oauth2.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "channel.method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.value", "secret"),
					resource.TestCheckNoResourceAttr(resourceName, "channel_json"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", "Inline Hook Channel Auth"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.oauth2.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "channel.method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.value", "secret_UPDATED"),
					resource.TestCheckNoResourceAttr(resourceName, "channel_json"),
				),
			},
			{ // No change whatsoever, this is akin to running 'terraform plan' with no changes planned
				Config: updatedConfig, // same as previous step with no changes
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", "Inline Hook Channel Auth"),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.oauth2.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "channel.method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.value", "secret_UPDATED"),
					resource.TestCheckNoResourceAttr(resourceName, "channel_json"),
				),
			},
		},
	})
}

func inlineHookExists(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
	_, resp, err := client.InlineHook.GetInlineHook(context.Background(), id)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return resp.StatusCode != http.StatusNotFound, nil
}

type secretHolder struct {
	secret string
}

func (s *secretHolder) CheckSecret(value string) error {
	if strings.Contains(value, s.secret) {
		return nil
	}
	return fmt.Errorf("%v doesn't seem to match expected value of %v", value, s.secret)
}
