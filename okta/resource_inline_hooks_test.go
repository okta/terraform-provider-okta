package okta

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccInlineHook(t *testing.T) {
	t.Skip("skipping test until EA feature is enabled in test account")
	ri := acctest.RandInt()
	resourceName := "okta_inline_hooks.test"
	mgr := newFixtureManager(inlineHook)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)

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
					resource.TestCheckResourceAttr(resourceName, "channel.auth_type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "channel.auth_key", "Authorization"),
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
					resource.TestCheckResourceAttr(resourceName, "channel.auth_type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "channel.auth_key", "Authorization"),
				),
			},
		},
	})
}

func hookExists(id string) (bool, error) {
	_, res, err := getSupplementFromMetadata(testAccProvider.Meta()).GetInlineHook(id)
	return err == nil && res.StatusCode != http.StatusNotFound, err
}
