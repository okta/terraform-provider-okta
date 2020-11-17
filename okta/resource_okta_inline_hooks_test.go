package okta

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOktaInlineHook_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "okta_inline_hook.test"
	mgr := newFixtureManager(inlineHook)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	activatedConfig := mgr.GetFixtures("basic_activated.tf", ri, t)
	registration := mgr.GetFixtures("registration.tf", ri, t)
	passwordImport := mgr.GetFixtures("password_import.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(inlineHook, inlineHookExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.oauth2.tokens.transform"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "channel.method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, inlineHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
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
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
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
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
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
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
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

func inlineHookExists(id string) (bool, error) {
	_, res, err := getSupplementFromMetadata(testAccProvider.Meta()).GetInlineHook(id)
	if err != nil && res.StatusCode != http.StatusNotFound {
		return false, err
	}
	return res.StatusCode != http.StatusNotFound, nil
}
