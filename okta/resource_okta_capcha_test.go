package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaCaptcha(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(captcha)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", captcha)
	resource.Test(
		t, resource.TestCase{
			PreCheck:          testAccPreCheck(t),
			ErrorCheck:        testAccErrorChecks(t),
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      createCheckResourceDestroy(captcha, doesCaptchaExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
						resource.TestCheckResourceAttr(resourceName, "type", "HCAPTCHA"),
						resource.TestCheckResourceAttr(resourceName, "site_key", "random_key"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)+"_updated"),
						resource.TestCheckResourceAttr(resourceName, "type", "HCAPTCHA"),
						resource.TestCheckResourceAttr(resourceName, "site_key", "random_key_updated")),
				},
			},
		})
}

func doesCaptchaExist(id string) (bool, error) {
	_, response, err := getSupplementFromMetadata(testAccProvider.Meta()).GetCaptcha(context.Background(), id)
	return doesResourceExist(response, err)
}
