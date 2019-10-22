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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(authServer, hookExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, hookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
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
					ensureResourceExists(resourceName, hookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
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
					ensureResourceExists(resourceName, hookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
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
					ensureResourceExists(resourceName, hookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "com.okta.user.pre-registration"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test1"),
					resource.TestCheckResourceAttr(resourceName, "channel.method", "POST"),
				),
			},
		},
	})
}

func hookExists(id string) (bool, error) {
	_, res, err := getSupplementFromMetadata(testAccProvider.Meta()).GetInlineHook(id)
	return err == nil && res.StatusCode != http.StatusNotFound, err
}
