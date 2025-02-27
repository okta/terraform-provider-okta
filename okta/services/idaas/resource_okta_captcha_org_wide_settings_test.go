package idaas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaCaptchaOrgWideSettings_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSCaptchaOrgWideSettings, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("updated.tf", t)
	empty := mgr.GetFixtures("empty.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSCaptchaOrgWideSettings)
	acctest.OktaResourceTest(
		t, resource.TestCase{
			PreCheck:                 acctest.AccPreCheck(t),
			ErrorCheck:               testAccErrorChecks(t),
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
			CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSCaptchaOrgWideSettings, doesCaptchaOrgWideSettingsExist),
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
						resource.TestCheckResourceAttr(resourceName, "(resources.OktaIDaaSCaptcha)_id", ""),
					),
				},
			},
		})
}

func doesCaptchaOrgWideSettingsExist(string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
	settings, _, err := client.GetOrgWideCaptchaSettings(context.Background())
	if err != nil {
		return false, err
	}
	return settings != nil && settings.CaptchaId != nil, nil
}
