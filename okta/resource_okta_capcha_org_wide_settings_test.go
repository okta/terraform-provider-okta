package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaCaptchaOrgWideSettings(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(captchaOrgWideSettings)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updated := mgr.GetFixtures("updated.tf", ri, t)
	empty := mgr.GetFixtures("empty.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", captchaOrgWideSettings)
	resource.Test(
		t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProvidersFactories,
			CheckDestroy:      createCheckResourceDestroy(captchaOrgWideSettings, doesCaptchaOrgWideSettingsExist),
			Steps: []resource.TestStep{
				{
					Config: config,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "enabled_for.#", "1"),
					),
				},
				{
					Config: updated,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "enabled_for.#", "3"),
					),
				},
				{
					Config: empty,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "enabled_for.#", "0"),
						resource.TestCheckResourceAttr(resourceName, "captcha_id", ""),
					),
				},
			},
		})
}

func doesCaptchaOrgWideSettingsExist(string) (bool, error) {
	settings, _, err := getSupplementFromMetadata(testAccProvider.Meta()).GetOrgWideCaptchaSettings(context.Background())
	if err != nil {
		return false, err
	}
	return settings != nil && settings.CaptchaId != nil, nil
}
