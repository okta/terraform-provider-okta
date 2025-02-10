package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaCaptcha_crud(t *testing.T) {
	mgr := newFixtureManager("resources", captcha, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", captcha)
	oktaResourceTest(
		t, resource.TestCase{
			PreCheck:          testAccPreCheck(t),
			ErrorCheck:        testAccErrorChecks(t),
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      checkResourceDestroy(captcha, doesCaptchaExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
						resource.TestCheckResourceAttr(resourceName, "type", "HCAPTCHA"),
						resource.TestCheckResourceAttr(resourceName, "site_key", "random_key"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)+"_updated"),
						resource.TestCheckResourceAttr(resourceName, "type", "HCAPTCHA"),
						resource.TestCheckResourceAttr(resourceName, "site_key", "random_key_updated")),
				},
			},
		})
}

func doesCaptchaExist(id string) (bool, error) {
	client := sdkSupplementClientForTest()
	_, response, err := client.GetCaptcha(context.Background(), id)
	return doesResourceExist(response, err)
}
