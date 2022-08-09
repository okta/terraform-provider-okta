package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaCaptchaOrgWideSettings(t *testing.T) {
	mgr := newFixtureManager(captchaOrgWideSettings, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	empty := mgr.GetFixtures("empty.tf", t)
	resourceName := fmt.Sprintf("%s.test", captchaOrgWideSettings)
	oktaResourceTest(
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
	client := apiSupplementForTest()
	settings, _, err := client.GetOrgWideCaptchaSettings(context.Background())
	if err != nil {
		return false, err
	}
	return settings != nil && settings.CaptchaId != nil, nil
}
